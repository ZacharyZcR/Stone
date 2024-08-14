// pkg/processing/http.go

package processing

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

// HandleHTTPConnection 处理HTTP连接
func HandleHTTPConnection(clientConn net.Conn, targetAddress string) {
	defer clientConn.Close()

	// 创建bufio.Reader
	reader := bufio.NewReader(clientConn)

	// 创建HTTP客户端
	client := &http.Client{}

	for {
		// 读取客户端请求
		request, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("读取HTTP请求失败:", err)
			}
			return
		}

		// 修改请求的目标地址
		request.URL.Host = targetAddress
		request.URL.Scheme = "http"
		request.RequestURI = ""

		// 发送请求到目标服务
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("发送请求到目标服务失败:", err)
			return
		}

		// 将响应写回客户端
		if err := response.Write(clientConn); err != nil {
			fmt.Println("写回客户端失败:", err)
			response.Body.Close()
			return
		}

		// 关闭响应体
		response.Body.Close()

		// 检查是否需要保持连接
		if !response.Close && response.Header.Get("Connection") != "close" {
			continue
		}
		break
	}
}
