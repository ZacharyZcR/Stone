package processing

import (
	"Stone/backend/pkg/monitoring"
	"Stone/backend/pkg/rules"
	"Stone/backend/pkg/utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
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
				// 检查是否是TLS连接尝试
				if strings.Contains(err.Error(), "malformed HTTP request") || strings.Contains(err.Error(), "invalid method") {
					fmt.Printf("检测到非HTTP协议连接 (可能是HTTPS/TLS): %v\n", err)
					// 返回426 Upgrade Required
					upgradeResponse := "HTTP/1.1 426 Upgrade Required\r\n" +
						"Upgrade: TLS/1.0, HTTP/1.1\r\n" +
						"Connection: Upgrade\r\n" +
						"Content-Type: text/html; charset=UTF-8\r\n" +
						"Content-Length: 97\r\n" +
						"\r\n" +
						"<html><body><h1>426 Upgrade Required</h1><p>This service requires HTTPS/TLS.</p></body></html>"
					
					clientConn.Write([]byte(upgradeResponse))
				} else {
					fmt.Printf("读取HTTP请求失败: %v\n", err)
				}
			}
			// 不记录TLS握手失败为错误日志，减少噪音
			if !strings.Contains(err.Error(), "malformed HTTP request") {
				utils.LogTraffic(clientIP, targetAddress, "", "", nil, "", err.Error())
			}
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
			fmt.Printf("目标服务不可用 %s: %v\n", targetAddress, err)
			utils.LogTraffic(clientIP, targetAddress, request.URL.String(), request.Method, request.Header, "", fmt.Sprintf("目标服务不可用: %v", err))
			
			// 返回502 Bad Gateway错误给客户端
			badGatewayResponse := "HTTP/1.1 502 Bad Gateway\r\n" +
				"Content-Type: text/html; charset=UTF-8\r\n" +
				"Content-Length: 85\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				"<html><body><h1>502 Bad Gateway</h1><p>The upstream server is down.</p></body></html>"
			
			clientConn.Write([]byte(badGatewayResponse))
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

	// 使用403 Forbidden状态码，符合WAF拦截语义
	statusCode := 403
	statusText := "Forbidden"

	// 生成随机长度的随机字符串（1000到3000个字符），减少资源消耗
	randomLength := rand.Intn(2001) + 1000
	randomString := make([]byte, randomLength)
	for i := range randomString {
		randomString[i] = byte(rand.Intn(94) + 33) // 可打印ASCII字符
	}

	// 将随机字符串作为HTML注释插入到HTML内容中
	htmlWithRandomString := []byte(fmt.Sprintf("%s\n<!-- %s -->", htmlContent, randomString))

	// 构造响应
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s",
		statusCode,
		statusText,
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
