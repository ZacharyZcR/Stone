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

type Metrics struct {
	Timestamp               time.Time `bson:"timestamp"`
	WebsiteRequestsTotal    int       `bson:"websiteRequestsTotal"`
	BlockedByBlacklistTotal int       `bson:"blockedByBlacklistTotal"`
	BlockedByRulesTotal     int       `bson:"blockedByRulesTotal"`
}

func IncrementMetric(metric string) error {
	if metricsCollection == nil {
		return fmt.Errorf("metrics collection is not initialized")
	}

	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{metric: 1},
		"$set": bson.M{"timestamp": time.Now()},
	}
	opts := options.Update().SetUpsert(true)

	result, err := metricsCollection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error updating metric %s: %v", metric, err)
		return err
	}

	log.Printf("Metric %s updated. Matched: %d, Modified: %d, Upserted: %d",
		metric, result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	return nil
}

func GetMetrics() (Metrics, error) {
	var metrics Metrics
	err := metricsCollection.FindOne(context.Background(), bson.M{}).Decode(&metrics)
	return metrics, err
}
