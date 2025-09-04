// pkg/utils/utils.go

package utils

import (
	"Stone/backend/pkg/logging"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	// 敏感信息正则表达式
	passwordRegex  = regexp.MustCompile(`(?i)(password|passwd|pwd)["':\s=]+[^"'\s&]+`)
	tokenRegex     = regexp.MustCompile(`(?i)(token|jwt|key|secret)["':\s=]+[^"'\s&]+`)
	cardRegex      = regexp.MustCompile(`\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}`)
	emailRegex     = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex     = regexp.MustCompile(`1[3-9]\d{9}`)
)

// sanitizeData 脱敏敏感信息
func sanitizeData(data string) string {
	if len(data) > 1000 {
		data = data[:1000] + "...[truncated]"
	}
	
	data = passwordRegex.ReplaceAllString(data, "$1***")
	data = tokenRegex.ReplaceAllString(data, "$1***")
	data = cardRegex.ReplaceAllString(data, "****-****-****-****")
	data = emailRegex.ReplaceAllString(data, "***@***.***")
	data = phoneRegex.ReplaceAllString(data, "***********")
	
	return data
}

// sanitizeHeaders 脱敏HTTP头部
func sanitizeHeaders(headers http.Header) http.Header {
	sanitized := make(http.Header)
	sensitiveHeaders := []string{
		"Authorization", "Cookie", "X-Auth-Token", 
		"X-Api-Key", "X-Session-Id", "X-Csrf-Token",
	}
	
	for key, values := range headers {
		sanitized[key] = make([]string, len(values))
		for i, value := range values {
			// 检查是否是敏感头部
			isSensitive := false
			for _, sensitive := range sensitiveHeaders {
				if strings.EqualFold(key, sensitive) {
					isSensitive = true
					break
				}
			}
			
			if isSensitive {
				sanitized[key][i] = "***"
			} else {
				sanitized[key][i] = sanitizeData(value)
			}
		}
	}
	
	return sanitized
}

// LogTraffic 记录流量日志
func LogTraffic(clientIP, targetIP, url, method string, headers http.Header, body, errorMsg string) {
	// 初始化日志数据
	logData := map[string]interface{}{
		"timestamp": time.Now(),
		"client_ip": clientIP,
		"target_ip": targetIP,
		"url":       sanitizeData(url),
		"method":    method,
		"headers":   sanitizeHeaders(headers),
		"body":      sanitizeData(body),
	}

	// 设置状态和错误信息
	if errorMsg != "" {
		logData["status"] = "failed"
		logData["error"] = sanitizeData(errorMsg)
	} else {
		logData["status"] = "success"
	}

	// 记录日志到存储
	if err := logging.LogTraffic(logData); err != nil {
		fmt.Printf("记录日志失败: %v\n", err)
	}
}
