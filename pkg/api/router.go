package api

import (
	"Stone/pkg/api/handlers"
	"Stone/pkg/logging"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// SetupRouter 设置API路由
func SetupRouter(configCollection *mongo.Collection, userCollection *mongo.Collection) *gin.Engine {
	router := gin.Default()

	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许跨域请求携带认证信息
		MaxAge:           12 * time.Hour,
	}))

	// 从数据库加载配置
	var configResult struct {
		Secrets struct {
			SessionSecret string `bson:"sessionSecret"`
			JWTSecret     string `bson:"jwtSecret"`
		} `bson:"secrets"`
	}
	err := configCollection.FindOne(context.Background(), bson.M{"type": "config"}).Decode(&configResult)
	if err != nil {
		logging.LogError(fmt.Errorf("加载配置失败: %v", err))
		return nil
	}

	// 设置Session中间件
	store := cookie.NewStore([]byte(configResult.Secrets.SessionSecret))
	router.Use(sessions.Sessions("mysession", store))

	// 验证JWT的API
	router.GET("/auth/check", handlers.CheckAuth(configResult.Secrets.JWTSecret))

	router.GET("/auth/qrcode", handlers.GenerateQRCode)
	router.POST("/auth/validate", handlers.ValidateTOTP(configResult.Secrets.JWTSecret))

	// 使用中间件进行鉴权
	authenticated := router.Group("/")
	authenticated.Use(AuthMiddleware(configResult.Secrets.JWTSecret))

	{
		// 系统状态API
		authenticated.GET("/status", handlers.GetStatus)

		// IP控制规则管理API
		authenticated.GET("/ip-control-rules", handlers.HandleIPControlRules)
		authenticated.GET("/ip-control-rules/:ip", handlers.HandleIPControlRules)
		authenticated.POST("/ip-control-rules", handlers.HandleIPControlRules)
		authenticated.DELETE("/ip-control-rules/:ip", handlers.HandleIPControlRules)

		// 拦截规则管理API
		authenticated.GET("/interception-rules", handlers.HandleInterceptionRules)
		authenticated.GET("/interception-rules/:name", handlers.HandleInterceptionRules)
		authenticated.POST("/interception-rules", handlers.HandleInterceptionRules)
		authenticated.DELETE("/interception-rules/:name", handlers.HandleInterceptionRules)

		// 日志查看API
		authenticated.GET("/logs", handlers.GetLogs)

		// 控制二维码接口状态的API
		authenticated.GET("/auth/qrcode/status", handlers.SetQRCodeStatus)
		authenticated.POST("/auth/qrcode/status", handlers.SetQRCodeStatus)

		// 用户管理API
		authenticated.GET("/users", handlers.HandleUsers)
		authenticated.GET("/users/:account", handlers.HandleUsers)
		authenticated.POST("/users", handlers.HandleUsers)
		authenticated.DELETE("/users/:account", handlers.HandleUsers)

		// 防火墙指标API
		authenticated.GET("/firewall/metrics", handlers.GetFirewallMetrics)
	}

	return router
}
