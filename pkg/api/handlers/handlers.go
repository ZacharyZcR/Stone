// pkg/api/handlers/handlers.go

package handlers

import (
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
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

// GetRules 获取当前规则
func GetRules(c *gin.Context) {
	// 获取当前规则
	currentRules := rules.GetCurrentRules()
	c.JSON(http.StatusOK, currentRules)
}

// UpdateRules 更新规则
func UpdateRules(c *gin.Context) {
	var newRules rules.Rules
	if err := c.ShouldBindJSON(&newRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新全局规则
	rules.UpdateRules(newRules.Whitelist, newRules.Blacklist)

	c.JSON(http.StatusOK, gin.H{"status": "rules updated"})
}
