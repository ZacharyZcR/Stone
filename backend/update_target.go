package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

	// 更新配置
	configCollection := client.Database("stoneDB").Collection("config")
	filter := bson.M{"type": "config"}
	update := bson.M{"$set": bson.M{"firewall.targetaddress": "localhost:9001"}}
	
	result, err := configCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("更新配置失败: %v\n", err)
		return
	}
	
	fmt.Printf("配置更新成功，匹配文档数: %d, 修改文档数: %d\n", result.MatchedCount, result.ModifiedCount)
}