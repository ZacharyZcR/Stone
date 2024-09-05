// pkg/rules/rules.go

package rules

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Pattern struct {
	Name   string `bson:"name" json:"name"`
	Regex  string `bson:"regex" json:"regex"`
	Method string `bson:"method" json:"method"` // 添加HTTP请求方法
}

// InterceptionRules 用于存储拦截规则
type InterceptionRules struct {
	Rules []Pattern `bson:"rules" json:"rules"`
}

// IPControlRules 用于存储IP控制规则
type IPControlRules struct {
	Whitelist []string `mapstructure:"whitelist"`
	Blacklist []string `mapstructure:"blacklist"`
}

type IPControlRule struct {
	IP        string `bson:"ip" json:"ip"`
	IsAllowed bool   `bson:"is_allowed" json:"-"`
	Type      string `json:"type"` // 用于解析请求体中的类型
}

var (
	interceptionRules InterceptionRules
	ipControlRules    IPControlRules
	rulesMutex        sync.RWMutex
	mongoCollection   *mongo.Collection // 假设已初始化
)

// SetMongoCollection 设置MongoDB集合
func SetMongoCollection(collection *mongo.Collection) {
	mongoCollection = collection
}

// LoadInterceptionRules 从MongoDB加载拦截规则
func LoadInterceptionRules(ctx context.Context) (*InterceptionRules, error) {
	var rules InterceptionRules
	err := mongoCollection.FindOne(ctx, bson.M{"type": "interception"}).Decode(&rules)
	if err != nil {
		return nil, fmt.Errorf("从MongoDB读取拦截规则失败: %w", err)
	}

	rulesMutex.Lock()
	interceptionRules = rules
	rulesMutex.Unlock()

	return &rules, nil
}

// LoadIPControlRules 从MongoDB加载IP控制规则
func LoadIPControlRules(ctx context.Context) (*IPControlRules, error) {
	var rules IPControlRules
	err := mongoCollection.FindOne(ctx, bson.M{"type": "ip_control"}).Decode(&rules)
	if err != nil {
		return nil, fmt.Errorf("从MongoDB读取IP控制规则失败: %w", err)
	}

	rulesMutex.Lock()
	ipControlRules = rules
	rulesMutex.Unlock()

	return &rules, nil
}

// CheckRequest 检查请求的URL、包体和头部
func CheckRequest(req *http.Request) bool {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	// 读取请求体
	var body []byte
	if req.Body != nil {
		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return false
		}
		req.Body.Close() // 关闭后重新设置Body，以便后续使用
		req.Body = ioutil.NopCloser(strings.NewReader(string(body)))
	}

	// 检查规则
	for _, pattern := range interceptionRules.Rules {
		// 检查HTTP方法
		if pattern.Method != "" && pattern.Method != req.Method {
			continue
		}

		// 检查URL
		matched, err := regexp.MatchString(pattern.Regex, req.URL.Path)
		if err != nil {
			continue // 如果正则表达式有问题，跳过此规则
		}
		if matched {
			return false
		}

		// 检查包体
		if len(body) > 0 {
			matched, err = regexp.MatchString(pattern.Regex, string(body))
			if err != nil {
				continue // 如果正则表达式有问题，跳过此规则
			}
			if matched {
				return false
			}
		}

		// 检查头部
		for _, values := range req.Header {
			for _, value := range values {
				matched, err = regexp.MatchString(pattern.Regex, value)
				if err != nil {
					continue // 如果正则表达式有问题，跳过此规则
				}
				if matched {
					return false
				}
			}
		}
	}

	return true
}

// IsAllowed 检查IP是否被允许
func IsAllowed(ip string) (allowed bool, inWhitelist bool) {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	// 检查黑名单
	for _, blockedIP := range ipControlRules.Blacklist {
		if blockedIP == ip {
			return false, false
		}
	}

	// 检查白名单
	for _, allowedIP := range ipControlRules.Whitelist {
		if allowedIP == ip {
			return true, true
		}
	}

	// 如果不在白名单或黑名单中，返回中性结果
	return true, false
}

// GetInterceptionRules 获取当前拦截规则
func GetInterceptionRules() InterceptionRules {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()
	return interceptionRules
}

// GetIPControlRules 获取当前IP控制规则
func GetIPControlRules() IPControlRules {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()
	return ipControlRules
}

// GetIPRule 获取特定IP的规则
func GetIPRule(ip string) (IPControlRule, bool) {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	for _, allowedIP := range ipControlRules.Whitelist {
		if allowedIP == ip {
			return IPControlRule{IP: ip, IsAllowed: true}, true
		}
	}
	for _, blockedIP := range ipControlRules.Blacklist {
		if blockedIP == ip {
			return IPControlRule{IP: ip, IsAllowed: false}, true
		}
	}
	return IPControlRule{}, false
}

// AddIPRule 添加新的IP规则
func AddIPRule(rule IPControlRule) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	// 根据类型设置IsAllowed
	if rule.Type == "whitelist" {
		rule.IsAllowed = true
	} else if rule.Type == "blacklist" {
		rule.IsAllowed = false
	} else {
		return fmt.Errorf("无效的IP规则类型")
	}

	if rule.IsAllowed {
		ipControlRules.Whitelist = append(ipControlRules.Whitelist, rule.IP)
	} else {
		ipControlRules.Blacklist = append(ipControlRules.Blacklist, rule.IP)
	}

	// 更新MongoDB中的IP控制规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "ip_control"},
		bson.M{
			"$set": bson.M{
				"whitelist": ipControlRules.Whitelist,
				"blacklist": ipControlRules.Blacklist,
			},
		},
	)
	return err
}

// DeleteIPRule 删除特定IP的规则
func DeleteIPRule(ip string) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	// 从白名单或黑名单中删除IP
	for i, allowedIP := range ipControlRules.Whitelist {
		if allowedIP == ip {
			ipControlRules.Whitelist = append(ipControlRules.Whitelist[:i], ipControlRules.Whitelist[i+1:]...)
			break
		}
	}
	for i, blockedIP := range ipControlRules.Blacklist {
		if blockedIP == ip {
			ipControlRules.Blacklist = append(ipControlRules.Blacklist[:i], ipControlRules.Blacklist[i+1:]...)
			break
		}
	}

	// 更新MongoDB中的IP控制规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "ip_control"},
		bson.M{
			"$set": bson.M{
				"whitelist": ipControlRules.Whitelist,
				"blacklist": ipControlRules.Blacklist,
			},
		},
	)
	return err
}

// GetInterceptionRule 获取特定名称的规则
func GetInterceptionRule(name string) (Pattern, bool) {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	for _, rule := range interceptionRules.Rules {
		if rule.Name == name {
			return rule, true
		}
	}
	return Pattern{}, false
}

// AddInterceptionRule 添加新的拦截规则
func AddInterceptionRule(rule Pattern) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	interceptionRules.Rules = append(interceptionRules.Rules, rule)

	// 更新MongoDB中的拦截规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "interception"},
		bson.M{
			"$set": bson.M{
				"rules": interceptionRules.Rules,
			},
		},
	)
	return err
}

// DeleteInterceptionRule 删除特定名称的规则
func DeleteInterceptionRule(name string) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	for i, rule := range interceptionRules.Rules {
		if rule.Name == name {
			interceptionRules.Rules = append(interceptionRules.Rules[:i], interceptionRules.Rules[i+1:]...)
			break
		}
	}

	// 更新MongoDB中的拦截规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "interception"},
		bson.M{
			"$set": bson.M{
				"rules": interceptionRules.Rules,
			},
		},
	)
	return err
}
