// pkg/api/handlers/rate_limit.go

package handlers

import (
	"Stone/backend/pkg/ratelimit"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetRateLimitRules 获取所有速率限制规则
func GetRateLimitRules(c *gin.Context) {
	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	rules := limiter.GetRules()
	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
		"total": len(rules),
	})
}

// GetRateLimitRule 获取单个规则
func GetRateLimitRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid rule ID",
		})
		return
	}

	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	rules := limiter.GetRules()
	for _, rule := range rules {
		if rule.ID == id {
			c.JSON(http.StatusOK, rule)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Rule not found",
	})
}

// CreateRateLimitRule 创建新的速率限制规则
func CreateRateLimitRule(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		Capacity    int    `json:"capacity" binding:"required,min=1"`
		RefillRate  int    `json:"refill_rate" binding:"required,min=1"`
		Window      int    `json:"window" binding:"required,min=1"` // 秒数
		Action      string `json:"action" binding:"required,oneof=block delay warn"`
		BlockTime   int    `json:"block_time" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	// 生成新ID
	rules := limiter.GetRules()
	maxID := 0
	for _, rule := range rules {
		if rule.ID > maxID {
			maxID = rule.ID
		}
	}

	newRule := ratelimit.RateLimitRule{
		ID:          maxID + 1,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Capacity:    req.Capacity,
		RefillRate:  req.RefillRate,
		Window:      time.Duration(req.Window) * time.Second,
		Action:      req.Action,
		BlockTime:   req.BlockTime,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		TriggeredCount: 0,
	}

	if err := limiter.AddRule(newRule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rule created successfully",
		"rule":    newRule,
	})
}

// UpdateRateLimitRule 更新速率限制规则
func UpdateRateLimitRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid rule ID",
		})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		Capacity    int    `json:"capacity" binding:"required,min=1"`
		RefillRate  int    `json:"refill_rate" binding:"required,min=1"`
		Window      int    `json:"window" binding:"required,min=1"` // 秒数
		Action      string `json:"action" binding:"required,oneof=block delay warn"`
		BlockTime   int    `json:"block_time" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	// 检查规则是否存在
	rules := limiter.GetRules()
	var existingRule *ratelimit.RateLimitRule
	for _, rule := range rules {
		if rule.ID == id {
			existingRule = &rule
			break
		}
	}

	if existingRule == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Rule not found",
		})
		return
	}

	updatedRule := ratelimit.RateLimitRule{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		Enabled:        req.Enabled,
		Capacity:       req.Capacity,
		RefillRate:     req.RefillRate,
		Window:         time.Duration(req.Window) * time.Second,
		Action:         req.Action,
		BlockTime:      req.BlockTime,
		CreatedAt:      existingRule.CreatedAt,
		TriggeredCount: existingRule.TriggeredCount,
	}

	if err := limiter.UpdateRule(id, updatedRule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rule updated successfully",
		"rule":    updatedRule,
	})
}

// DeleteRateLimitRule 删除速率限制规则
func DeleteRateLimitRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid rule ID",
		})
		return
	}

	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	// 检查规则是否存在
	rules := limiter.GetRules()
	found := false
	for _, rule := range rules {
		if rule.ID == id {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Rule not found",
		})
		return
	}

	if err := limiter.DeleteRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rule deleted successfully",
	})
}

// GetRateLimitStats 获取速率限制统计信息
func GetRateLimitStats(c *gin.Context) {
	limiter := ratelimit.GetGlobalLimiter()
	if limiter == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Rate limiter not initialized",
		})
		return
	}

	stats := limiter.GetStats()
	c.JSON(http.StatusOK, stats)
}

// ResetRateLimitCounters 重置速率限制计数器(管理功能)
func ResetRateLimitCounters(c *gin.Context) {
	// 这个功能通过重新初始化limiter实现
	// 在生产环境中可能需要更精细的控制
	c.JSON(http.StatusOK, gin.H{
		"message": "Rate limit counters reset functionality not implemented",
		"note": "Counters will automatically reset after cleanup period",
	})
}