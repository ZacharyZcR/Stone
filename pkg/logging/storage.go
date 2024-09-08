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

// SetMongoCollection 设置MongoDB集合
func SetMongoCollection(collection *mongo.Collection) {
	mongoCollection = collection
}

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
	// 确保 Redis 和 MongoDB 客户端已初始化
	if redisClient == nil || mongoCollection == nil {
		return fmt.Errorf("Redis或MongoDB客户端未初始化")
	}

	// 将日志数据转换为JSON字符串
	logDataJSON, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 将日志保存到Redis
	redisKey := fmt.Sprintf("log:%d", time.Now().UnixNano())
	err = redisClient.Set(ctx, redisKey, logDataJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("保存到Redis失败: %v", err)
	}

	// 将日志保存到MongoDB
	_, err = mongoCollection.InsertOne(ctx, logData)
	if err != nil {
		// 如果MongoDB存储失败，可以选择在Redis中标记此日志为未同步，稍后重试
		redisClient.Set(ctx, fmt.Sprintf("%s:unsynced", redisKey), logDataJSON, 0)
		return fmt.Errorf("保存到MongoDB失败: %v", err)
	}

	return nil
}

// FetchLogsFromMongoWithFilters 从MongoDB中检索日志，支持过滤和分页
func FetchLogsFromMongoWithFilters(ctx context.Context, page, pageSize int, startDateTime, endDateTime time.Time, ip, status string) ([]bson.M, int64, error) {
	// 构建过滤条件
	filter := bson.D{}

	// 添加时间过滤条件
	if !startDateTime.IsZero() || !endDateTime.IsZero() {
		timeFilter := bson.D{}
		if !startDateTime.IsZero() {
			timeFilter = append(timeFilter, bson.E{"$gte", startDateTime})
		}
		if !endDateTime.IsZero() {
			timeFilter = append(timeFilter, bson.E{"$lte", endDateTime})
		}
		filter = append(filter, bson.E{"timestamp", timeFilter})
	}

	// 添加IP过滤条件
	if ip != "" {
		filter = append(filter, bson.E{"client_ip", ip})
	}

	// 添加状态过滤条件
	if status != "" {
		if status == "blocked" {
			filter = append(filter, bson.E{"status", bson.M{"$ne": "success"}})
		} else if status == "passed" {
			filter = append(filter, bson.E{"status", "success"})
		}
	}

	// 计算总记录数
	totalCount, err := mongoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("计算总记录数失败: %v", err)
	}

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"timestamp", -1}}) // 按时间倒序排列
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	// 执行查询
	cursor, err := mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("检索日志失败: %v", err)
	}
	defer cursor.Close(ctx)

	var logs []bson.M
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, fmt.Errorf("解析日志失败: %v", err)
	}

	return logs, totalCount, nil
}
