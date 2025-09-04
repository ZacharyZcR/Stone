package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	Role         string    `bson:"role"`
	Created      time.Time `bson:"created"`
	LastLogin    time.Time `bson:"last_login"`
	Active       bool      `bson:"active"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: go run create_admin.go <用户名> <密码>")
		fmt.Println("示例: go run create_admin.go admin SecurePassword123")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]

	if len(password) < 8 {
		log.Fatal("密码必须至少8位字符")
	}

	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27019"))
	if err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	userCollection := client.Database("stoneDB").Collection("users")

	// 检查用户是否已存在
	var existingUser User
	err = userCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		log.Fatalf("用户 '%s' 已存在", username)
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建管理员用户
	admin := User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		Created:      time.Now(),
		Active:       true,
	}

	_, err = userCollection.InsertOne(context.Background(), admin)
	if err != nil {
		log.Fatalf("创建管理员用户失败: %v", err)
	}

	fmt.Printf("管理员用户 '%s' 创建成功\n", username)
	fmt.Println("请设置环境变量:")
	fmt.Println("export JWT_SECRET=\"your-secure-jwt-secret-here\"")
}