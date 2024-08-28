// cmd/stone.go

package main

import (
	"Stone/pkg/api"
	"Stone/pkg/capture"
	"Stone/pkg/config"
	"Stone/pkg/logging"
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"fmt"
	"log"
	"time"
)

func main() {
	// 记录程序启动时间
	monitoring.StartTime = time.Now()

	// 初始化日志
	logging.LogInfo("启动Stone防火墙")

	// 加载配置
	cfg, err := config.LoadConfig("pkg/config/default.yaml")
	if err != nil {
		logging.LogError(fmt.Errorf("加载配置失败: %v", err))
		return
	}

	// 初始化存储（Redis 和 MongoDB）
	err = logging.InitStorage("localhost:6379", "mongodb://localhost:27017", "stoneDB", "trafficLogs")
	if err != nil {
		logging.LogError(fmt.Errorf("初始化存储失败: %v", err))
		return
	}

	// 加载拦截规则
	_, err = rules.LoadInterceptionRules("pkg/rules/interception_rules.yaml")
	if err != nil {
		logging.LogError(fmt.Errorf("加载拦截规则失败: %v", err))
		return
	}

	// 加载IP控制规则
	_, err = rules.LoadIPControlRules("pkg/rules/ip_control.yaml")
	if err != nil {
		logging.LogError(fmt.Errorf("加载IP控制规则失败: %v", err))
		return
	}

	// 初始化监控
	monitoring.InitMonitoring()

	// 输出加载的配置
	logging.LogInfo(fmt.Sprintf("服务器将在端口 %d 上运行", cfg.Server.Port))
	logging.LogInfo(fmt.Sprintf("防火墙模式: %s", cfg.Firewall.Mode))
	logging.LogInfo(fmt.Sprintf("规则文件: %s", cfg.Firewall.RulesFile))

	// 启动流量捕获
	go func() {
		if err := capture.StartCapture(cfg.Server.Port, cfg.Firewall.TargetAddress); err != nil {
			logging.LogError(fmt.Errorf("启动流量捕获失败: %v", err))
		}
	}()

	// 启动API服务
	router := api.SetupRouter()
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("启动API服务失败: %v", err)
	}
}
