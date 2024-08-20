// pkg/capture/capture.go

package capture

import (
	"Stone/pkg/processing"
	"Stone/pkg/rules"
	"Stone/pkg/utils"
	"fmt"
	"net"
)

// StartCapture 启动流量捕获
func StartCapture(port int, targetAddress string) error {
	// 监听指定端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("监听端口失败: %w", err)
	}
	defer listener.Close()

	fmt.Printf("流量捕获已启动，监听端口: %d\n", port)

	// 接受并处理连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
			continue
		}

		clientIP, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		allowed, _ := rules.IsAllowed(clientIP)
		if !allowed {
			fmt.Printf("IP在黑名单中，连接已阻断: %s\n", clientIP)
			utils.LogTraffic(clientIP, "blocked", "", "")
			conn.Close()
			continue
		}

		go processing.HandleHTTPConnection(conn, targetAddress)
	}
}
