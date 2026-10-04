package repository

import (
	"database/sql"
	"health_app/configs"
	"health_app/internal/models"
	"log"
	"time"
)

// HealthRecordRepository 定义健康记录的数据访问接口
type HealthRecordRepository interface {
	Create(record *models.HealthRecord, results []models.TestResult) (int64, error)
	FindByUserID(userID int) ([]*models.HealthRecordWithResults, error)
	FindByID(id int) (*models.HealthRecordWithResults, error)
	Update(id int, record *models.HealthRecord, results []models.TestResult) error
	Delete(id int) error
	FindByDateAndUser(userID int, date string) (*models.HealthRecordWithResults, error)
	GetTrendData(userID int, indicator string, unit string) ([]map[string]interface{}, error)
	GetAllIndicators(userID int) ([]string, error)
}

// SQLiteHealthRecordRepository SQLite实现的健康记录仓库
type SQLiteHealthRecordRepository struct {
	db *sql.DB
}

// NewHealthRecordRepository 创建新的健康记录仓库实例
func NewHealthRecordRepository() HealthRecordRepository {
	return &SQLiteHealthRecordRepository{
		db: configs.DB,
	}
}

// Create 创建新的健康记录及其检测结果
func (r *SQLiteHealthRecordRepository) Create(record *models.HealthRecord, results []models.TestResult) (int64, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return 0, sql.ErrConnDone
	}

	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[Repository] 开启事务失败: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	// 插入健康记录
	insertRecordQuery := `
		INSERT INTO health_records (user_id, record_date, category) 
		VALUES (?, ?, ?)`
	result, err := tx.Exec(insertRecordQuery, record.UserID, record.RecordDate.Format("2006-01-02"), record.Category)
	if err != nil {
		log.Printf("[Repository] 插入健康记录失败: %v", err)
		return 0, err
	}

	recordID, err := result.LastInsertId()
	if err != nil {
		log.Printf("[Repository] 获取记录ID失败: %v", err)
		return 0, err
	}

	// 插入所有检测结果
	for _, result := range results {
		insertResultQuery := `
			INSERT INTO test_results (record_id, indicator, value, unit, reference, abnormal) 
			VALUES (?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(insertResultQuery, recordID, result.Indicator, result.Value, result.Unit, result.Reference, result.Abnormal)
		if err != nil {
			log.Printf("[Repository] 插入检测结果失败: %v", err)
			return 0, err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Printf("[Repository] 提交事务失败: %v", err)
		return 0, err
	}

	return recordID, nil
}

// FindByUserID 根据用户ID查找所有健康记录
func (r *SQLiteHealthRecordRepository) FindByUserID(userID int) ([]*models.HealthRecordWithResults, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	// 使用 JOIN 一次性查询所有记录及关联的检测结果
	query := `
		SELECT
			hr.id, hr.user_id, hr.record_date, hr.category, hr.created_at, hr.updated_at,
			tr.id, tr.indicator, tr.value, tr.unit, tr.reference, tr.abnormal, tr.created_at, tr.updated_at
		FROM health_records hr
		LEFT JOIN test_results tr ON hr.id = tr.record_id
		WHERE hr.user_id = ?
		ORDER BY hr.record_date DESC, hr.id, tr.indicator`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("[Repository] 查询健康记录失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 用 map 做去重，用 slice 保持 JOIN 结果的顺序（已按 date DESC 排序）
	recordMap := make(map[int]*models.HealthRecordWithResults)
	var records []*models.HealthRecordWithResults
	for rows.Next() {
		var hr models.HealthRecord
		var recordDate string
		var trID int
		var indicator string
		var value float64
		var unit, reference string
		var abnormal bool
		var trCreatedAt, trUpdatedAt time.Time

		err := rows.Scan(
			&hr.ID, &hr.UserID, &recordDate, &hr.Category, &hr.CreatedAt, &hr.UpdatedAt,
			&trID, &indicator, &value, &unit, &reference, &abnormal, &trCreatedAt, &trUpdatedAt,
		)
		if err != nil {
			log.Printf("[Repository] 扫描记录失败: %v", err)
			return nil, err
		}

		hr.RecordDate, _ = time.Parse("2006-01-02", recordDate)

		if _, exists := recordMap[hr.ID]; !exists {
			rec := &models.HealthRecordWithResults{
				HealthRecord: hr,
				Results:      []models.TestResult{},
			}
			recordMap[hr.ID] = rec
			records = append(records, rec) // 按 JOIN 结果顺序追加
		}

		// 如果有检测结果则追加
		if trID > 0 {
			recordMap[hr.ID].Results = append(recordMap[hr.ID].Results, models.TestResult{
				ID:        trID,
				RecordID:  hr.ID,
				Indicator: indicator,
				Value:     value,
				Unit:      unit,
				Reference: reference,
				Abnormal:  abnormal,
				CreatedAt: trCreatedAt,
				UpdatedAt: trUpdatedAt,
			})
		}
	}
	return records, nil
}

// FindByID 根据ID查找健康记录
func (r *SQLiteHealthRecordRepository) FindByID(id int) (*models.HealthRecordWithResults, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	// 查询健康记录
	recordQuery := `
		SELECT id, user_id, record_date, category, created_at, updated_at 
		FROM health_records 
		WHERE id = ?`
	row := r.db.QueryRow(recordQuery, id)

	var record models.HealthRecord
	var recordDate string

	err := row.Scan(&record.ID, &record.UserID, &recordDate, &record.Category, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		log.Printf("[Repository] 查询单个健康记录失败: %v", err)
		return nil, err
	}

	// 解析日期字符串
	record.RecordDate, err = time.Parse("2006-01-02", recordDate)
	if err != nil {
		log.Printf("[Repository] 解析日期失败: %v", err)
		return nil, err
	}

	// 获取该记录的检测结果
	results, err := r.getTestResultsByRecordID(record.ID)
	if err != nil {
		log.Printf("[Repository] 获取检测结果失败: %v", err)
		return nil, err
	}

	recordWithResults := &models.HealthRecordWithResults{
		HealthRecord: record,
		Results:      results,
	}

	return recordWithResults, nil
}

// Update 更新健康记录
func (r *SQLiteHealthRecordRepository) Update(id int, record *models.HealthRecord, results []models.TestResult) error {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return sql.ErrConnDone
	}

	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[Repository] 开启事务失败: %v", err)
		return err
	}
	defer tx.Rollback()

	// 更新健康记录基本信息
	updateRecordQuery := `
		UPDATE health_records 
		SET record_date = ?, category = ? 
		WHERE id = ?`
	_, err = tx.Exec(updateRecordQuery, record.RecordDate.Format("2006-01-02"), record.Category, id)
	if err != nil {
		log.Printf("[Repository] 更新记录信息失败: %v", err)
		return err
	}

	// 删除现有的检测结果
	deleteResultsQuery := `DELETE FROM test_results WHERE record_id = ?`
	_, err = tx.Exec(deleteResultsQuery, id)
	if err != nil {
		log.Printf("[Repository] 删除现有结果失败: %v", err)
		return err
	}

	// 插入新的检测结果
	for _, result := range results {
		insertResultQuery := `
			INSERT INTO test_results (record_id, indicator, value, unit, reference, abnormal) 
			VALUES (?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(insertResultQuery, id, result.Indicator, result.Value, result.Unit, result.Reference, result.Abnormal)
		if err != nil {
			log.Printf("[Repository] 插入检测结果失败: %v", err)
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Printf("[Repository] 提交事务失败: %v", err)
		return err
	}

	return nil
}

// Delete 删除健康记录
func (r *SQLiteHealthRecordRepository) Delete(id int) error {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return sql.ErrConnDone
	}

	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[Repository] 开启事务失败: %v", err)
		return err
	}
	defer tx.Rollback()

	// 删除检测结果
	deleteResultsQuery := `DELETE FROM test_results WHERE record_id = ?`
	_, err = tx.Exec(deleteResultsQuery, id)
	if err != nil {
		log.Printf("[Repository] 删除检测结果失败: %v", err)
		return err
	}

	// 删除健康记录
	deleteRecordQuery := `DELETE FROM health_records WHERE id = ?`
	_, err = tx.Exec(deleteRecordQuery, id)
	if err != nil {
		log.Printf("[Repository] 删除健康记录失败: %v", err)
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Printf("[Repository] 提交事务失败: %v", err)
		return err
	}

	return nil
}

