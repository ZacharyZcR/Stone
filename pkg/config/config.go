// pkg/config/config.go

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 结构体用于存储配置参数
type Config struct {
	Server   ServerConfig
	Firewall FirewallConfig
}

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Port int
}

// FirewallConfig 防火墙相关配置
type FirewallConfig struct {
	Mode          string
	RulesFile     string
	TargetAddress string // 添加目标地址字段
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
