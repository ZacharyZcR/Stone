package handlers

import (
	"Stone/pkg/monitoring"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"runtime"
	"time"
)

// GetStatus 获取系统状态
func GetStatus(c *gin.Context) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取CPU使用率"})
		return
	}

	// 获取内存使用率
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取内存使用率"})
		return
	}

	// 获取磁盘使用率
	diskStat, err := disk.Usage("/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取磁盘使用率"})
		return
	}

	// 获取网络流量
	netIO, err := net.IOCounters(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取网络流量"})
		return
	}

	// 获取系统负载
	loadStat, err := load.Avg()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取系统负载"})
		return
	}

	// 获取系统启动时间
	uptime := time.Since(monitoring.StartTime).String()

	// 获取打开文件描述符数量
	numFDs := runtime.NumGoroutine() // 示例，实际需要更复杂的操作

	// 获取线程和进程数量
	numThreads := runtime.NumGoroutine()   // 示例，实际需要更复杂的操作
	numProcesses := runtime.NumGoroutine() // 示例，实际需要更复杂的操作

	c.JSON(http.StatusOK, gin.H{
		"status":            "running",
		"uptime":            uptime,
		"cpu_usage_percent": cpuPercent[0],
		"memory_usage":      vmStat.UsedPercent,
		"disk_usage":        diskStat.UsedPercent,
		"network_in":        netIO[0].BytesRecv,
		"network_out":       netIO[0].BytesSent,
		"load_average":      loadStat.Load1,
		"open_file_desc":    numFDs,
		"threads":           numThreads,
		"processes":         numProcesses,
	})
}

var metricsCollection *mongo.Collection

// SetMetricsCollection 设置 metrics 集合
func SetMetricsCollection(collection *mongo.Collection) {
	metricsCollection = collection
}

// GetFirewallMetrics 获取防火墙相关的指标
func GetFirewallMetrics(c *gin.Context) {
	if metricsCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Metrics collection is not initialized"})
		return
	}

	var metrics struct {
		Timestamp               time.Time `bson:"timestamp"`
		WebsiteRequestsTotal    int       `bson:"websiteRequestsTotal"`
		BlockedByBlacklistTotal int       `bson:"blockedByBlacklistTotal"`
		BlockedByRulesTotal     int       `bson:"blockedByRulesTotal"`
	}

	err := metricsCollection.FindOne(
		context.Background(),
		bson.M{},
		options.FindOne().SetSort(bson.D{{"timestamp", -1}}),
	).Decode(&metrics)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No metrics data found"})
		} else {
			log.Printf("Error retrieving metrics: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics data"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success_requests":   metrics.WebsiteRequestsTotal,
		"blacklist_requests": metrics.BlockedByBlacklistTotal,
		"rules_requests":     metrics.BlockedByRulesTotal,
		"last_updated":       metrics.Timestamp,
	})
}
