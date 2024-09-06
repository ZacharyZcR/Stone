package handlers

import (
	"Stone/pkg/logging"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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
