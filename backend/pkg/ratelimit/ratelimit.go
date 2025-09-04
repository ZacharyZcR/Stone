// pkg/ratelimit/ratelimit.go

package ratelimit

import (
	"context"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RateLimitRule 速率限制规则
type RateLimitRule struct {
	ID          int    `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Enabled     bool   `bson:"enabled" json:"enabled"`
	
	// 令牌桶参数
	Capacity    int           `bson:"capacity" json:"capacity"`       // 桶容量(最大令牌数)
	RefillRate  int           `bson:"refill_rate" json:"refill_rate"` // 每秒补充令牌数
	Window      time.Duration `bson:"window" json:"window"`           // 时间窗口(秒)
	
	// 处理策略
	Action      string `bson:"action" json:"action"`             // block/delay/warn
	BlockTime   int    `bson:"block_time" json:"block_time"`     // 阻断时间(秒)
	
	CreatedAt   string `bson:"created_at" json:"created_at"`
	TriggeredCount int `bson:"triggered_count" json:"triggered_count"`
}

// TokenBucket 令牌桶
type TokenBucket struct {
	tokens      int       // 当前令牌数
	capacity    int       // 桶容量
	refillRate  int       // 每秒补充速率
	lastRefill  time.Time // 上次补充时间
	mu          sync.Mutex
}

// RateLimiter 速率限制器
type RateLimiter struct {
	buckets     sync.Map  // string -> *TokenBucket
	rules       []RateLimitRule
	rulesMutex  sync.RWMutex
	mongoCollection *mongo.Collection
	
	// LRU清理
	cleanup     *time.Ticker
	maxBuckets  int
}

// DefaultRule 默认规则
var DefaultRule = RateLimitRule{
	ID:          0,
	Name:        "default",
	Description: "Default rate limiting rule",
	Enabled:     true,
	Capacity:    60,           // 60个令牌
	RefillRate:  10,           // 每秒补充10个
	Window:      60 * time.Second, // 60秒窗口
	Action:      "block",
	BlockTime:   300,          // 阻断5分钟
}

// globalLimiter 全局实例
var globalLimiter *RateLimiter

// Init 初始化速率限制器
func Init(collection *mongo.Collection) *RateLimiter {
	rl := &RateLimiter{
		mongoCollection: collection,
		cleanup:        time.NewTicker(5 * time.Minute),
		maxBuckets:     10000, // 最多1万个IP的令牌桶
	}
	
	// 加载规则
	rl.loadRules()
	
	// 启动清理协程
	go rl.cleanupBuckets()
	
	globalLimiter = rl
	return rl
}

// GetGlobalLimiter 获取全局限制器
func GetGlobalLimiter() *RateLimiter {
	return globalLimiter
}

// loadRules 从MongoDB加载规则
func (rl *RateLimiter) loadRules() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var result struct {
		Rules []RateLimitRule `bson:"rules"`
	}
	
	err := rl.mongoCollection.FindOne(ctx, bson.M{"type": "rate_limit"}).Decode(&result)
	if err != nil {
		// 使用默认规则
		rl.rulesMutex.Lock()
		rl.rules = []RateLimitRule{DefaultRule}
		rl.rulesMutex.Unlock()
		return
	}
	
	rl.rulesMutex.Lock()
	rl.rules = result.Rules
	if len(rl.rules) == 0 {
		rl.rules = []RateLimitRule{DefaultRule}
	}
	rl.rulesMutex.Unlock()
}

// newTokenBucket 创建令牌桶
func newTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity, // 初始满桶
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// refill 补充令牌
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	
	// 计算应该补充的令牌数
	tokensToAdd := int(elapsed * float64(tb.refillRate))
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// consume 消费令牌
func (tb *TokenBucket) consume(tokens int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	tb.refill()
	
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

// IsAllowed 检查IP是否被允许(主入口)
func IsAllowed(ip string) (allowed bool, action string) {
	if globalLimiter == nil {
		return true, ""
	}
	return globalLimiter.checkIP(ip)
}

// checkIP 检查单个IP
func (rl *RateLimiter) checkIP(ip string) (bool, string) {
	rl.rulesMutex.RLock()
	if len(rl.rules) == 0 {
		rl.rulesMutex.RUnlock()
		return true, ""
	}
	
	// 使用第一个启用的规则(简化逻辑)
	var rule RateLimitRule
	found := false
	for _, r := range rl.rules {
		if r.Enabled {
			rule = r
			found = true
			break
		}
	}
	rl.rulesMutex.RUnlock()
	
	if !found {
		return true, ""
	}
	
	// 获取或创建令牌桶
	bucket := rl.getBucket(ip, rule.Capacity, rule.RefillRate)
	
	// 尝试消费1个令牌
	if bucket.consume(1) {
		return true, ""
	}
	
	// 令牌不足，根据策略处理
	return false, rule.Action
}

// getBucket 获取或创建IP的令牌桶
func (rl *RateLimiter) getBucket(ip string, capacity, refillRate int) *TokenBucket {
	if bucket, ok := rl.buckets.Load(ip); ok {
		return bucket.(*TokenBucket)
	}
	
	// 创建新桶
	newBucket := newTokenBucket(capacity, refillRate)
	actual, loaded := rl.buckets.LoadOrStore(ip, newBucket)
	if loaded {
		return actual.(*TokenBucket)
	}
	
	return newBucket
}

// cleanupBuckets 清理过期的令牌桶
func (rl *RateLimiter) cleanupBuckets() {
	for range rl.cleanup.C {
		bucketCount := 0
		cutoff := time.Now().Add(-10 * time.Minute) // 10分钟未使用就删除
		
		rl.buckets.Range(func(key, value interface{}) bool {
			bucket := value.(*TokenBucket)
			bucket.mu.Lock()
			isOld := bucket.lastRefill.Before(cutoff)
			bucket.mu.Unlock()
			
			if isOld {
				rl.buckets.Delete(key)
			} else {
				bucketCount++
			}
			
			// 如果桶数量过多，删除一些旧的
			if bucketCount > rl.maxBuckets {
				rl.buckets.Delete(key)
				bucketCount--
			}
			
			return true
		})
	}
}

// GetRules 获取所有规则
func (rl *RateLimiter) GetRules() []RateLimitRule {
	rl.rulesMutex.RLock()
	defer rl.rulesMutex.RUnlock()
	
	rules := make([]RateLimitRule, len(rl.rules))
	copy(rules, rl.rules)
	return rules
}

// AddRule 添加规则
func (rl *RateLimiter) AddRule(rule RateLimitRule) error {
	rl.rulesMutex.Lock()
	rl.rules = append(rl.rules, rule)
	rl.rulesMutex.Unlock()
	
	return rl.saveRules()
}

// UpdateRule 更新规则
func (rl *RateLimiter) UpdateRule(id int, rule RateLimitRule) error {
	rl.rulesMutex.Lock()
	for i, r := range rl.rules {
		if r.ID == id {
			rl.rules[i] = rule
			break
		}
	}
	rl.rulesMutex.Unlock()
	
	return rl.saveRules()
}

// DeleteRule 删除规则
func (rl *RateLimiter) DeleteRule(id int) error {
	rl.rulesMutex.Lock()
	for i, r := range rl.rules {
		if r.ID == id {
			rl.rules = append(rl.rules[:i], rl.rules[i+1:]...)
			break
		}
	}
	rl.rulesMutex.Unlock()
	
	return rl.saveRules()
}

// saveRules 保存规则到MongoDB
func (rl *RateLimiter) saveRules() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := rl.mongoCollection.UpdateOne(
		ctx,
		bson.M{"type": "rate_limit"},
		bson.M{
			"$set": bson.M{
				"rules": rl.rules,
				"updated_at": time.Now().Format("2006-01-02 15:04:05"),
			},
		},
		nil,
	)
	
	return err
}

// GetStats 获取统计信息
func (rl *RateLimiter) GetStats() map[string]interface{} {
	bucketCount := 0
	rl.buckets.Range(func(key, value interface{}) bool {
		bucketCount++
		return true
	})
	
	return map[string]interface{}{
		"active_buckets": bucketCount,
		"rules_count":   len(rl.rules),
		"max_buckets":   rl.maxBuckets,
	}
}