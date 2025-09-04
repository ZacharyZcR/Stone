package processing

import (
	"Stone/backend/pkg/health"
	"Stone/backend/pkg/monitoring"
	"Stone/backend/pkg/ratelimit"
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
	"sync"
	"time"
)

// ErrorType represents different types of HTTP errors
type ErrorType int

const (
	ErrorTLSRequired ErrorType = iota
	ErrorBadGateway
	ErrorForbidden
	ErrorInternal
	ErrorRateLimit
)

// HTTPProxy handles HTTP connections with connection pooling
type HTTPProxy struct {
	backendPool  *BackendPool
	errorHandler *ErrorHandler
}

// BackendPool manages HTTP client connections
type BackendPool struct {
	clients    sync.Map      // string -> *http.Client
	timeout    time.Duration
	maxClients int
	cleanup    *time.Ticker
}

// ErrorHandler manages error responses
type ErrorHandler struct {
	templates map[ErrorType][]byte
}

// NewHTTPProxy creates a new HTTP proxy with connection pooling
func NewHTTPProxy() *HTTPProxy {
	return &HTTPProxy{
		backendPool:  newBackendPool(),
		errorHandler: newErrorHandler(),
	}
}

// newBackendPool creates a backend connection pool with cleanup
func newBackendPool() *BackendPool {
	bp := &BackendPool{
		timeout:    30 * time.Second,
		maxClients: 100,
		cleanup:    time.NewTicker(5 * time.Minute),
	}
	
	// Start cleanup routine
	go bp.cleanupRoutine()
	return bp
}

// cleanupRoutine periodically cleans idle connections
func (bp *BackendPool) cleanupRoutine() {
	for range bp.cleanup.C {
		clientCount := 0
		bp.clients.Range(func(key, value interface{}) bool {
			clientCount++
			// For now, just count. In production, implement LRU eviction
			return true
		})
		
		if clientCount > bp.maxClients {
			log.Printf("Warning: %d HTTP clients in pool, consider increasing maxClients", clientCount)
		}
	}
}

