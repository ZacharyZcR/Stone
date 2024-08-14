// pkg/rules/rules.go

package rules

import (
	"fmt"
	"github.com/spf13/viper"
)

// Rules 结构体用于存储白名单和黑名单
type Rules struct {
	Whitelist []string
	Blacklist []string
}

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

	return &rules, nil
}

// IsAllowed 检查IP是否被允许
func (r *Rules) IsAllowed(ip string) bool {
	// 检查黑名单
	for _, blockedIP := range r.Blacklist {
		if blockedIP == ip {
			return false
		}
	}

	// 检查白名单
	for _, allowedIP := range r.Whitelist {
		if allowedIP == ip {
			return true
		}
	}

	// 默认不允许
	return false
}

// UpdateRules 动态更新规则
func (r *Rules) UpdateRules(newWhitelist, newBlacklist []string) {
	r.Whitelist = newWhitelist
	r.Blacklist = newBlacklist
}
