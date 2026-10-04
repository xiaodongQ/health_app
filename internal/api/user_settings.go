package api

import (
	"health_app/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserSetting(c *gin.Context) {
	userSettingsRepo := repository.NewUserSettingsRepository()
	userIDStr := c.Query("user_id")
	key := c.Query("key")
	if userIDStr == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id or key"})
		return
	}
	userID, _ := strconv.Atoi(userIDStr)
	value, err := userSettingsRepo.Get(userID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": value})
}

func SetUserSetting(c *gin.Context) {
	userSettingsRepo := repository.NewUserSettingsRepository()
	var req struct {
		UserID int    `json:"user_id"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userSettingsRepo.Set(req.UserID, req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "setting saved"})
}
