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
	logData := map[string]interface{}{
		"timestamp": time.Now(),
		"client_ip": clientIP,
		"target_ip": targetIP,
		"url":       url,
		"method":    method,
		"headers":   headers,
		"body":      body,
		"status":    "success",
	}
	if errorMsg != "" {
		logData["status"] = "failed"
		logData["error"] = errorMsg
	}

	if err := logging.LogTraffic(logData); err != nil {
		fmt.Printf("记录日志失败: %v\n", err)
	}
}
