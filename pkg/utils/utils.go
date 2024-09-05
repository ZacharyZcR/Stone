// pkg/utils/utils.go

package utils

import (
	"Stone/pkg/logging"
	"fmt"
	"net/http"
	"time"
)

// LogTraffic 记录流量日志
func LogTraffic(clientIP, targetIP, url, method string, headers http.Header, body, errorMsg string) {
	// 初始化日志数据
	logData := map[string]interface{}{
		"timestamp": time.Now(),
		"client_ip": clientIP,
		"target_ip": targetIP,
		"url":       url,
		"method":    method,
		"headers":   headers,
		"body":      body,
	}

	// 设置状态和错误信息
	if errorMsg != "" {
		logData["status"] = "failed"
		logData["error"] = errorMsg
	} else {
		logData["status"] = "success"
	}

	// 记录日志到存储
	if err := logging.LogTraffic(logData); err != nil {
		fmt.Printf("记录日志失败: %v\n", err)
	}
}
