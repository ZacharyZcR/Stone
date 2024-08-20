// pkg/utils/utils.go

package utils

import (
	"Stone/pkg/logging"
	"fmt"
	"time"
)

// LogTraffic 记录流量日志
func LogTraffic(clientIP, status, url, errorMsg string) {
	logData := map[string]interface{}{
		"timestamp": time.Now(),
		"client_ip": clientIP,
		"status":    status,
	}
	if url != "" {
		logData["url"] = url
	}
	if errorMsg != "" {
		logData["error"] = errorMsg
	}

	if err := logging.LogTraffic(logData); err != nil {
		fmt.Printf("记录日志失败: %v\n", err)
	}
}
