// test_ratelimit.go - 测试速率限制功能

package main

import (
	"Stone/backend/pkg/ratelimit"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	fmt.Println("测试速率限制功能...")

	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// 初始化速率限制器
	rateLimitCollection := client.Database("stoneDB").Collection("rate_limits")
	limiter := ratelimit.Init(rateLimitCollection)
	
	fmt.Println("速率限制器初始化成功")

	// 测试IP
	testIP := "192.168.1.100"

	// 模拟连续请求
	fmt.Printf("测试IP %s 的速率限制:\n", testIP)
	
	successCount := 0
	blockedCount := 0
	
	// 发送70个请求，应该前60个通过，后10个被限制
	for i := 1; i <= 70; i++ {
		allowed, action := ratelimit.IsAllowed(testIP)
		if allowed {
			successCount++
			fmt.Printf("请求 %d: 通过 ✓\n", i)
		} else {
			blockedCount++
			fmt.Printf("请求 %d: 被限制 ✗ (动作: %s)\n", i, action)
		}
		
		// 稍微间隔一下
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Printf("\n测试结果:\n")
	fmt.Printf("成功请求: %d\n", successCount)
	fmt.Printf("被限制请求: %d\n", blockedCount)
	
	// 获取统计信息
	stats := limiter.GetStats()
	fmt.Printf("统计信息: %+v\n", stats)
	
	// 测试规则管理
	fmt.Println("\n测试规则管理功能:")
	rules := limiter.GetRules()
	fmt.Printf("当前规则数量: %d\n", len(rules))
	for i, rule := range rules {
		fmt.Printf("规则 %d: %s (容量:%d, 补充率:%d/s)\n", 
			i+1, rule.Name, rule.Capacity, rule.RefillRate)
	}
	
	fmt.Println("\n速率限制功能测试完成!")
}