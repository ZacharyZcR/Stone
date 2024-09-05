package main

import (
	"Stone/pkg/api"
	"Stone/pkg/capture"
	"Stone/pkg/config"
	"Stone/pkg/logging"
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	monitoring.StartTime = time.Now()
	logging.LogInfo("启动Stone防火墙")

	// 初始化存储（Redis 和 MongoDB）
	err := logging.InitStorage("localhost:6379", "mongodb://localhost:27017", "stoneDB", "logs")
	if err != nil {
		logging.LogError(fmt.Errorf("初始化存储失败: %v", err))
		return
	}

	// 获取MongoDB集合
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logging.LogError(fmt.Errorf("无法连接到MongoDB: %v", err))
		return
	}
	defer client.Disconnect(context.Background())

	configCollection := client.Database("stoneDB").Collection("config")
	rulesCollection := client.Database("stoneDB").Collection("rules")

	// 设置集合
	config.SetMongoCollection(configCollection)
	rules.SetMongoCollection(rulesCollection)

	// 从MongoDB加载配置
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		logging.LogError(fmt.Errorf("加载配置失败: %v", err))
		return
	}

	// 从MongoDB加载规则
	_, err = rules.LoadInterceptionRules(context.Background())
	if err != nil {
		logging.LogError(fmt.Errorf("加载拦截规则失败: %v", err))
		return
	}

	_, err = rules.LoadIPControlRules(context.Background())
	if err != nil {
		logging.LogError(fmt.Errorf("加载IP控制规则失败: %v", err))
		return
	}

	monitoring.InitMonitoring()

	logging.LogInfo(fmt.Sprintf("服务器将在端口 %d 上运行", cfg.Server.Port))
	logging.LogInfo(fmt.Sprintf("防火墙模式: %s", cfg.Firewall.Mode))
	logging.LogInfo(fmt.Sprintf("规则文件: %s", cfg.Firewall.RulesFile))

	go func() {
		if err := capture.StartCapture(cfg.Server.Port, cfg.Firewall.TargetAddress); err != nil {
			logging.LogError(fmt.Errorf("启动流量捕获失败: %v", err))
		}
	}()

	router := api.SetupRouter()
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("启动API服务失败: %v", err)
	}
}
