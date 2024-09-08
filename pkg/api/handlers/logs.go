package handlers

import (
	"Stone/pkg/logging"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetLogs 查看日志
func GetLogs(c *gin.Context) {
	// 解析查询参数
	startDateTimeStr := c.Query("startDateTime")
	endDateTimeStr := c.Query("endDateTime")
	ip := c.Query("ip")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

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

	// 验证状态参数
	if status != "" && status != "blocked" && status != "passed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态参数，必须为 'blocked' 或 'passed'"})
		return
	}

	// 获取日志
	logs, totalCount, err := logging.FetchLogsFromMongoWithFilters(context.Background(), page, pageSize, startDateTime, endDateTime, ip, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取日志"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":       logs,
		"totalCount": totalCount,
		"page":       page,
		"pageSize":   pageSize,
	})
}
