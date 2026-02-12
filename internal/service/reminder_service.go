package service

import (
	"errors"
	"health_app/internal/models"
	"health_app/internal/repository"
)

// ReminderService 提醒服务
type ReminderService struct {
	repo repository.ReminderRepository
}

// NewReminderService 创建新的提醒服务实例
func NewReminderService(repo repository.ReminderRepository) *ReminderService {
	return &ReminderService{
		repo: repo,
	}
}

// CreateReminder 创建提醒
func (s *ReminderService) CreateReminder(reminder *models.Reminder) (int64, error) {
	// 验证必需字段
	if reminder.UserID == 0 || reminder.Type == "" || reminder.Title == "" {
		return 0, errors.New("user ID, type, and title are required")
	}

	return s.repo.Create(reminder)
}

// GetReminders 获取用户的提醒列表
func (s *ReminderService) GetReminders(userID int, activeOnly bool) ([]*models.Reminder, error) {
	return s.repo.FindByUserID(userID, activeOnly)
}

// GetReminderByID 根据ID获取单个提醒
func (s *ReminderService) GetReminderByID(id int) (*models.Reminder, error) {
	return s.repo.FindByID(id)
}

// UpdateReminder 更新提醒
func (s *ReminderService) UpdateReminder(id int, reminder *models.Reminder) error {
	// 验证必需字段
	if reminder.Type == "" || reminder.Title == "" {
		return errors.New("type and title are required")
	}

	return s.repo.Update(id, reminder)
}

// DeleteReminder 删除提醒
func (s *ReminderService) DeleteReminder(id int) error {
	return s.repo.Delete(id)
}

// GetUpcomingReminders 获取即将到来的提醒
func (s *ReminderService) GetUpcomingReminders(userID int, daysAhead int) ([]*models.Reminder, error) {
	return s.repo.GetUpcomingReminders(userID, daysAhead)
}