// FindByDateAndUser 根据用户ID和日期查找健康记录
func (r *SQLiteHealthRecordRepository) FindByDateAndUser(userID int, date string) (*models.HealthRecordWithResults, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	// 查询指定用户和日期的健康记录
	recordQuery := `
		SELECT id, user_id, record_date, category, created_at, updated_at 
		FROM health_records 
		WHERE user_id = ? AND record_date = ?`
	row := r.db.QueryRow(recordQuery, userID, date)

	var record models.HealthRecord
	var recordDate string

	err := row.Scan(&record.ID, &record.UserID, &recordDate, &record.Category, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有找到记录，返回nil而不是错误
		}
		log.Printf("[Repository] 查询指定日期记录失败: %v", err)
		return nil, err
	}

	// 解析日期字符串
	record.RecordDate, err = time.Parse("2006-01-02", recordDate)
	if err != nil {
		log.Printf("[Repository] 解析日期失败: %v", err)
		return nil, err
	}

	// 获取该记录的检测结果
	results, err := r.getTestResultsByRecordID(record.ID)
	if err != nil {
		log.Printf("[Repository] 获取检测结果失败: %v", err)
		return nil, err
	}

	recordWithResults := &models.HealthRecordWithResults{
		HealthRecord: record,
		Results:      results,
	}

	return recordWithResults, nil
}

