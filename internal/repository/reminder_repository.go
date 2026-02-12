package repository

import (
	"database/sql"
	"health_app/configs"
	"health_app/internal/models"
	"log"
	"time"
)

// ReminderRepository 定义提醒的数据访问接口
type ReminderRepository interface {
	Create(reminder *models.Reminder) (int64, error)
	FindByUserID(userID int, activeOnly bool) ([]*models.Reminder, error)
	FindByID(id int) (*models.Reminder, error)
	Update(id int, reminder *models.Reminder) error
	Delete(id int) error
	GetUpcomingReminders(userID int, daysAhead int) ([]*models.Reminder, error)
}

// SQLiteReminderRepository SQLite实现的提醒仓库
type SQLiteReminderRepository struct {
	db *sql.DB
}

// NewReminderRepository 创建新的提醒仓库实例
func NewReminderRepository() ReminderRepository {
	return &SQLiteReminderRepository{
		db: configs.DB,
	}
}

// Create 创建新提醒
func (r *SQLiteReminderRepository) Create(reminder *models.Reminder) (int64, error) {
	insertReminderQuery := `
		INSERT INTO reminders (user_id, type, title, next_date, frequency, is_active) 
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(insertReminderQuery,
		reminder.UserID, reminder.Type, reminder.Title,
		reminder.NextDate.Format("2006-01-02"), reminder.Frequency, reminder.IsActive)
	if err != nil {
		log.Printf("[Repository] 插入提醒失败: %v", err)
		return 0, err
	}

	reminderID, err := result.LastInsertId()
	if err != nil {
		log.Printf("[Repository] 获取提醒ID失败: %v", err)
		return 0, err
	}

	return reminderID, nil
}

// FindByUserID 根据用户ID查找提醒
func (r *SQLiteReminderRepository) FindByUserID(userID int, activeOnly bool) ([]*models.Reminder, error) {
	var query string
	var args []interface{}

	if activeOnly {
		query = `SELECT id, user_id, type, title, next_date, frequency, is_active, created_at, updated_at 
		         FROM reminders WHERE user_id = ? AND is_active = ? ORDER BY next_date ASC`
		args = []interface{}{userID, true}
	} else {
		query = `SELECT id, user_id, type, title, next_date, frequency, is_active, created_at, updated_at 
		         FROM reminders WHERE user_id = ? ORDER BY next_date ASC`
		args = []interface{}{userID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("[Repository] 查询提醒失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reminders []*models.Reminder
	for rows.Next() {
		var reminder models.Reminder
		var nextDateString string
		err := rows.Scan(
			&reminder.ID, &reminder.UserID, &reminder.Type, &reminder.Title,
			&nextDateString, &reminder.Frequency, &reminder.IsActive,
			&reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			log.Printf("[Repository] 扫描提醒失败: %v", err)
			return nil, err
		}

		// 解析日期字符串
		reminder.NextDate, err = time.Parse("2006-01-02", nextDateString)
		if err != nil {
			log.Printf("[Repository] 解析日期失败: %v", err)
			return nil, err
		}

		reminders = append(reminders, &reminder)
	}

	return reminders, nil
}

// FindByID 根据ID查找提醒
func (r *SQLiteReminderRepository) FindByID(id int) (*models.Reminder, error) {
	var reminder models.Reminder
	var nextDateString string

	query := `SELECT id, user_id, type, title, next_date, frequency, is_active, created_at, updated_at 
	          FROM reminders WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(
		&reminder.ID, &reminder.UserID, &reminder.Type, &reminder.Title,
		&nextDateString, &reminder.Frequency, &reminder.IsActive,
		&reminder.CreatedAt, &reminder.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 提醒不存在，返回nil而不是错误
		}
		log.Printf("[Repository] 查询单个提醒失败: %v", err)
		return nil, err
	}

	// 解析日期字符串
	reminder.NextDate, err = time.Parse("2006-01-02", nextDateString)
	if err != nil {
		log.Printf("[Repository] 解析日期失败: %v", err)
		return nil, err
	}

	return &reminder, nil
}

// Update 更新提醒
func (r *SQLiteReminderRepository) Update(id int, reminder *models.Reminder) error {
	updateReminderQuery := `
		UPDATE reminders 
		SET type = ?, title = ?, next_date = ?, frequency = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`
	_, err := r.db.Exec(updateReminderQuery,
		reminder.Type, reminder.Title, reminder.NextDate.Format("2006-01-02"),
		reminder.Frequency, reminder.IsActive, id)
	if err != nil {
		log.Printf("[Repository] 更新提醒失败: %v", err)
		return err
	}

	return nil
}

// Delete 删除提醒
func (r *SQLiteReminderRepository) Delete(id int) error {
	deleteReminderQuery := "DELETE FROM reminders WHERE id = ?"
	result, err := r.db.Exec(deleteReminderQuery, id)
	if err != nil {
		log.Printf("[Repository] 删除提醒失败: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[Repository] 获取影响行数失败: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // 提醒不存在
	}

	return nil
}

// GetUpcomingReminders 获取即将到期的提醒
func (r *SQLiteReminderRepository) GetUpcomingReminders(userID int, daysAhead int) ([]*models.Reminder, error) {
	// 计算截止日期
	cutoffDate := time.Now().AddDate(0, 0, daysAhead)

	// 查询即将到期的提醒
	query := `
		SELECT id, user_id, type, title, next_date, frequency, is_active, created_at, updated_at 
		FROM reminders 
		WHERE user_id = ? 
		AND is_active = ? 
		AND next_date <= ?
		ORDER BY next_date ASC`

	rows, err := r.db.Query(query, userID, true, cutoffDate.Format("2006-01-02"))
	if err != nil {
		log.Printf("[Repository] 查询即将到期提醒失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reminders []*models.Reminder
	for rows.Next() {
		var reminder models.Reminder
		var nextDateString string
		err := rows.Scan(
			&reminder.ID, &reminder.UserID, &reminder.Type, &reminder.Title,
			&nextDateString, &reminder.Frequency, &reminder.IsActive,
			&reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			log.Printf("[Repository] 扫描即将到期提醒失败: %v", err)
			return nil, err
		}

		// 解析日期字符串
		reminder.NextDate, err = time.Parse("2006-01-02", nextDateString)
		if err != nil {
			log.Printf("[Repository] 解析日期失败: %v", err)
			return nil, err
		}

		reminders = append(reminders, &reminder)
	}

	return reminders, nil
}
