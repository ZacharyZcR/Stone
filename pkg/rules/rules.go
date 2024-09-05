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
	Name  string `mapstructure:"name"`
	Regex string `mapstructure:"regex"`
}

// InterceptionRules 用于存储拦截规则
type InterceptionRules struct {
	URLPatterns  []Pattern `mapstructure:"url_patterns"`
	BodyPatterns []Pattern `mapstructure:"body_patterns"`
}

// IPControlRules 用于存储IP控制规则
type IPControlRules struct {
	Whitelist []string `mapstructure:"whitelist"`
	Blacklist []string `mapstructure:"blacklist"`
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

// CheckRequest 检查请求的URL和包体
func CheckRequest(req *http.Request) bool {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	// 检查URL
	for _, pattern := range interceptionRules.URLPatterns {
		matched, err := regexp.MatchString(pattern.Regex, req.URL.Path)
		if err != nil {
			continue // 如果正则表达式有问题，跳过此规则
		}
		if matched {
			return false
		}
	}

	// 检查包体
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return false
		}
		req.Body.Close() // 关闭后重新设置Body，以便后续使用
		req.Body = ioutil.NopCloser(strings.NewReader(string(body)))

		for _, pattern := range interceptionRules.BodyPatterns {
			matched, err := regexp.MatchString(pattern.Regex, string(body))
			if err != nil {
				continue // 如果正则表达式有问题，跳过此规则
			}
			if matched {
				return false
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

// UpdateIPControlRules 动态更新IP控制规则
func UpdateIPControlRules(newWhitelist, newBlacklist []string) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	ipControlRules.Whitelist = newWhitelist
	ipControlRules.Blacklist = newBlacklist

	// 更新MongoDB中的IP控制规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "ip_control"},
		bson.M{
			"$set": bson.M{
				"whitelist": newWhitelist,
				"blacklist": newBlacklist,
			},
		},
	)
	return err
}

// UpdateInterceptionRules 动态更新拦截规则
func UpdateInterceptionRules(newURLPatterns, newBodyPatterns []Pattern) error {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	interceptionRules.URLPatterns = newURLPatterns
	interceptionRules.BodyPatterns = newBodyPatterns

	// 更新MongoDB中的拦截规则
	_, err := mongoCollection.UpdateOne(
		context.Background(),
		bson.M{"type": "interception"},
		bson.M{
			"$set": bson.M{
				"url_patterns":  newURLPatterns,
				"body_patterns": newBodyPatterns,
			},
		},
	)
	return err
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