// GetTrendData 获取特定指标的趋势数据
func (r *SQLiteHealthRecordRepository) GetTrendData(userID int, indicator string, unit string) ([]map[string]interface{}, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	var trendQuery string
	var rows *sql.Rows
	var err error

	if unit != "" {
		// 如果指定了单位，则精确匹配指标名称和单位
		trendQuery = `
			SELECT hr.record_date, tr.value, tr.unit, tr.abnormal
			FROM test_results tr
			JOIN health_records hr ON tr.record_id = hr.id
			WHERE hr.user_id = ? AND tr.indicator = ? AND tr.unit = ?
			ORDER BY hr.record_date ASC`
		rows, err = r.db.Query(trendQuery, userID, indicator, unit)
	} else {
		// 如果没有指定单位，则只按指标名称匹配（兼容旧的调用方式）
		trendQuery = `
			SELECT hr.record_date, tr.value, tr.unit, tr.abnormal
			FROM test_results tr
			JOIN health_records hr ON tr.record_id = hr.id
			WHERE hr.user_id = ? AND tr.indicator = ?
			ORDER BY hr.record_date ASC`
		rows, err = r.db.Query(trendQuery, userID, indicator)
	}

	if err != nil {
		log.Printf("[Repository] 查询趋势数据失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var trends []map[string]interface{}
	for rows.Next() {
		var recordDate string
		var value float64
		var dbUnit string
		var abnormal bool

		err := rows.Scan(&recordDate, &value, &dbUnit, &abnormal)
		if err != nil {
			log.Printf("[Repository] 扫描趋势数据失败: %v", err)
			return nil, err
		}

		trend := map[string]interface{}{
			"date":     recordDate,
			"value":    value,
			"unit":     dbUnit,
			"abnormal": abnormal,
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetAllIndicators 获取所有指标名称
func (r *SQLiteHealthRecordRepository) GetAllIndicators(userID int) ([]string, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	// 查询用户的所有指标名称
	indicatorsQuery := `
		SELECT DISTINCT tr.indicator
		FROM test_results tr
		JOIN health_records hr ON tr.record_id = hr.id
		WHERE hr.user_id = ?
		ORDER BY tr.indicator`
	rows, err := r.db.Query(indicatorsQuery, userID)
	if err != nil {
		log.Printf("[Repository] 查询指标名称失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var indicators []string
	for rows.Next() {
		var indicator string
		err := rows.Scan(&indicator)
		if err != nil {
			log.Printf("[Repository] 扫描指标名称失败: %v", err)
			return nil, err
		}
		indicators = append(indicators, indicator)
	}

	return indicators, nil
}

// getTestResultsByRecordID 根据记录ID获取检测结果
func (r *SQLiteHealthRecordRepository) getTestResultsByRecordID(recordID int) ([]models.TestResult, error) {
	if r.db == nil {
		log.Printf("[Repository] 数据库连接未初始化")
		return nil, sql.ErrConnDone
	}

	resultsQuery := `
		SELECT id, indicator, value, unit, reference, abnormal, created_at, updated_at
		FROM test_results
		WHERE record_id = ?
		ORDER BY indicator`
	resultRows, err := r.db.Query(resultsQuery, recordID)
	if err != nil {
		log.Printf("[Repository] 查询检测结果失败: %v", err)
		return nil, err
	}
	defer resultRows.Close()

	var results []models.TestResult
	for resultRows.Next() {
		var result models.TestResult
		err := resultRows.Scan(&result.ID, &result.Indicator, &result.Value, &result.Unit, &result.Reference, &result.Abnormal, &result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			log.Printf("[Repository] 扫描检测结果失败: %v", err)
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
