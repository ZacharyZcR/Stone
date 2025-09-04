package handlers

import (
	"Stone/backend/pkg/health"
	"Stone/backend/pkg/monitoring"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	// 获取真实的打开文件描述符数量
	numFDs := getFileDescriptorCount()

	// 获取真实的线程数量
	numThreads := getThreadCount()

	// 获取真实的进程数量
	numProcesses := getProcessCount()

	// 获取健康检查状态
	systemHealth := health.CheckSystemHealth()
	
	// 获取内存指标
	inMemoryMetrics := monitoring.GetInMemoryMetrics()

	status := "running"
	if systemHealth.DegradedMode {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            status,
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
		
		// Health information
		"health": systemHealth,
		
		// Current session metrics (always available)
		"session_metrics": inMemoryMetrics,
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

func GetFirewallMetrics(c *gin.Context) {
	if metricsCollection == nil {
		// In degraded mode, return in-memory metrics only
		log.Printf("MongoDB not available, returning in-memory metrics")
		inMemoryMetrics := monitoring.GetInMemoryMetrics()
		
		// Return a single day entry with current session data
		today := time.Now().Format("2006-01-02")
		response := []gin.H{
			{
				"date":               today,
				"success_requests":   inMemoryMetrics.WebsiteRequestsTotal,
				"blacklist_requests": inMemoryMetrics.BlockedByBlacklistTotal,
				"rules_requests":     inMemoryMetrics.BlockedByRulesTotal,
			},
		}
		
		c.JSON(http.StatusOK, response)
		return
	}

	// 定义北京时区
	beijingLocation, _ := time.LoadLocation("Asia/Shanghai")

	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	// 如果没有提供参数，使用当天的北京日期
	if startDateStr == "" || endDateStr == "" {
		now := time.Now().In(beijingLocation)
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, beijingLocation)
		endDate = startDate
	} else {
		// 解析日期字符串为北京时间
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, beijingLocation)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, beijingLocation)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	// 计算请求的天数
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	// 调整查询时间范围
	startDateUTC := startDate.UTC()
	endDateUTC := endDate.Add(24 * time.Hour).UTC()

	// 查询指定日期范围内的指标
	filter := bson.M{
		"date": bson.M{
			"$gte": startDateUTC,
			"$lt":  endDateUTC,
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

	// 构造响应，确保返回的数组长度等于请求的天数
	response := make([]gin.H, days)
	metricsMap := make(map[string]DailyMetrics)

	// 将查询到的指标数据放入 map 中，使用北京时间的日期作为键
	for _, m := range metrics {
		dateStr := m.Date.In(beijingLocation).Format("2006-01-02")
		metricsMap[dateStr] = m
	}

	// 填充响应数组
	for i := 0; i < days; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		dateStr := currentDate.Format("2006-01-02")

		if m, exists := metricsMap[dateStr]; exists {
			response[i] = gin.H{
				"date":               dateStr,
				"success_requests":   m.WebsiteRequestsTotal,
				"blacklist_requests": m.BlockedByBlacklistTotal,
				"rules_requests":     m.BlockedByRulesTotal,
			}
		} else {
			response[i] = gin.H{
				"date":               dateStr,
				"success_requests":   0,
				"blacklist_requests": 0,
				"rules_requests":     0,
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// getFileDescriptorCount 获取当前进程的文件描述符数量
func getFileDescriptorCount() int {
	pid := os.Getpid()
	fdDir := "/proc/" + strconv.Itoa(pid) + "/fd"
	
	files, err := ioutil.ReadDir(fdDir)
	if err != nil {
		log.Printf("Error reading file descriptors: %v", err)
		return 0
	}
	
	return len(files)
}

// getThreadCount 获取当前进程的线程数量
func getThreadCount() int {
	pid := os.Getpid()
	statusFile := "/proc/" + strconv.Itoa(pid) + "/status"
	
	content, err := ioutil.ReadFile(statusFile)
	if err != nil {
		log.Printf("Error reading process status: %v", err)
		return 0
	}
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Threads:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if count, err := strconv.Atoi(parts[1]); err == nil {
					return count
				}
			}
		}
	}
	
	return 0
}

// getProcessCount 获取系统进程总数
func getProcessCount() int {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error getting process list: %v", err)
		return 0
	}
	
	return len(processes)
}