// newErrorHandler creates an error handler with pre-compiled templates
func newErrorHandler() *ErrorHandler {
	eh := &ErrorHandler{
		templates: make(map[ErrorType][]byte),
	}
	
	// Pre-compile error responses to avoid string concatenation
	eh.templates[ErrorTLSRequired] = []byte(
		"HTTP/1.1 426 Upgrade Required\r\n" +
		"Upgrade: TLS/1.0, HTTP/1.1\r\n" +
		"Connection: Upgrade\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Length: 97\r\n" +
		"\r\n" +
		"<html><body><h1>426 Upgrade Required</h1><p>This service requires HTTPS/TLS.</p></body></html>")
	
	eh.templates[ErrorBadGateway] = []byte(
		"HTTP/1.1 502 Bad Gateway\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Length: 85\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"<html><body><h1>502 Bad Gateway</h1><p>The upstream server is down.</p></body></html>")
	
	eh.templates[ErrorInternal] = []byte(
		"HTTP/1.1 500 Internal Server Error\r\n" +
		"Content-Length: 0\r\n" +
		"Connection: close\r\n" +
		"\r\n")
	
	eh.templates[ErrorRateLimit] = []byte(
		"HTTP/1.1 429 Too Many Requests\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Length: 98\r\n" +
		"Retry-After: 300\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"<html><body><h1>429 Too Many Requests</h1><p>Rate limit exceeded. Try again later.</p></body></html>")
	
	return eh
}

// Global proxy instance - initialized once
var defaultProxy *HTTPProxy

func init() {
	defaultProxy = NewHTTPProxy()
}

// HandleHTTPConnection - simplified public interface
func HandleHTTPConnection(clientConn net.Conn, targetAddress string) {
	defaultProxy.HandleConnection(clientConn, targetAddress)
}

// HandleConnection processes HTTP connections
func (p *HTTPProxy) HandleConnection(conn net.Conn, target string) {
	defer conn.Close()
	
	clientIP := p.getClientIP(conn)
	reader := bufio.NewReader(conn)
	
	for {
		req, err := p.parseRequest(reader)
		if err != nil {
			p.handleRequestError(conn, clientIP, target, err)
			return
		}
		
		if blocked := p.checkAndBlockIfNeeded(conn, req, clientIP, target); blocked {
			return // Request was blocked, response already sent
		}
		
		resp, err := p.forwardRequest(req, target)
		if err != nil {
			p.errorHandler.sendError(conn, ErrorBadGateway)
			utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", err.Error())
			return
		}
		
		if err := p.writeResponse(conn, resp); err != nil {
			resp.Body.Close()
			utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", err.Error())
			return
		}
		
		p.logSuccess(clientIP, target, req)
		resp.Body.Close()
		
		if p.shouldClose(resp) {
			break
		}
	}
}

// getClientIP extracts and normalizes client IP
func (p *HTTPProxy) getClientIP(conn net.Conn) string {
	clientIP, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	return convertIPv6ToIPv4(clientIP)
}

// parseRequest reads and validates HTTP request
func (p *HTTPProxy) parseRequest(reader *bufio.Reader) (*http.Request, error) {
	return http.ReadRequest(reader)
}

// handleRequestError processes request parsing errors
func (p *HTTPProxy) handleRequestError(conn net.Conn, clientIP, target string, err error) {
	if err == io.EOF {
		return
	}
	
	if p.isTLSError(err) {
		log.Printf("TLS connection attempt from %s: %v", clientIP, err)
		p.errorHandler.sendError(conn, ErrorTLSRequired)
		return
	}
	
	log.Printf("Request parsing failed from %s: %v", clientIP, err)
	utils.LogTraffic(clientIP, target, "", "", nil, "", err.Error())
}

// isTLSError checks if error indicates TLS connection attempt
func (p *HTTPProxy) isTLSError(err error) bool {
	errorStr := err.Error()
	return strings.Contains(errorStr, "malformed HTTP request") ||
		strings.Contains(errorStr, "invalid method")
}

// checkAndBlockIfNeeded performs security checks and sends blocking response if needed
func (p *HTTPProxy) checkAndBlockIfNeeded(conn net.Conn, req *http.Request, clientIP, target string) bool {
	if health.IsSystemInDegradedMode() {
		log.Printf("Degraded mode - skipping security checks for %s", clientIP)
		return false
	}
	
	// 1. IP黑白名单检查
	allowed, inWhitelist := rules.IsAllowed(clientIP)
	if !allowed {
		log.Printf("IP blacklisted: %s", clientIP)
		utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", "IP blacklisted")
		monitoring.IncrementMetric("blockedByBlacklistTotal")
		p.errorHandler.sendBlockedResponse(conn, "blocked.html")
		return true
	}
	
	// 2. 速率限制检查（白名单IP豁免）
	if !inWhitelist {
		rateLimited, action := ratelimit.IsAllowed(clientIP)
		if !rateLimited {
			log.Printf("Rate limit exceeded for %s, action: %s", clientIP, action)
			utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", "Rate limit exceeded")
			monitoring.IncrementMetric("blockedByRateLimitTotal")
			
			switch action {
			case "block":
				p.errorHandler.sendError(conn, ErrorRateLimit)
			case "delay":
				time.Sleep(2 * time.Second) // 简单延迟
				p.errorHandler.sendError(conn, ErrorRateLimit)
			default:
				p.errorHandler.sendError(conn, ErrorRateLimit)
			}
			return true
		}
	}
	
	// 3. WAF规则检查
	if !inWhitelist && !rules.CheckRequest(req) {
		log.Printf("Request blocked by rules from %s", clientIP)
		utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", "Blocked by rules")
		monitoring.IncrementMetric("blockedByRulesTotal")
		p.errorHandler.sendBlockedResponse(conn, "blocked.html")
		return true
	}
	
	return false
}

// forwardRequest sends request to backend with connection pooling
func (p *HTTPProxy) forwardRequest(req *http.Request, target string) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = target
	req.RequestURI = ""
	
	client := p.backendPool.getClient(target)
	return client.Do(req)
}

// writeResponse streams response back to client
func (p *HTTPProxy) writeResponse(conn net.Conn, resp *http.Response) error {
	return resp.Write(conn)
}

// logSuccess records successful request
func (p *HTTPProxy) logSuccess(clientIP, target string, req *http.Request) {
	monitoring.IncrementMetric("websiteRequestsTotal")
	utils.LogTraffic(clientIP, target, req.URL.String(), req.Method, req.Header, "", "")
}

// shouldClose determines if connection should be closed
func (p *HTTPProxy) shouldClose(resp *http.Response) bool {
	return resp.Close || resp.Header.Get("Connection") == "close"
}

// getClient returns HTTP client for target with connection pooling
func (bp *BackendPool) getClient(target string) *http.Client {
	if client, ok := bp.clients.Load(target); ok {
		return client.(*http.Client)
	}
	
	// Create optimized transport for this target
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,   // Increased for better connection reuse
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	
	client := &http.Client{
		Timeout:   bp.timeout,
		Transport: transport,
	}
	
	bp.clients.Store(target, client)
	return client
}

// sendError sends pre-compiled error responses (zero allocation)
func (eh *ErrorHandler) sendError(conn net.Conn, errType ErrorType) {
	if response, ok := eh.templates[errType]; ok {
		conn.Write(response)
	} else {
		conn.Write(eh.templates[ErrorInternal])
	}
}

// sendBlockedResponse sends WAF blocking response
func (eh *ErrorHandler) sendBlockedResponse(conn net.Conn, filePath string) {
	htmlContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read blocked response file: %v", err)
		eh.sendError(conn, ErrorForbidden)
		return
	}
	
	// Generate anti-fingerprint noise
	randomLength := rand.Intn(2001) + 1000
	randomString := make([]byte, randomLength)
	for i := range randomString {
		randomString[i] = byte(rand.Intn(94) + 33)
	}
	
	htmlWithNoise := fmt.Sprintf("%s\n<!-- %s -->", htmlContent, randomString)
	response := fmt.Sprintf("HTTP/1.1 403 Forbidden\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s",
		len(htmlWithNoise),
		htmlWithNoise)
	
	conn.Write([]byte(response))
}

// convertIPv6ToIPv4 normalizes IPv6-mapped IPv4 addresses
func convertIPv6ToIPv4(ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return ipAddress
	}
	
	if ip4 := ip.To4(); ip4 != nil {
		return ip4.String()
	}
	
	return ipAddress
}
