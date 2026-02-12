package api

import (
	"encoding/json"
	"health_app/configs"
	"health_app/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

// User represents the structure for user data
type User = models.User

// RegisterUser handles user registration (只需要用户名和密码)
func RegisterUser(c *gin.Context) {
	log.Printf("[API] 收到注册请求: %s %s", c.Request.Method, c.Request.URL.Path)

	var user User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		log.Printf("[API] 注册失败: 无效的请求体 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Printf("[API] 注册用户: name=%s", user.Name)

	// Insert user into the database (email 为空)
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	_, err = configs.DB.Exec(query, user.Name, "", user.Password)
	if err != nil {
		log.Printf("[API] 注册失败: 数据库错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	log.Printf("[API] 用户注册成功: name=%s", user.Name)
	c.JSON(201, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user login (使用用户名和密码登录)
func LoginUser(c *gin.Context) {
	log.Printf("[API] 收到登录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	var user User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		log.Printf("[API] 登录失败: 无效的请求体 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Printf("[API] 用户登录: name=%s", user.Name)

	// Check user credentials (使用 name 和 password 登录)
	var dbUser User
	query := "SELECT id, name, email, password FROM users WHERE name = ? AND password = ?"
	err = configs.DB.QueryRow(query, user.Name, user.Password).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password)
	if err != nil {
		log.Printf("[API] 登录失败: 无效的用户名或密码 - name=%s", user.Name)
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	log.Printf("[API] 用户登录成功: id=%d, name=%s", dbUser.ID, dbUser.Name)
	c.JSON(200, dbUser)
}
