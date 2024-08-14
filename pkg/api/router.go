// pkg/api/router.go

package api

import (
	"Stone/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置API路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 系统状态API
	router.GET("/status", handlers.GetStatus)

	// 规则管理API
	router.GET("/rules", handlers.GetRules)
	router.POST("/rules", handlers.UpdateRules)

	return router
}
