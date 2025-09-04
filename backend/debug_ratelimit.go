package main

import (
	"Stone/backend/pkg/ratelimit"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("无法连接到MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// 初始化速率限制器
	rateLimitCollection := client.Database("stoneDB").Collection("rate_limits")
	limiter := ratelimit.Init(rateLimitCollection)
	ratelimit.SetGlobalLimiter(limiter)

	testIP := "127.0.0.1"
	fmt.Printf("测试IP: %s\n", testIP)
	
	// 连续发送10个请求测试
	for i := 1; i <= 10; i++ {
		allowed, action := ratelimit.IsAllowed(testIP)
		status := "ALLOWED"
		if !allowed {
			status = fmt.Sprintf("BLOCKED (%s)", action)
		}
		fmt.Printf("请求 %2d: %s\n", i, status)
	}
	
	// 获取统计信息
	stats := ratelimit.GetStats()
	fmt.Printf("\n统计信息:\n")
	fmt.Printf("  活跃令牌桶数: %d\n", stats.ActiveBuckets)
	fmt.Printf("  最大令牌桶数: %d\n", stats.MaxBuckets)
	fmt.Printf("  规则数量: %d\n", stats.RulesCount)
}