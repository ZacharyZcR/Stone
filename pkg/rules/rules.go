// pkg/rules/rules.go

package rules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/viper"
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
)

// LoadInterceptionRules 加载拦截规则文件
func LoadInterceptionRules(rulesPath string) (*InterceptionRules, error) {
	viper.SetConfigFile(rulesPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取拦截规则文件失败: %w", err)
	}

	var rules InterceptionRules
	if err := viper.Unmarshal(&rules); err != nil {
		return nil, fmt.Errorf("解析拦截规则文件失败: %w", err)
	}

	rulesMutex.Lock()
	interceptionRules = rules
	rulesMutex.Unlock()

	return &rules, nil
}

// LoadIPControlRules 加载IP控制规则文件
func LoadIPControlRules(rulesPath string) (*IPControlRules, error) {
	viper.SetConfigFile(rulesPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取IP控制规则文件失败: %w", err)
	}

	var rules IPControlRules
	if err := viper.Unmarshal(&rules); err != nil {
		return nil, fmt.Errorf("解析IP控制规则文件失败: %w", err)
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
func UpdateIPControlRules(newWhitelist, newBlacklist []string) {
	rulesMutex.Lock()
	ipControlRules.Whitelist = newWhitelist
	ipControlRules.Blacklist = newBlacklist
	rulesMutex.Unlock()
}

// UpdateInterceptionRules 动态更新拦截规则
func UpdateInterceptionRules(newURLPatterns, newBodyPatterns []Pattern) {
	rulesMutex.Lock()
	defer rulesMutex.Unlock()

	interceptionRules.URLPatterns = newURLPatterns
	interceptionRules.BodyPatterns = newBodyPatterns
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
