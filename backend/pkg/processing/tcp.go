// pkg/processing/tcp.go

package processing

import (
	"fmt"
	"io"
	"net"
)

// HandleTCPConnection 处理TCP连接并转发流量
func HandleTCPConnection(clientConn net.Conn, targetAddress string) {
	defer clientConn.Close()

	// 连接到目标服务
	targetConn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		fmt.Println("无法连接到目标服务:", err)
		return
	}
	defer targetConn.Close()

	// 双向转发数据
	go func() {
		if _, err := io.Copy(targetConn, clientConn); err != nil {
			fmt.Println("转发到目标服务失败:", err)
		}
	}()

	if _, err := io.Copy(clientConn, targetConn); err != nil {
		fmt.Println("从目标服务接收数据失败:", err)
	}
}
