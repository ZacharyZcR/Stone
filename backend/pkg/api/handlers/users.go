package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// MongoDB 集合
var userCollection *mongo.Collection

type User struct {
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `bson:"password_hash" json:"-"`
	Role         string    `bson:"role" json:"role"`
	Created      time.Time `bson:"created" json:"created"`
	LastLogin    time.Time `bson:"last_login" json:"last_login"`
	Active       bool      `bson:"active" json:"active"`
}

// SetUserCollection 设置用户集合
func SetUserCollection(collection *mongo.Collection) {
	userCollection = collection
}

// CreateUser 创建新用户（仅管理员）
func CreateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密码(至少8位)和角色是必需的"})
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

	newUser := User{
		Username:     request.Username,
		PasswordHash: string(hashedPassword),
		Role:         request.Role,
		Created:      time.Now(),
		Active:       true,
	}

	_, err = userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建用户"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用户创建成功", "username": request.Username})
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	cursor, err := userCollection.Find(context.Background(), bson.M{"active": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户列表"})
		return
	}
	defer cursor.Close(context.Background())

	var users []User
	if err := cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析用户数据"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUser 删除用户（软删除）
func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名是必需的"})
		return
	}

	_, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"username": username},
		bson.M{"$set": bson.M{"active": false}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法删除用户"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}
