package api

import (
	"encoding/json"
	"log"
	"strconv"

	"health_app/internal/models"
	"health_app/internal/repository"
	"health_app/internal/service"

	"github.com/gin-gonic/gin"
)

var reminderService *service.ReminderService

// getReminderService 获取提醒服务实例（懒加载）
func getReminderService() *service.ReminderService {
	if reminderService == nil {
		reminderRepo := repository.NewReminderRepository()
		reminderService = service.NewReminderService(reminderRepo)
	}
	return reminderService
}

// AddReminder 添加提醒
func AddReminder(c *gin.Context) {
	log.Printf("[API] 收到添加提醒请求: %s %s", c.Request.Method, c.Request.URL.Path)

	var reminder models.Reminder
	err := json.NewDecoder(c.Request.Body).Decode(&reminder)
	if err != nil {
		log.Printf("[API] 添加提醒失败: 无效的请求体 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// 验证必需字段
	if reminder.UserID == 0 || reminder.Type == "" || reminder.Title == "" || reminder.NextDate.IsZero() {
		log.Printf("[API] 添加提醒失败: 缺少必需字段 - UserID: %d, Type: %s, Title: %s, NextDate: %v",
			reminder.UserID, reminder.Type, reminder.Title, reminder.NextDate)
		c.JSON(400, gin.H{"error": "User ID, Type, Title, and Next Date are required"})
		return
	}

	// 验证提醒类型
	if reminder.Type != "medication" && reminder.Type != "test" {
		log.Printf("[API] 添加提醒失败: 无效的提醒类型 - %s", reminder.Type)
		c.JSON(400, gin.H{"error": "Invalid reminder type. Must be 'medication' or 'test'"})
		return
	}

	// 使用服务层添加提醒
	service := getReminderService()
	reminderID, err := service.CreateReminder(&reminder)
	if err != nil {
		log.Printf("[API] 添加提醒失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to create reminder: " + err.Error()})
		return
	}

	log.Printf("[API] 提醒添加成功: reminder_id=%d, user_id=%d", reminderID, reminder.UserID)
	c.JSON(201, gin.H{
		"message": "Reminder added successfully",
		"id":      reminderID,
	})
}

// GetReminders 获取用户的提醒
func GetReminders(c *gin.Context) {
	log.Printf("[API] 收到获取提醒请求: %s %s", c.Request.Method, c.Request.URL.Path)

	userIDStr := c.Query("user_id")
	activeOnly := c.Query("active_only")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取提醒失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	// 使用服务层获取提醒
	service := getReminderService()
	activeOnlyBool := activeOnly == "true"
	reminders, err := service.GetReminders(userID, activeOnlyBool)
	if err != nil {
		log.Printf("[API] 获取提醒失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve reminders: " + err.Error()})
		return
	}

	log.Printf("[API] 提醒查询成功: user_id=%d, 共 %d 条提醒", userID, len(reminders))
	c.JSON(200, reminders)
}

// UpdateReminder 更新提醒
func UpdateReminder(c *gin.Context) {
	log.Printf("[API] 收到更新提醒请求: %s %s", c.Request.Method, c.Request.URL.Path)

	reminderIDStr := c.Query("id")
	reminderID, err := strconv.Atoi(reminderIDStr)
	if err != nil {
		log.Printf("[API] 更新提醒失败: 无效的提醒ID - %s", reminderIDStr)
		c.JSON(400, gin.H{"error": "Invalid Reminder ID"})
		return
	}

	var reminder models.Reminder
	err = json.NewDecoder(c.Request.Body).Decode(&reminder)
	if err != nil {
		log.Printf("[API] 更新提醒失败: 无效的请求体 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// 使用服务层更新提醒
	service := getReminderService()
	err = service.UpdateReminder(reminderID, &reminder)
	if err != nil {
		log.Printf("[API] 更新提醒失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to update reminder: " + err.Error()})
		return
	}

	log.Printf("[API] 提醒更新成功: reminder_id=%d", reminderID)
	c.JSON(200, gin.H{
		"message": "Reminder updated successfully",
		"id":      reminderID,
	})
}

// DeleteReminder 删除提醒
func DeleteReminder(c *gin.Context) {
	log.Printf("[API] 收到删除提醒请求: %s %s", c.Request.Method, c.Request.URL.Path)

	reminderIDStr := c.Query("id")
	reminderID, err := strconv.Atoi(reminderIDStr)
	if err != nil {
		log.Printf("[API] 删除提醒失败: 无效的提醒ID - %s", reminderIDStr)
		c.JSON(400, gin.H{"error": "Invalid Reminder ID"})
		return
	}

	// 使用服务层删除提醒
	service := getReminderService()
	err = service.DeleteReminder(reminderID)
	if err != nil {
		log.Printf("[API] 删除提醒失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to delete reminder: " + err.Error()})
		return
	}

	log.Printf("[API] 提醒删除成功: reminder_id=%d", reminderID)
	c.JSON(200, gin.H{
		"message": "Reminder deleted successfully",
		"id":      reminderID,
	})
}

// GetUpcomingReminders 获取即将到期的提醒
func GetUpcomingReminders(c *gin.Context) {
	log.Printf("[API] 收到获取即将到期提醒请求: %s %s", c.Request.Method, c.Request.URL.Path)

	// 获取查询参数
	userIDStr := c.Query("user_id")
	daysStr := c.Query("days_ahead") // 默认为7天

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取即将到期提醒失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	daysAhead := 7 // 默认值
	if daysStr != "" {
		daysAhead, err = strconv.Atoi(daysStr)
		if err != nil {
			log.Printf("[API] 获取即将到期提醒失败: 无效的天数参数 - %s", daysStr)
			c.JSON(400, gin.H{"error": "Invalid Days Ahead parameter"})
			return
		}
	}

	// 使用服务层获取即将到期的提醒
	service := getReminderService()
	reminders, err := service.GetUpcomingReminders(userID, daysAhead)
	if err != nil {
		log.Printf("[API] 获取即将到期提醒失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve upcoming reminders: " + err.Error()})
		return
	}

	log.Printf("[API] 即将到期提醒查询成功: user_id=%d, days_ahead=%d, 共 %d 条提醒", userID, daysAhead, len(reminders))
	c.JSON(200, reminders)
}
