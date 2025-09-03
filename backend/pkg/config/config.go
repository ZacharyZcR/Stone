// pkg/config/config.go

package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Server   ServerConfig
	Firewall FirewallConfig
}

type ServerConfig struct {
	Port int `bson:"port"`
}

type FirewallConfig struct {
	Mode          string `bson:"mode"`
	RulesFile     string `bson:"rulesfile"`
	TargetAddress string `bson:"targetaddress"`
}

var mongoCollection *mongo.Collection

// SetMongoCollection 设置MongoDB集合
func SetMongoCollection(collection *mongo.Collection) {
	mongoCollection = collection
}

// LoadConfig 从MongoDB加载配置文件
func LoadConfig(ctx context.Context) (*Config, error) {
	var config Config
	err := mongoCollection.FindOne(ctx, bson.M{"type": "config"}).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("从MongoDB读取配置文件失败: %w", err)
	}
	return &config, nil
}
