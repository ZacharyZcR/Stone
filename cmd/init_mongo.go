// cmd/init_mongo.go

package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 初始化MongoDB客户端
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("无法连接到MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// 获取MongoDB集合
	configCollection := client.Database("stoneDB").Collection("config")
	rulesCollection := client.Database("stoneDB").Collection("rules")
	logsCollection := client.Database("stoneDB").Collection("logs")

	// 插入配置文档
	configDoc := bson.M{
		"type": "config",
		"server": bson.M{
			"port": 8082,
		},
		"firewall": bson.M{
			"mode":          "main",
			"rulesfile":     "pkg/rules/rules.yaml",
			"targetaddress": "localhost:80",
		},
	}

	_, err = configCollection.InsertOne(context.Background(), configDoc)
	if err != nil {
		fmt.Printf("插入配置文档失败: %v\n", err)
		return
	}

	// 插入规则文档
	interceptionRulesDoc := bson.M{
		"type": "interception",
		"url_patterns": []bson.M{
			{"name": "Admin Access", "regex": "/admin"},
			{"name": "Login Access", "regex": "/login"},
		},
		"body_patterns": []bson.M{
			{"name": "SQL Injection - Drop", "regex": "DROP TABLE"},
			{"name": "SQL Injection - Select", "regex": "SELECT \\* FROM"},
		},
	}

	ipControlRulesDoc := bson.M{
		"type":      "ip_control",
		"whitelist": []string{"192.168.1.100", "10.0.0.1"},
		"blacklist": []string{"192.168.1.200", "10.0.0.2"},
	}

	_, err = rulesCollection.InsertOne(context.Background(), interceptionRulesDoc)
	if err != nil {
		fmt.Printf("插入拦截规则文档失败: %v\n", err)
		return
	}

	_, err = rulesCollection.InsertOne(context.Background(), ipControlRulesDoc)
	if err != nil {
		fmt.Printf("插入IP控制规则文档失败: %v\n", err)
		return
	}

	// 初始化一个空的日志集合
	_, err = logsCollection.InsertOne(context.Background(), bson.M{"initialized": true})
	if err != nil {
		fmt.Printf("初始化日志集合失败: %v\n", err)
		return
	}

	fmt.Println("初始化完成")
}
