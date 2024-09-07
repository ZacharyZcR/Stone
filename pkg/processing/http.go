package processing

import (
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"Stone/pkg/utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
)

// HandleHTTPConnection 处理HTTP连接
func HandleHTTPConnection(clientConn net.Conn, targetAddress string) {
	defer clientConn.Close()

	// 获取客户端IP
	clientIP, _, _ := net.SplitHostPort(clientConn.RemoteAddr().String())

	// 尝试将IPv6地址转换为IPv4地址
	clientIP = convertIPv6ToIPv4(clientIP)

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
			monitoring.IncrementMetric("blockedByBlacklistTotal")
			sendBlockedResponse(clientConn, "blocked.html")
			return
		}

		// 如果IP不在白名单，进行URL和包体检查
		if !inWhitelist && !rules.CheckRequest(request) {
			fmt.Println("检测到危险请求，连接已阻断")
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", "Blocked by rules")
			monitoring.IncrementMetric("blockedByRulesTotal")
			sendBlockedResponse(clientConn, "blocked.html")
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
		err = monitoring.IncrementMetric("websiteRequestsTotal")
		if err != nil {
			log.Printf("Failed to increment websiteRequestsTotal: %v", err)
		}
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

func sendBlockedResponse(conn net.Conn, filePath string) {
	// 读取HTML文件内容
	htmlContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("无法读取被阻断响应文件:", err)
		return
	}

	// 生成随机状态码（200到503之间）
	randomStatusCode := rand.Intn(304) + 200

	// 生成随机长度的随机字符串（5000到10000个字符）
	randomLength := rand.Intn(5001) + 5000
	randomString := make([]byte, randomLength)
	for i := range randomString {
		randomString[i] = byte(rand.Intn(94) + 33) // 可打印ASCII字符
	}

	// 将随机字符串作为HTML注释插入到HTML内容中
	htmlWithRandomString := []byte(fmt.Sprintf("%s\n<!-- %s -->", htmlContent, randomString))

	// 构造响应
	response := fmt.Sprintf("HTTP/1.1 %d \r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s",
		randomStatusCode,
		len(htmlWithRandomString),
		htmlWithRandomString)

	// 发送响应
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("写回被阻断响应失败:", err)
	}
}

// 新增函数: 尝试将IPv6地址转换为IPv4地址
func convertIPv6ToIPv4(ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return ipAddress // 如果解析失败,返回原始地址
	}

	if ip.To4() != nil {
		return ip.To4().String() // 如果是IPv4或者可以转换为IPv4,返回IPv4地址
	}

	// 对于无法转换的IPv6地址,保持原样
	return ipAddress
}
