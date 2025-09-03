// pkg/logging/logging.go

package logging

import (
	"log"
	"os"
)

// 初始化日志记录器
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// 创建日志文件
	logFile, err := os.OpenFile("stone.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}

	// 初始化日志记录器
	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo 记录信息日志
func LogInfo(message string) {
	InfoLogger.Println(message)
}

// LogError 记录错误日志
func LogError(err error) {
	ErrorLogger.Println(err)
}
