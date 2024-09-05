// pkg/capture/capture.go

package capture

import (
	"Stone/pkg/processing"
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

		go processing.HandleHTTPConnection(conn, targetAddress)
	}
}
