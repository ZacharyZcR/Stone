// pkg/rules/rules.go

package rules

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

// Rules 结构体用于存储白名单和黑名单
type Rules struct {
	Whitelist []string
	Blacklist []string
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

	// 默认不允许
	return false
}

// UpdateRules 动态更新规则
func UpdateRules(newWhitelist, newBlacklist []string) {
	rulesMutex.Lock()
	currentRules.Whitelist = newWhitelist
	currentRules.Blacklist = newBlacklist
	rulesMutex.Unlock()
}

// GetCurrentRules 获取当前规则
func GetCurrentRules() Rules {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()
	return currentRules
}
