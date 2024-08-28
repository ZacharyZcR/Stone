// pkg/logging/storage.go

package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	redisClient     *redis.Client
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	ctx             = context.Background()
)

// InitStorage 初始化Redis和MongoDB连接
func InitStorage(redisAddr, mongoURI, mongoDB, mongoCollectionName string) error {
	// 初始化Redis客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 测试Redis连接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("无法连接到Redis: %v", err)
	}

	// 初始化MongoDB客户端
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("无法连接到MongoDB: %v", err)
	}

	// 测试MongoDB连接
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("无法连接到MongoDB: %v", err)
	}

	// 获取MongoDB集合
	mongoCollection = mongoClient.Database(mongoDB).Collection(mongoCollectionName)

	return nil
}

// LogTraffic 保存流量日志到Redis和MongoDB
func LogTraffic(logData map[string]interface{}) error {
	// 将日志数据转换为JSON字符串
	logDataJSON, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 将日志保存到Redis
	err = redisClient.Set(ctx, fmt.Sprintf("log:%d", time.Now().UnixNano()), logDataJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("保存到Redis失败: %v", err)
	}

	// 将日志保存到MongoDB
	_, err = mongoCollection.InsertOne(ctx, logData)
	if err != nil {
		return fmt.Errorf("保存到MongoDB失败: %v", err)
	}

	return nil
}

// FetchLogsFromMongoWithFilters 从MongoDB中检索日志，支持过滤
func FetchLogsFromMongoWithFilters(ctx context.Context, limit int64, startDateTime, endDateTime time.Time, ip string) ([]bson.M, error) {
	// 构建过滤条件
	filter := bson.D{}

	// 添加时间过滤条件
	timeFilter := bson.D{}
	if !startDateTime.IsZero() {
		timeFilter = append(timeFilter, bson.E{"$gte", startDateTime})
	}
	if !endDateTime.IsZero() {
		timeFilter = append(timeFilter, bson.E{"$lte", endDateTime})
	}
	if len(timeFilter) > 0 {
		filter = append(filter, bson.E{"timestamp", timeFilter})
	}

	// 添加IP过滤条件
	if ip != "" {
		filter = append(filter, bson.E{"client_ip", ip})
	}

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"timestamp", -1}}) // 按时间倒序排列

	// 执行查询
	cursor, err := mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("检索日志失败: %v", err)
	}
	defer cursor.Close(ctx)

	var logs []bson.M
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("解析日志失败: %v", err)
	}

	return logs, nil
}
