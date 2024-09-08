package handlers

import (
	"Stone/pkg/rules"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		// 检查 IP 是否为空
		if newRule.IP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IP address cannot be empty"})
			return
		}
		if err := rules.AddIPRule(newRule); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule added"})
	case http.MethodDelete:
		ip := c.Param("ip")
		// 检查 IP 是否为空
		if ip == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IP address cannot be empty"})
			return
		}
		if err := rules.DeleteIPRule(ip); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "IP rule deleted"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}
