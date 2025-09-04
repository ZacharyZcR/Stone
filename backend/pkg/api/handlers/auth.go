package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

// Login 用户名密码登录
func Login(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码是必需的"})
			return
		}

		var user User
		err := userCollection.FindOne(context.Background(), bson.M{
			"username": request.Username,
			"active":   true,
		}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 更新最后登录时间
		_, err = userCollection.UpdateOne(
			context.Background(),
			bson.M{"username": request.Username},
			bson.M{"$set": bson.M{"last_login": time.Now()}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新登录时间失败"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Minute * 15).Unix(), // 15分钟过期
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":    tokenString,
			"username": user.Username,
			"role":     user.Role,
		})
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名(3-20位)和密码(至少8位)是必需的"})
		return
	}

	// 检查用户是否已存在
	var existingUser User
	err := userCollection.FindOne(context.Background(), bson.M{"username": request.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户 - 默认为普通用户角色
	newUser := User{
		Username:     request.Username,
		PasswordHash: string(hashedPassword),
		Role:         "user", // 默认角色
		Created:      time.Now(),
		Active:       true,
	}

	_, err = userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "注册成功",
		"username": request.Username,
	})
}

// CheckAuth 验证JWT并返回认证状态
func CheckAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "未提供令牌"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "无效的令牌格式"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
			c.JSON(http.StatusOK, gin.H{
				"authenticated": true,
				"username":     claims["username"],
				"role":         claims["role"],
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false, "error": "无效的令牌"})
		}
	}
}

