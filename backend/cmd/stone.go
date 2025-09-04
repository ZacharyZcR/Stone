package main

import (
	"Stone/backend/pkg/api"
	"Stone/backend/pkg/api/handlers"
	"Stone/backend/pkg/capture"
	"Stone/backend/pkg/config"
	"Stone/backend/pkg/logging"
	"Stone/backend/pkg/monitoring"
	"Stone/backend/pkg/rules"
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
	err := logging.InitStorage("localhost:6379", "mongodb://localhost:27019", "stoneDB", "logs")
	if err != nil {
		logging.LogError(fmt.Errorf("初始化存储失败: %v", err))
		return
	}

	// 获取MongoDB集合
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27019"))
	if err != nil {
		logging.LogError(fmt.Errorf("无法连接到MongoDB: %v", err))
		return
	}
	defer client.Disconnect(context.Background())

	configCollection := client.Database("stoneDB").Collection("config")
	rulesCollection := client.Database("stoneDB").Collection("rules")
	userCollection := client.Database("stoneDB").Collection("users") // 新增的用户集合
	metricsCollection := client.Database("stoneDB").Collection("metrics")

	// 设置集合
	config.SetMongoCollection(configCollection)
	rules.SetMongoCollection(rulesCollection)
	monitoring.SetMongoCollection(metricsCollection)
	handlers.SetUserCollection(userCollection) // 设置用户集合
	handlers.SetMetricsCollection(metricsCollection)

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

	logging.LogInfo(fmt.Sprintf("服务器将在端口 %d 上运行", cfg.Server.Port))
	logging.LogInfo(fmt.Sprintf("防火墙模式: %s", cfg.Firewall.Mode))
	logging.LogInfo(fmt.Sprintf("规则文件: %s", cfg.Firewall.RulesFile))

	go func() {
		if err := capture.StartCapture(cfg.Server.Port, cfg.Firewall.TargetAddress); err != nil {
			logging.LogError(fmt.Errorf("启动流量捕获失败: %v", err))
		}
	}()

	router := api.SetupRouter(configCollection, userCollection) // 传递用户集合
	if router == nil {
		log.Fatal("创建路由失败：请设置JWT_SECRET环境变量")
	}
	if err := router.Run(":8083"); err != nil {
		log.Fatalf("启动API服务失败: %v", err)
	}
}
