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

// FetchIPStatsFromMongo 从MongoDB中检索IP访问统计
func FetchIPStatsFromMongo(ctx context.Context, startDateTime, endDateTime time.Time, status string) ([]bson.M, error) {
	// 构建过滤条件
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": startDateTime,
			"$lte": endDateTime,
		},
	}

	// 添加状态过滤条件
	if status == "blocked" {
		filter["status"] = bson.M{"$ne": "success"}
	} else if status == "passed" {
		filter["status"] = "success"
	}

	// 构建聚合管道
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$client_ip",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	// 执行聚合查询
	cursor, err := mongoCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("执行聚合查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析聚合结果失败: %v", err)
	}

	return results, nil
}

// AttackerProfile 结构体定义攻击者画像
type AttackerProfile struct {
	DailyAttacks       map[string]DailyStats `json:"daily_attacks"`
	MostActiveDay      string                `json:"most_active_day"`
	HourlyDistribution map[int]int           `json:"hourly_distribution"`
	MostVisitedURL     string                `json:"most_visited_url"`
	TotalVisits        int                   `json:"total_visits"`
	TotalNormalVisits  int                   `json:"total_normal_visits"`
	TotalAttacks       int                   `json:"total_attacks"`
}

// DailyStats 结构体定义每日统计
type DailyStats struct {
	Total   int `json:"total"`
	Attacks int `json:"attacks"`
	Normal  int `json:"normal"`
}

// FetchAttackerProfile 从MongoDB中获取攻击者画像
func FetchAttackerProfile(ctx context.Context, ip string, endTime time.Time) (AttackerProfile, error) {
	startTime := endTime.AddDate(0, 0, -6) // 获取7天的数据（包括当天）
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())

	profile := AttackerProfile{
		DailyAttacks:       make(map[string]DailyStats),
		HourlyDistribution: make(map[int]int),
		TotalVisits:        0,
		TotalNormalVisits:  0,
		TotalAttacks:       0,
	}

	// 构建基础过滤条件
	filter := bson.M{
		"client_ip": ip,
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	// 1. 获取每日攻击和正常访问次数分布
	dailyPipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$timestamp"}},
			"total": bson.M{"$sum": 1},
			"attacks": bson.M{"$sum": bson.M{
				"$cond": []interface{}{bson.M{"$ne": []interface{}{"$status", "success"}}, 1, 0},
			}},
		}},
		{"$sort": bson.M{"_id": 1}},
	}

	dailyCursor, err := mongoCollection.Aggregate(ctx, dailyPipeline)
	if err != nil {
		return profile, err
	}
	defer dailyCursor.Close(ctx)

	var dailyResults []bson.M
	if err = dailyCursor.All(ctx, &dailyResults); err != nil {
		return profile, err
	}

	// 填充7天的数据，没有数据的日期用0填充
	for i := 0; i < 7; i++ {
		date := startTime.AddDate(0, 0, i).Format("2006-01-02")
		profile.DailyAttacks[date] = DailyStats{Total: 0, Attacks: 0, Normal: 0}
	}

	var maxCount int
	// 在处理每日数据时，累计总数
	for _, result := range dailyResults {
		date := result["_id"].(string)
		total := int(result["total"].(int32))
		attacks := int(result["attacks"].(int32))
		normal := total - attacks

		profile.DailyAttacks[date] = DailyStats{
			Total:   total,
			Attacks: attacks,
			Normal:  normal,
		}

		profile.TotalVisits += total
		profile.TotalAttacks += attacks
		profile.TotalNormalVisits += normal

		if total > maxCount {
			maxCount = total
			profile.MostActiveDay = date
		}
	}

	mostActiveDayStart, _ := time.Parse("2006-01-02", profile.MostActiveDay)
	mostActiveDayEnd := mostActiveDayStart.Add(24 * time.Hour)

	hourlyFilter := bson.M{
		"client_ip": ip,
		"timestamp": bson.M{
			"$gte": mostActiveDayStart,
			"$lt":  mostActiveDayEnd,
		},
	}

	hourlyPipeline := []bson.M{
		{"$match": hourlyFilter},
		{"$project": bson.M{
			"hour": bson.M{
				"$hour": bson.M{
					"$add": []interface{}{"$timestamp", 8 * 60 * 60 * 1000}, // 添加8小时
				},
			},
		}},
		{"$group": bson.M{
			"_id":   "$hour",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"_id": 1}},
	}

	hourlyCursor, err := mongoCollection.Aggregate(ctx, hourlyPipeline)
	if err != nil {
		return profile, err
	}
	defer hourlyCursor.Close(ctx)

	var hourlyResults []bson.M
	if err = hourlyCursor.All(ctx, &hourlyResults); err != nil {
		return profile, err
	}

	// 填充24小时的数据，没有数据的小时用0填充
	for i := 0; i < 24; i++ {
		profile.HourlyDistribution[i] = 0
	}

	for _, result := range hourlyResults {
		hour := int(result["_id"].(int32))
		count := int(result["count"].(int32))
		profile.HourlyDistribution[hour] = count
	}

	// 3. 获取7天内最常访问的URL
	urlPipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$url",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 1},
	}

	urlCursor, err := mongoCollection.Aggregate(ctx, urlPipeline)
	if err != nil {
		return profile, err
	}
	defer urlCursor.Close(ctx)

	var urlResults []bson.M
	if err = urlCursor.All(ctx, &urlResults); err != nil {
		return profile, err
	}

	if len(urlResults) > 0 {
		profile.MostVisitedURL = urlResults[0]["_id"].(string)
	}

	return profile, nil
}
