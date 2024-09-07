// pkg/processing/http.go

package processing

import (
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"Stone/pkg/utils"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

// HandleHTTPConnection 处理HTTP连接
func HandleHTTPConnection(clientConn net.Conn, targetAddress string) {
	defer clientConn.Close()

	// 获取客户端IP
	clientIP, _, _ := net.SplitHostPort(clientConn.RemoteAddr().String())

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
			utils.LogTraffic(clientIP, targetAddress, "", "", nil, "", err.Error())
			return
		}

		// 检查IP是否在黑名单
		allowed, inWhitelist := rules.IsAllowed(clientIP)
		if !allowed {
			fmt.Printf("IP在黑名单中，连接已阻断: %s\n", clientIP)
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", "IP在黑名单中")
			monitoring.BlockedByBlacklistTotal.Inc()
			return
		}

		// 如果IP不在白名单，进行URL和包体检查
		if !inWhitelist && !rules.CheckRequest(request) {
			fmt.Println("检测到危险请求，连接已阻断")
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", "Blocked by rules")
			monitoring.BlockedByRulesTotal.Inc()
			return
		}

		// 设置目标地址
		request.URL.Scheme = "http"
		request.URL.Host = targetAddress
		request.RequestURI = ""

		// 发送请求到目标服务
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("发送请求到目标服务失败:", err)
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", err.Error())
			return
		}

		// 将响应写回客户端
		if err := response.Write(clientConn); err != nil {
			fmt.Println("写回客户端失败:", err)
			response.Body.Close()
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", err.Error())
			return
		}

		// 请求成功，更新访问计数
		monitoring.WebsiteRequestsTotal.WithLabelValues("success").Inc()
		utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", "")

		// 关闭响应体
		response.Body.Close()

		// 检查是否需要保持连接
		if !response.Close && response.Header.Get("Connection") != "close" {
			continue
		}
		break
	}
}
