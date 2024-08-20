// pkg/capture/capture.go

package capture

import (
	"Stone/pkg/processing"
	"Stone/pkg/rules"
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
		if !rules.IsAllowed(clientIP) {
			fmt.Printf("拒绝连接: %s\n", clientIP)
			conn.Close()
			continue
		}

		go func(conn net.Conn) {
			// 简单协议识别，根据需要选择处理器
			// 这里假设所有流量都是HTTP
			processing.HandleHTTPConnection(conn, targetAddress)
		}(conn)
	}
}
