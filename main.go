package main

import (
	"health_app/configs"
	"health_app/internal/api"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 确保数据库连接初始化完成 - 必须在任何API路由注册之前执行
	configs.InitDB()

	// 设置Gin模式为发布模式以提高性能
	gin.SetMode(gin.ReleaseMode)

	// 初始化路由器
	r := gin.Default()

	// 获取当前文件的绝对路径，用于定位 web 目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	webDir := filepath.Join(dir, "web")

	// API路由组 - 首先注册API路由
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", func(c *gin.Context) {
			log.Printf("[API] 健康检查请求")
			c.JSON(200, gin.H{"message": "Health App API is running"})
		})

		// 用户认证相关接口
		apiGroup.POST("/users/register", api.RegisterUser)
		apiGroup.POST("/users/login", api.LoginUser)

		// 健康记录相关接口
		apiGroup.POST("/health_records", api.AddHealthRecord)
		apiGroup.GET("/health_records", api.GetHealthRecords)
		apiGroup.GET("/health_records/:id", api.GetHealthRecordByID)
		apiGroup.PUT("/health_records/:id", api.UpdateHealthRecord)
		apiGroup.DELETE("/health_records/:id", api.DeleteHealthRecord)
		apiGroup.GET("/health_records/date", api.GetHealthRecordByDate) // 根据用户ID和日期获取健康记录
		apiGroup.GET("/trends", api.GetTrendData)
		apiGroup.GET("/indicators", api.GetAllIndicators)
		apiGroup.GET("/health_config", api.GetHealthConfig) // 新增健康配置API

		// 提醒功能相关接口
		apiGroup.POST("/reminders", api.AddReminder)
		apiGroup.GET("/reminders", api.GetReminders)
		apiGroup.PUT("/reminders", api.UpdateReminder)
		apiGroup.DELETE("/reminders", api.DeleteReminder)
		apiGroup.GET("/reminders/upcoming", api.GetUpcomingReminders)

		// 用户设置
		apiGroup.GET("/settings", api.GetUserSetting)
		apiGroup.POST("/settings", api.SetUserSetting)
	}

	// 静态文件服务 - 使用 /static 路径避免与API路由冲突
	r.Static("/static", webDir)

	// 主页路由
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(webDir, "index.html"))
	})

	// 处理前端路由 - 对于非API且非静态资源的请求，返回index.html
	r.NoRoute(func(c *gin.Context) {
		// 检查是否是API请求
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// 检查是否是静态资源请求
		ext := filepath.Ext(c.Request.URL.Path)
		if ext != "" { // 如果有扩展名，则认为是静态资源
			// 尝试查找静态文件
			filePath := filepath.Join(webDir, c.Request.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}
		}

		// 对于没有扩展名的请求（前端路由），返回index.html
		c.File(filepath.Join(webDir, "index.html"))
	})

	// Start the server
	log.Println("Starting server on :8081...")
	log.Fatal(r.Run(":8081"))
}
