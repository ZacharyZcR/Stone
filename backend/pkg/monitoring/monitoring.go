package monitoring

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync/atomic"
	"time"
)

var metricsCollection *mongo.Collection

// In-memory counters for degraded mode
var (
	websiteRequestsTotal    int64 = 0
	blockedByBlacklistTotal int64 = 0
	blockedByRulesTotal     int64 = 0
)

func SetMongoCollection(collection *mongo.Collection) {
	metricsCollection = collection
}

type DailyMetrics struct {
	Date                    time.Time `bson:"date"`
	WebsiteRequestsTotal    int       `bson:"websiteRequestsTotal"`
	BlockedByBlacklistTotal int       `bson:"blockedByBlacklistTotal"`
	BlockedByRulesTotal     int       `bson:"blockedByRulesTotal"`
}

// InMemoryMetrics represents current session metrics
type InMemoryMetrics struct {
	WebsiteRequestsTotal    int64 `json:"website_requests_total"`
	BlockedByBlacklistTotal int64 `json:"blocked_by_blacklist_total"`
	BlockedByRulesTotal     int64 `json:"blocked_by_rules_total"`
}

func IncrementMetric(metric string) error {
	// Always increment in-memory counters
	switch metric {
	case "websiteRequestsTotal":
		atomic.AddInt64(&websiteRequestsTotal, 1)
	case "blockedByBlacklistTotal":
		atomic.AddInt64(&blockedByBlacklistTotal, 1)
	case "blockedByRulesTotal":
		atomic.AddInt64(&blockedByRulesTotal, 1)
	}

	// Try to persist to MongoDB if available
	if metricsCollection == nil {
		// MongoDB not initialized - only use in-memory counters
		return nil
	}

	// 定义北京时区
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// If timezone fails, only increment in-memory
		log.Printf("Failed to load Beijing timezone, using in-memory counters: %v", err)
		return nil
	}

	// 获取当前北京时间（去掉时分秒）
	now := time.Now().In(beijingLocation)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, beijingLocation)

	// 将北京时间转换为 UTC 时间存储
	todayUTC := today.UTC()

	filter := bson.M{"date": todayUTC}
	update := bson.M{
		"$inc":         bson.M{metric: 1},
		"$setOnInsert": bson.M{"date": todayUTC},
	}
	opts := options.Update().SetUpsert(true)

	result, err := metricsCollection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		// MongoDB failed, but we still have in-memory counters
		log.Printf("MongoDB metric update failed (using in-memory): %v", err)
		return nil
	}

	log.Printf("Metric %s updated for date %v (UTC: %v). Matched: %d, Modified: %d, Upserted: %d",
		metric, today, todayUTC, result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	return nil
}

// GetInMemoryMetrics returns current session metrics
func GetInMemoryMetrics() InMemoryMetrics {
	return InMemoryMetrics{
		WebsiteRequestsTotal:    atomic.LoadInt64(&websiteRequestsTotal),
		BlockedByBlacklistTotal: atomic.LoadInt64(&blockedByBlacklistTotal),
		BlockedByRulesTotal:     atomic.LoadInt64(&blockedByRulesTotal),
	}
}
