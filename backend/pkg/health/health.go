// pkg/health/health.go

package health

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync/atomic"
	"time"
)

// HealthStatus represents the health state of a service
type HealthStatus struct {
	IsHealthy   bool      `json:"is_healthy"`
	LastChecked time.Time `json:"last_checked"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
}

// SystemHealth represents overall system health
type SystemHealth struct {
	MongoDB     HealthStatus `json:"mongodb"`
	RulesEngine HealthStatus `json:"rules_engine"`
	IsHealthy   bool         `json:"is_healthy"`
	DegradedMode bool        `json:"degraded_mode"`
}

var (
	// Use atomic values for thread-safe access
	mongoHealthy        int64 = 1 // 1 = healthy, 0 = unhealthy
	rulesEngineHealthy  int64 = 1
	systemInDegradedMode int64 = 0 // 1 = degraded, 0 = normal

	mongoCollection     *mongo.Collection
	lastHealthCheck     time.Time
	healthCheckInterval = 30 * time.Second
)

// SetMongoCollection sets the MongoDB collection for health checks
func SetMongoCollection(collection *mongo.Collection) {
	mongoCollection = collection
}

// IsMongoHealthy returns true if MongoDB is healthy
func IsMongoHealthy() bool {
	return atomic.LoadInt64(&mongoHealthy) == 1
}

// IsRulesEngineHealthy returns true if rules engine is healthy
func IsRulesEngineHealthy() bool {
	return atomic.LoadInt64(&rulesEngineHealthy) == 1
}

// IsSystemInDegradedMode returns true if system is in degraded mode
func IsSystemInDegradedMode() bool {
	return atomic.LoadInt64(&systemInDegradedMode) == 1
}

// checkMongoHealth checks MongoDB connection
func checkMongoHealth() HealthStatus {
	if mongoCollection == nil {
		atomic.StoreInt64(&mongoHealthy, 0)
		return HealthStatus{
			IsHealthy:   false,
			LastChecked: time.Now(),
			ErrorMsg:    "MongoDB collection not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Simple ping to check connection
	err := mongoCollection.Database().Client().Ping(ctx, nil)
	if err != nil {
		atomic.StoreInt64(&mongoHealthy, 0)
		return HealthStatus{
			IsHealthy:   false,
			LastChecked: time.Now(),
			ErrorMsg:    err.Error(),
		}
	}

	atomic.StoreInt64(&mongoHealthy, 1)
	return HealthStatus{
		IsHealthy:   true,
		LastChecked: time.Now(),
	}
}

// checkRulesEngineHealth checks if rules engine is functional
func checkRulesEngineHealth() HealthStatus {
	// For now, rules engine is healthy if MongoDB is healthy
	// since rules are stored in MongoDB
	mongoHealthy := IsMongoHealthy()
	
	if mongoHealthy {
		atomic.StoreInt64(&rulesEngineHealthy, 1)
	} else {
		atomic.StoreInt64(&rulesEngineHealthy, 0)
	}

	return HealthStatus{
		IsHealthy:   mongoHealthy,
		LastChecked: time.Now(),
		ErrorMsg:    func() string {
			if !mongoHealthy {
				return "Rules engine depends on MongoDB"
			}
			return ""
		}(),
	}
}

// CheckSystemHealth performs a comprehensive health check
func CheckSystemHealth() SystemHealth {
	mongoStatus := checkMongoHealth()
	rulesStatus := checkRulesEngineHealth()
	
	isHealthy := mongoStatus.IsHealthy && rulesStatus.IsHealthy
	degradedMode := !isHealthy
	
	// Update degraded mode flag
	if degradedMode {
		atomic.StoreInt64(&systemInDegradedMode, 1)
		log.Printf("System entering degraded mode - MongoDB: %v, Rules: %v", 
			mongoStatus.IsHealthy, rulesStatus.IsHealthy)
	} else {
		if atomic.LoadInt64(&systemInDegradedMode) == 1 {
			log.Printf("System exiting degraded mode - all services healthy")
		}
		atomic.StoreInt64(&systemInDegradedMode, 0)
	}
	
	lastHealthCheck = time.Now()
	
	return SystemHealth{
		MongoDB:      mongoStatus,
		RulesEngine:  rulesStatus,
		IsHealthy:    isHealthy,
		DegradedMode: degradedMode,
	}
}

// StartHealthMonitor starts a background health monitoring routine
func StartHealthMonitor() {
	go func() {
		ticker := time.NewTicker(healthCheckInterval)
		defer ticker.Stop()
		
		// Initial health check
		CheckSystemHealth()
		
		for range ticker.C {
			CheckSystemHealth()
		}
	}()
	
	log.Printf("Health monitor started, checking every %v", healthCheckInterval)
}

// GetLastHealthCheck returns the timestamp of the last health check
func GetLastHealthCheck() time.Time {
	return lastHealthCheck
}