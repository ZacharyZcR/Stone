// pkg/logging/storage.go

package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// 至少要有MongoDB客户端初始化
	if mongoCollection == nil {
		return fmt.Errorf("MongoDB客户端未初始化")
	}

	// 将日志数据转换为JSON字符串
	logDataJSON, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 尝试将日志保存到Redis（非阻塞）
	if redisClient != nil {
		redisKey := fmt.Sprintf("log:%d", time.Now().UnixNano())
		err := redisClient.Set(ctx, redisKey, logDataJSON, time.Hour).Err()
		if err != nil {
			// Redis失败只记录，不阻塞业务
			fmt.Printf("Redis保存失败(非阻塞): %v\n", err)
		}
	}

	// 将日志保存到MongoDB（主要存储）
	_, err = mongoCollection.InsertOne(ctx, logData)
	if err != nil {
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
	// 基础信息
	ThreatLevel  string `json:"threat_level"`
	AttackCount  int    `json:"attack_count"`
	FirstSeen    string `json:"first_seen"`
	LastSeen     string `json:"last_seen"`
	
	// 地理位置信息（简化处理）
	Country string `json:"country"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
	ASN     string `json:"asn"`
	
	// 攻击类型统计
	AttackTypes []AttackType `json:"attack_types"`
	
	// 详细统计数据
	DailyAttacks       map[string]DailyStats `json:"daily_attacks"`
	MostActiveDay      string                `json:"most_active_day"`
	HourlyDistribution map[int]int           `json:"hourly_distribution"`
	MostVisitedURL     string                `json:"most_visited_url"`
	TotalVisits        int                   `json:"total_visits"`
	TotalNormalVisits  int                   `json:"total_normal_visits"`
	TotalAttacks       int                   `json:"total_attacks"`
}

// AttackType 攻击类型统计
type AttackType struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// DailyStats 结构体定义每日统计
type DailyStats struct {
	Total   int `json:"total"`
	Attacks int `json:"attacks"`
	Normal  int `json:"normal"`
}

func FetchAttackerProfile(ctx context.Context, ip string, endTime time.Time) (AttackerProfile, error) {
	// 定义北京时区
	beijingLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return AttackerProfile{}, fmt.Errorf("failed to load Beijing timezone: %v", err)
	}

	// 将endTime转换为北京时间
	endTime = endTime.In(beijingLoc)
	startTime := endTime.AddDate(0, 0, -6) // 获取7天的数据（包括当天）
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, beijingLoc)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, beijingLoc)

	profile := AttackerProfile{
		DailyAttacks:       make(map[string]DailyStats),
		HourlyDistribution: make(map[int]int),
		TotalVisits:        0,
		TotalNormalVisits:  0,
		TotalAttacks:       0,
	}

	// 构建基础过滤条件（使用UTC时间进行查询）
	filter := bson.M{
		"client_ip": ip,
		"timestamp": bson.M{
			"$gte": startTime.UTC(),
			"$lte": endTime.UTC(),
		},
	}

	// 1. 获取每日攻击和正常访问次数分布
	dailyPipeline := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date": bson.M{"$dateToString": bson.M{
				"format": "%Y-%m-%d",
				"date":   bson.M{"$add": []interface{}{"$timestamp", 8 * 60 * 60 * 1000}}, // 转换为北京时间
			}},
			"status": 1,
		}},
		{"$group": bson.M{
			"_id":   "$date",
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

	mostActiveDayStart, _ := time.ParseInLocation("2006-01-02", profile.MostActiveDay, beijingLoc)
	mostActiveDayEnd := mostActiveDayStart.Add(24 * time.Hour)

	hourlyFilter := bson.M{
		"client_ip": ip,
		"timestamp": bson.M{
			"$gte": mostActiveDayStart.UTC(),
			"$lt":  mostActiveDayEnd.UTC(),
		},
	}

	hourlyPipeline := []bson.M{
		{"$match": hourlyFilter},
		{"$project": bson.M{
			"hour": bson.M{
				"$hour": bson.M{
					"$add": []interface{}{"$timestamp", 8 * 60 * 60 * 1000}, // 添加8小时转换为北京时间
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

	// 填充基础信息字段
	profile.AttackCount = profile.TotalAttacks
	
	// 获取首次和最后访问时间
	firstLastPipeline := []bson.M{
		{"$match": bson.M{"client_ip": ip}},
		{"$group": bson.M{
			"_id": nil,
			"first_seen": bson.M{"$min": "$timestamp"},
			"last_seen":  bson.M{"$max": "$timestamp"},
		}},
	}
	
	firstLastCursor, err := mongoCollection.Aggregate(ctx, firstLastPipeline)
	if err == nil {
		defer firstLastCursor.Close(ctx)
		var firstLastResults []bson.M
		if firstLastCursor.All(ctx, &firstLastResults) == nil && len(firstLastResults) > 0 {
			result := firstLastResults[0]
			if firstSeen, ok := result["first_seen"].(primitive.DateTime); ok {
				profile.FirstSeen = time.Unix(int64(firstSeen)/1000, 0).Format(time.RFC3339)
			}
			if lastSeen, ok := result["last_seen"].(primitive.DateTime); ok {
				profile.LastSeen = time.Unix(int64(lastSeen)/1000, 0).Format(time.RFC3339)
			}
		}
	}
	
	// 分析攻击类型（基于日志状态和URL模式）
	profile.AttackTypes = analyzeAttackTypes(ctx, ip)
	
	// 计算威胁等级
	profile.ThreatLevel = calculateThreatLevel(profile.TotalAttacks, len(profile.AttackTypes))
	
	// 地理位置信息（简化处理，基于IP分析）
	profile.Country, profile.City, profile.ISP, profile.ASN = getIPLocationInfo(ip)

	return profile, nil
}

// analyzeAttackTypes 分析攻击类型
func analyzeAttackTypes(ctx context.Context, ip string) []AttackType {
	if mongoCollection == nil {
		return []AttackType{}
	}

	// 基于URL模式分析攻击类型
	patterns := map[string]*regexp.Regexp{
		"SQL注入":  regexp.MustCompile(`(?i)(union|select|insert|delete|drop|script|javascript|onload)`),
		"XSS攻击":  regexp.MustCompile(`(?i)(script|alert|onerror|onload|javascript)`),
		"路径遍历":   regexp.MustCompile(`(?i)(\.\.\/|\.\.\\|\/etc\/|\\windows\\)`),
		"命令注入":   regexp.MustCompile(`(?i)(;|&&|\|\||cmd|exec|system)`),
		"文件上传":   regexp.MustCompile(`(?i)(\.php|\.asp|\.jsp|upload)`),
	}

	attackTypes := []AttackType{}
	
	// 查询被拦截的请求
	filter := bson.M{
		"client_ip": ip,
		"status":    bson.M{"$ne": "success"},
	}
	
	cursor, err := mongoCollection.Find(ctx, filter)
	if err != nil {
		return attackTypes
	}
	defer cursor.Close(ctx)

	typeCounts := make(map[string]int)
	
	for cursor.Next(ctx) {
		var log bson.M
		if err := cursor.Decode(&log); err != nil {
			continue
		}
		
		url, ok := log["url"].(string)
		if !ok {
			continue
		}
		
		// 检查URL匹配的攻击类型
		for attackType, pattern := range patterns {
			if pattern.MatchString(url) {
				typeCounts[attackType]++
			}
		}
	}
	
	// 转换为结果格式
	for attackType, count := range typeCounts {
		if count > 0 {
			attackTypes = append(attackTypes, AttackType{
				Type:  attackType,
				Count: count,
			})
		}
	}
	
	return attackTypes
}

// calculateThreatLevel 计算威胁等级
func calculateThreatLevel(attackCount int, attackTypeCount int) string {
	if attackCount >= 50 || attackTypeCount >= 4 {
		return "高危"
	} else if attackCount >= 20 || attackTypeCount >= 2 {
		return "中等"
	} else if attackCount > 0 {
		return "低危"
	}
	return "未知"
}

// getIPLocationInfo 获取IP地理位置信息（简化版）
func getIPLocationInfo(ip string) (country, city, isp, asn string) {
	// 简化的IP地理位置分析
	// 在实际应用中，这里应该调用第三方IP地理位置API
	
	// 检查是否为私有IP
	if isPrivateIP(ip) {
		return "本地网络", "内网", "私有网络", "Private"
	}
	
	// 基于IP地址段的简单判断（示例）
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "未知", "未知", "未知", "未知"
	}
	
	// 简单的地理位置判断逻辑
	if strings.HasPrefix(ip, "1.") || strings.HasPrefix(ip, "14.") {
		return "中国", "北京", "中国电信", "AS4134"
	} else if strings.HasPrefix(ip, "8.8.") || strings.HasPrefix(ip, "8.4.") {
		return "美国", "加利福尼亚", "Google", "AS15169"
	} else if strings.HasPrefix(ip, "114.114.") {
		return "中国", "南京", "114DNS", "AS4837"
	}
	
	// 默认值
	return "未知", "未知", "未知", "未知"
}

// isPrivateIP 检查是否为私有IP
func isPrivateIP(ip string) bool {
	private := false
	IP := net.ParseIP(ip)
	if IP == nil {
		return false
	}
	
	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	
	private = private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	
	return private
}
