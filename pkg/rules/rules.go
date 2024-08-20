// pkg/rules/rules.go

package rules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Rules 结构体用于存储各种规则
type Rules struct {
	Whitelist    []string
	Blacklist    []string
	URLPatterns  []string // URL模式
	BodyPatterns []string // 包体模式
}

var (
	currentRules Rules
	rulesMutex   sync.RWMutex
)

// LoadRules 加载规则文件
func LoadRules(rulesPath string) (*Rules, error) {
	viper.SetConfigFile(rulesPath)
	viper.SetConfigType("yaml")

	// 读取规则文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取规则文件失败: %w", err)
	}

	var rules Rules
	if err := viper.Unmarshal(&rules); err != nil {
		return nil, fmt.Errorf("解析规则文件失败: %w", err)
	}

	// 初始化全局规则
	rulesMutex.Lock()
	currentRules = rules
	rulesMutex.Unlock()

	return &rules, nil
}

// IsAllowed 检查IP是否被允许
func IsAllowed(ip string) bool {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	// 检查黑名单
	for _, blockedIP := range currentRules.Blacklist {
		if blockedIP == ip {
			return false
		}
	}

	// 检查白名单
	for _, allowedIP := range currentRules.Whitelist {
		if allowedIP == ip {
			return true
		}
	}

	// 如果不在白名单或黑名单中，返回中性结果
	return true
}

// CheckRequest 检查请求的URL和包体
func CheckRequest(req *http.Request) bool {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()

	// 检查URL
	for _, pattern := range currentRules.URLPatterns {
		if strings.Contains(req.URL.String(), pattern) {
			return false
		}
	}

	// 检查包体
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false
	}
	req.Body.Close() // 关闭后重新设置Body，以便后续使用
	req.Body = ioutil.NopCloser(strings.NewReader(string(body)))

	for _, pattern := range currentRules.BodyPatterns {
		if strings.Contains(string(body), pattern) {
			return false
		}
	}

	return true
}

// UpdateRules 动态更新规则
func UpdateRules(newWhitelist, newBlacklist, newURLPatterns, newBodyPatterns []string) {
	rulesMutex.Lock()
	currentRules.Whitelist = newWhitelist
	currentRules.Blacklist = newBlacklist
	currentRules.URLPatterns = newURLPatterns
	currentRules.BodyPatterns = newBodyPatterns
	rulesMutex.Unlock()
}

// GetCurrentRules 获取当前规则
func GetCurrentRules() Rules {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()
	return currentRules
}
