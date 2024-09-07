package handlers

import (
	"Stone/pkg/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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

// GetFirewallMetrics 获取防火墙相关的指标
func GetFirewallMetrics(c *gin.Context) {
	// 获取所有指标
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取指标数据"})
		return
	}

	var successCount, blacklistCount, rulesCount float64

	// 遍历指标，查找相关的计数器
	for _, mf := range metrics {
		switch mf.GetName() {
		case "website_requests_total":
			for _, metric := range mf.GetMetric() {
				for _, label := range metric.GetLabel() {
					if label.GetName() == "status" && label.GetValue() == "success" {
						successCount += metric.GetCounter().GetValue()
					}
				}
			}
		case "blocked_by_blacklist_total":
			for _, metric := range mf.GetMetric() {
				blacklistCount += metric.GetCounter().GetValue()
			}
		case "blocked_by_rules_total":
			for _, metric := range mf.GetMetric() {
				rulesCount += metric.GetCounter().GetValue()
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success_requests":   successCount,
		"blacklist_requests": blacklistCount,
		"rules_requests":     rulesCount,
	})
}
