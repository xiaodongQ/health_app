package service

import (
	"errors"
	"health_app/internal/models"
	"health_app/internal/repository"
)

// HealthRecordService 健康记录服务
type HealthRecordService struct {
	repo repository.HealthRecordRepository
}

// NewHealthRecordService 创建新的健康记录服务实例
func NewHealthRecordService(repo repository.HealthRecordRepository) *HealthRecordService {
	return &HealthRecordService{
		repo: repo,
	}
}

// AddHealthRecord 添加健康记录
func (s *HealthRecordService) AddHealthRecord(record *models.HealthRecord, results []models.TestResult) (int64, error) {
	// 验证必需字段
	if record.UserID == 0 || record.RecordDate.IsZero() || record.Category == "" {
		return 0, errors.New("user ID, record date, and category are required")
	}

	return s.repo.Create(record, results)
}

// GetHealthRecords 获取用户的所有健康记录
func (s *HealthRecordService) GetHealthRecords(userID int) ([]*models.HealthRecordWithResults, error) {
	return s.repo.FindByUserID(userID)
}

// GetHealthRecordByID 根据ID获取单个健康记录
func (s *HealthRecordService) GetHealthRecordByID(id int) (*models.HealthRecordWithResults, error) {
	return s.repo.FindByID(id)
}

// UpdateHealthRecord 更新健康记录
func (s *HealthRecordService) UpdateHealthRecord(id int, record *models.HealthRecord, results []models.TestResult) error {
	// 验证必需字段
	if record.RecordDate.IsZero() || record.Category == "" {
		return errors.New("record date and category are required")
	}

	return s.repo.Update(id, record, results)
}

// DeleteHealthRecord 删除健康记录
func (s *HealthRecordService) DeleteHealthRecord(id int) error {
	return s.repo.Delete(id)
}

// GetHealthRecordByDate 根据用户和日期获取健康记录
func (s *HealthRecordService) GetHealthRecordByDate(userID int, date string) (*models.HealthRecordWithResults, error) {
	return s.repo.FindByDateAndUser(userID, date)
}

// GetTrendData 获取特定指标的趋势数据
func (s *HealthRecordService) GetTrendData(userID int, indicator string, unit string) ([]map[string]interface{}, error) {
	return s.repo.GetTrendData(userID, indicator, unit)
}

// GetAllIndicators 获取所有指标名称
func (s *HealthRecordService) GetAllIndicators(userID int) ([]string, error) {
	return s.repo.GetAllIndicators(userID)
}
