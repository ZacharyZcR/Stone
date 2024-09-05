// pkg/api/handlers/handlers.go

package handlers

import (
	"Stone/pkg/logging"
	"Stone/pkg/monitoring"
	"Stone/pkg/rules"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pquerna/otp/totp"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"image/png"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// GetStatus 获取系统状态
func GetStatus(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 获取系统启动时间
	uptime := time.Since(monitoring.StartTime).String()

	// 获取处理的请求总数（假设使用Prometheus）
	requestsTotal, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取请求总数"})
		return
	}

	var totalRequests float64
	for _, mf := range requestsTotal {
		if mf.GetName() == "stone_requests_total" {
			for _, metric := range mf.GetMetric() {
				totalRequests += metric.GetCounter().GetValue()
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "running",
		"uptime":           uptime,
		"goroutines":       runtime.NumGoroutine(),
		"memory_alloc":     memStats.Alloc,
		"total_requests":   totalRequests,
		"cpu_cores":        runtime.NumCPU(),
		"cpu_architecture": runtime.GOARCH,
		"go_version":       runtime.Version(),
	})
}

// GetLogs 查看日志
func GetLogs(c *gin.Context) {
	// 解析查询参数
	startDateTimeStr := c.Query("startDateTime")
	endDateTimeStr := c.Query("endDateTime")
	ip := c.Query("ip")

	// 解析时间参数
	var startDateTime, endDateTime time.Time
	var err error
	if startDateTimeStr != "" {
		startDateTime, err = time.Parse(time.RFC3339, startDateTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式，必须为RFC3339格式"})
			return
		}
	}
	if endDateTimeStr != "" {
		endDateTime, err = time.Parse(time.RFC3339, endDateTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式，必须为RFC3339格式"})
			return
		}
	}

	// 获取日志
	logs, err := logging.FetchLogsFromMongoWithFilters(context.Background(), 100, startDateTime, endDateTime, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取日志"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// HandleIPControlRules 处理IP控制规则的操作
func HandleIPControlRules(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		ip := c.Param("ip")
		if ip == "" {
			// 获取所有IP控制规则
			currentRules := rules.GetIPControlRules()
			c.JSON(http.StatusOK, currentRules)
		} else {
			// 获取特定IP的规则
			rule, found := rules.GetIPRule(ip)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
				return
			}
			c.JSON(http.StatusOK, rule)
		}
	case http.MethodPost:
		var newRule rules.IPControlRule
		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := rules.AddIPRule(newRule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule added"})
	case http.MethodDelete:
		ip := c.Param("ip")
		if err := rules.DeleteIPRule(ip); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule deleted"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

// HandleInterceptionRules 处理拦截规则的操作
func HandleInterceptionRules(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		name := c.Param("name")
		if name == "" {
			// 获取所有拦截规则
			currentRules := rules.GetInterceptionRules()
			c.JSON(http.StatusOK, currentRules)
		} else {
			// 获取特定名称的规则
			rule, found := rules.GetInterceptionRule(name)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
				return
			}
			c.JSON(http.StatusOK, rule)
		}
	case http.MethodPost:
		var newRule rules.Pattern
		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := rules.AddInterceptionRule(newRule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add interception rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Interception rule added"})
	case http.MethodDelete:
		name := c.Param("name")
		if err := rules.DeleteInterceptionRule(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interception rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Interception rule deleted"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

// MongoDB 集合
var totpCollection *mongo.Collection

// SetTOTPCollection 设置MongoDB集合
func SetTOTPCollection(collection *mongo.Collection) {
	totpCollection = collection
}

// 随机生成一个16位的字符串
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("无法生成随机字符串: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// GenerateQRCode 生成二维码并返回给客户端
func GenerateQRCode(c *gin.Context) {
	// 生成一个随机的账户名称
	accountName, err := generateRandomString(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成账户名称"})
		return
	}

	// 生成一个新的密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "StoneAdmin",
		AccountName: accountName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成密钥"})
		return
	}

	// 将密钥存储在MongoDB中
	_, err = totpCollection.InsertOne(context.Background(), bson.M{
		"account": accountName,
		"secret":  key.Secret(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法存储密钥"})
		return
	}

	// 生成二维码
	img, err := key.Image(200, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成二维码"})
		return
	}

	// 将二维码编码为PNG并返回
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法编码二维码"})
		return
	}

	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

// ValidateTOTP 验证用户输入的TOTP代码
func ValidateTOTP(c *gin.Context) {
	var request struct {
		Code    string `json:"code" binding:"required"`
		Account string `json:"account" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码和账户是必需的"})
		return
	}

	// 从MongoDB中获取存储的密钥
	var result struct {
		Secret string `bson:"secret"`
	}
	err := totpCollection.FindOne(context.Background(), bson.M{"account": request.Account}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法找到密钥"})
		return
	}

	// 验证用户输入的TOTP代码
	if totp.Validate(request.Code, result.Secret) {
		// 生成JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"account": request.Account,
			"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3天有效期
		})

		// 使用环境变量中的密钥签名JWT
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "验证码有效", "token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "验证码无效"})
	}
}

// CheckAuth 验证JWT并返回认证状态
func CheckAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取JWT
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "未提供令牌"})
			return
		}

		// 去掉 "Bearer " 前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "无效的令牌格式"})
			return
		}

		// 验证JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 确保使用的是预期的签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "无效的令牌"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.JSON(http.StatusOK, gin.H{"authenticated": true, "account": claims["account"]})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "无效的令牌"})
		}
	}
}
