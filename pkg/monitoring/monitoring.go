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

	// 获取当前日期（去掉时分秒）
	today := time.Now().Truncate(24 * time.Hour)

	filter := bson.M{"date": today}
	update := bson.M{
		"$inc":         bson.M{metric: 1},
		"$setOnInsert": bson.M{"date": today},
	}
	opts := options.Update().SetUpsert(true)

	result, err := metricsCollection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error updating metric %s for date %v: %v", metric, today, err)
		return err
	}

	log.Printf("Metric %s updated for date %v. Matched: %d, Modified: %d, Upserted: %d",
		metric, today, result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	return nil
}

func GetMetrics(days int) ([]DailyMetrics, error) {
	if metricsCollection == nil {
		return nil, fmt.Errorf("metrics collection is not initialized")
	}

	// 计算开始日期
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -days+1)

	// 查询指定日期范围内的指标
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	opts := options.Find().SetSort(bson.D{{"date", 1}})

	cursor, err := metricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var metrics []DailyMetrics
	if err = cursor.All(context.Background(), &metrics); err != nil {
		return nil, err
	}

	// 如果某天没有数据，填充零值
	filledMetrics := fillMissingDays(metrics, startDate, endDate)

	return filledMetrics, nil
}

func fillMissingDays(metrics []DailyMetrics, startDate, endDate time.Time) []DailyMetrics {
	filledMetrics := make([]DailyMetrics, 0)
	metricsMap := make(map[time.Time]DailyMetrics)

	for _, m := range metrics {
		metricsMap[m.Date] = m
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if m, exists := metricsMap[d]; exists {
			filledMetrics = append(filledMetrics, m)
		} else {
			filledMetrics = append(filledMetrics, DailyMetrics{Date: d})
		}
	}

	return filledMetrics
}
