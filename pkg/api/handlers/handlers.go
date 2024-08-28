// pkg/api/handlers/handlers.go

package handlers

import (
	"Stone/pkg/logging"
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"runtime"
	"time"
)

// GetStatus 获取系统状态
func GetStatus(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 获取系统启动时间
	uptime := time.Since(monitoring.StartTime).String()

	// 获取处理的请求总数（假设使用Prometheus）
	requestsTotal, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取请求总数"})
		return
	}

	var totalRequests float64
	for _, mf := range requestsTotal {
		if mf.GetName() == "stone_requests_total" {
			for _, metric := range mf.GetMetric() {
				totalRequests += metric.GetCounter().GetValue()
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "running",
		"uptime":           uptime,
		"goroutines":       runtime.NumGoroutine(),
		"memory_alloc":     memStats.Alloc,
		"total_requests":   totalRequests,
		"cpu_cores":        runtime.NumCPU(),
		"cpu_architecture": runtime.GOARCH,
		"go_version":       runtime.Version(),
	})
}

// GetIPControlRules 获取当前IP控制规则
func GetIPControlRules(c *gin.Context) {
	currentRules := rules.GetIPControlRules()
	c.JSON(http.StatusOK, currentRules)
}

// GetInterceptionRules 获取当前拦截规则
func GetInterceptionRules(c *gin.Context) {
	currentRules := rules.GetInterceptionRules()
	c.JSON(http.StatusOK, currentRules)
}

// UpdateIPControlRules 更新IP控制规则
func UpdateIPControlRules(c *gin.Context) {
	var newRules rules.IPControlRules
	if err := c.ShouldBindJSON(&newRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules.UpdateIPControlRules(newRules.Whitelist, newRules.Blacklist)
	c.JSON(http.StatusOK, gin.H{"status": "IP control rules updated"})
}

// UpdateInterceptionRules 更新拦截规则
func UpdateInterceptionRules(c *gin.Context) {
	var newRules struct {
		URLPatterns  []rules.Pattern `json:"url_patterns"`
		BodyPatterns []rules.Pattern `json:"body_patterns"`
	}

	if err := c.ShouldBindJSON(&newRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules.UpdateInterceptionRules(newRules.URLPatterns, newRules.BodyPatterns)
	c.JSON(http.StatusOK, gin.H{"status": "Interception rules updated"})
}

// GetLogs 查看日志
func GetLogs(c *gin.Context) {
	// 解析查询参数
	startDateTimeStr := c.Query("startDateTime")
	endDateTimeStr := c.Query("endDateTime")
	ip := c.Query("ip")

	// 解析时间参数
	var startDateTime, endDateTime time.Time
	var err error
	if startDateTimeStr != "" {
		startDateTime, err = time.Parse(time.RFC3339, startDateTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式，必须为RFC3339格式"})
			return
		}
	}
	if endDateTimeStr != "" {
		endDateTime, err = time.Parse(time.RFC3339, endDateTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式，必须为RFC3339格式"})
			return
		}
	}

	// 获取日志
	logs, err := logging.FetchLogsFromMongoWithFilters(context.Background(), 100, startDateTime, endDateTime, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取日志"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
