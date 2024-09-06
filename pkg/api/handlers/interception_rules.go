package handlers

import (
	"Stone/pkg/rules"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
