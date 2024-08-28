// pkg/api/router.go

package api

import (
	"Stone/pkg/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// SetupRouter 设置API路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许跨域请求携带认证信息
		MaxAge:           12 * time.Hour,
	}))

	// 系统状态API
	router.GET("/status", handlers.GetStatus)

	// IP控制规则管理API
	router.GET("/ip-control-rules", handlers.GetIPControlRules)
	router.POST("/ip-control-rules", handlers.UpdateIPControlRules)

	// 拦截规则管理API
	router.GET("/interception-rules", handlers.GetInterceptionRules)
	router.POST("/interception-rules", handlers.UpdateInterceptionRules)

	// 日志查看API
	router.GET("/logs", handlers.GetLogs)

	return router
}
