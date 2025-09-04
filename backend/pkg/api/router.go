package api

import (
	"Stone/backend/pkg/api/handlers"
	"Stone/backend/pkg/logging"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

// SetupRouter 设置API路由
func SetupRouter(configCollection *mongo.Collection, userCollection *mongo.Collection) *gin.Engine {
	router := gin.Default()

	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://10.31.2.243:8084", 
			"http://localhost:8084",
			"http://localhost:8085", 
			"http://localhost:8086",
			"http://127.0.0.1:8084",
			"http://127.0.0.1:8085",
			"http://127.0.0.1:8086",
		}, // 允许的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许跨域请求携带认证信息
		MaxAge:           12 * time.Hour,
	}))

	// 从环境变量获取密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logging.LogError(fmt.Errorf("JWT_SECRET environment variable must be set"))
		return nil
	}

	// Session middleware removed - using JWT only

	// 认证API
	router.POST("/auth/register", handlers.Register) // 公开注册
	router.POST("/auth/login", handlers.Login(jwtSecret))
	router.GET("/auth/check", handlers.CheckAuth(jwtSecret))

	// 使用中间件进行鉴权
	authenticated := router.Group("/")
	authenticated.Use(AuthMiddleware(jwtSecret))

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

		// 新增的IP统计API
		authenticated.GET("/ip-stats", handlers.GetIPStats)

		// 攻击者画像
		authenticated.GET("/attacker-profile", handlers.GetAttackerProfile)


		// 用户管理API
		authenticated.GET("/users", handlers.GetUsers)
		authenticated.POST("/users", handlers.CreateUser)
		authenticated.DELETE("/users/:username", handlers.DeleteUser)

		// 防火墙指标API
		authenticated.GET("/firewall/metrics", handlers.GetFirewallMetrics)
	}

	return router
}
