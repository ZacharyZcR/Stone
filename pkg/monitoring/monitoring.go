package monitoring

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var metricsCollection *mongo.Collection

func SetMongoCollection(collection *mongo.Collection) {
	metricsCollection = collection
}

type DailyMetrics struct {
	Date                    time.Time `bson:"date"`
	WebsiteRequestsTotal    int       `bson:"websiteRequestsTotal"`
	BlockedByBlacklistTotal int       `bson:"blockedByBlacklistTotal"`
	BlockedByRulesTotal     int       `bson:"blockedByRulesTotal"`
}

func IncrementMetric(metric string) error {
	if metricsCollection == nil {
		return fmt.Errorf("metrics collection is not initialized")
	}

	// 定义北京时区
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("failed to load Beijing timezone: %v", err)
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
		log.Printf("Error updating metric %s for date %v (UTC: %v): %v", metric, today, todayUTC, err)
		return err
	}

	log.Printf("Metric %s updated for date %v (UTC: %v). Matched: %d, Modified: %d, Upserted: %d",
		metric, today, todayUTC, result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	return nil
}
