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

type DailyMetrics struct {
	Date                    time.Time `bson:"date"`
	WebsiteRequestsTotal    int       `bson:"websiteRequestsTotal"`
	BlockedByBlacklistTotal int       `bson:"blockedByBlacklistTotal"`
	BlockedByRulesTotal     int       `bson:"blockedByRulesTotal"`
}

// GetFirewallMetrics 获取防火墙相关的指标
func GetFirewallMetrics(c *gin.Context) {
	if metricsCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Metrics collection is not initialized"})
		return
	}

	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	// 如果没有提供参数，使用当天的日期
	if startDateStr == "" || endDateStr == "" {
		startDate = time.Now().UTC().Truncate(24 * time.Hour)
		endDate = startDate
	} else {
		// 解析日期字符串
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	// 查询指定日期范围内的指标
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate.Add(24*time.Hour - time.Second), // 包括结束日期的全天
		},
	}
	opts := options.Find().SetSort(bson.D{{"date", 1}})

	cursor, err := metricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Printf("Error retrieving metrics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics data"})
		return
	}
	defer cursor.Close(context.Background())

	var metrics []DailyMetrics
	if err = cursor.All(context.Background(), &metrics); err != nil {
		log.Printf("Error decoding metrics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode metrics data"})
		return
	}

	if len(metrics) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No metrics data found for the specified date range"})
		return
	}

	// 构造响应
	response := make([]gin.H, len(metrics))
	for i, m := range metrics {
		response[i] = gin.H{
			"date":               m.Date.Format("2006-01-02"),
			"success_requests":   m.WebsiteRequestsTotal,
			"blacklist_requests": m.BlockedByBlacklistTotal,
			"rules_requests":     m.BlockedByRulesTotal,
		}
	}

	c.JSON(http.StatusOK, response)
}
