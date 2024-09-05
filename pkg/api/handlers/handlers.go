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

// HandleIPControlRules 处理IP控制规则的操作
func HandleIPControlRules(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		ip := c.Param("ip")
		if ip == "" {
			// 获取所有IP控制规则
			currentRules := rules.GetIPControlRules()
			c.JSON(http.StatusOK, currentRules)
		} else {
			// 获取特定IP的规则
			rule, found := rules.GetIPRule(ip)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
				return
			}
			c.JSON(http.StatusOK, rule)
		}
	case http.MethodPost:
		var newRule rules.IPControlRule
		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := rules.AddIPRule(newRule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule added"})
	case http.MethodDelete:
		ip := c.Param("ip")
		if err := rules.DeleteIPRule(ip); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule deleted"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

// HandleInterceptionRules 处理拦截规则的操作
func HandleInterceptionRules(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		name := c.Param("name")
		if name == "" {
			// 获取所有拦截规则
			currentRules := rules.GetInterceptionRules()
			c.JSON(http.StatusOK, currentRules)
		} else {
			// 获取特定名称的规则
			rule, found := rules.GetInterceptionRule(name)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
				return
			}
			c.JSON(http.StatusOK, rule)
		}
	case http.MethodPost:
		var newRule rules.Pattern
		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := rules.AddInterceptionRule(newRule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add interception rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Interception rule added"})
	case http.MethodDelete:
		name := c.Param("name")
		if err := rules.DeleteInterceptionRule(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interception rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Interception rule deleted"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}
