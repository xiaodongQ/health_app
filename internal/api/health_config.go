package api

import (
	"log"

	"health_app/configs"

	"github.com/gin-gonic/gin"
)

// GetHealthConfig 返回所有健康检查类别和对应指标的配置
func GetHealthConfig(c *gin.Context) {
	log.Printf("[API] 收到获取健康配置请求: %s %s", c.Request.Method, c.Request.URL.Path)
	c.JSON(200, configs.DefaultHealthConfig)
}
