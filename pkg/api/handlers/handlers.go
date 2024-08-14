// pkg/api/handlers/handlers.go

package handlers

import (
	"Stone/pkg/rules"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStatus 获取系统状态
func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"uptime": "24h", // 示例数据
	})
}

// GetRules 获取当前规则
func GetRules(c *gin.Context) {
	// 示例：假设我们有一个全局变量存储规则
	currentRules := rules.Rules{
		Whitelist: []string{"192.168.1.100"},
		Blacklist: []string{"192.168.1.200"},
	}
	c.JSON(http.StatusOK, currentRules)
}

// UpdateRules 更新规则
func UpdateRules(c *gin.Context) {
	var newRules rules.Rules
	if err := c.ShouldBindJSON(&newRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新全局规则（示例）
	// 实际应用中应有更好的管理方式
	// globalRules.UpdateRules(newRules.Whitelist, newRules.Blacklist)

	c.JSON(http.StatusOK, gin.H{"status": "rules updated"})
}
