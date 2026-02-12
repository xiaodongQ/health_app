package models

import "time"

// HealthRecord 表示一次完整的健康检查记录
type HealthRecord struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	RecordDate time.Time `json:"record_date"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TestResult 表示单个检测指标结果
type TestResult struct {
	ID        int       `json:"id"`
	RecordID  int       `json:"record_id"`
	Indicator string    `json:"indicator"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Reference string    `json:"reference"` // 参考范围
	Abnormal  bool      `json:"abnormal"`  // 是否异常
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HealthRecordWithResults 表示包含检测结果的健康记录
type HealthRecordWithResults struct {
	HealthRecord
	Results []TestResult `json:"results"`
}

// Reminder 表示提醒事项
type Reminder struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`      // medication 或 test
	Title     string    `json:"title"`     // 提醒标题
	NextDate  time.Time `json:"next_date"` // 下次提醒日期
	Frequency string    `json:"frequency"` // 频率
	IsActive  bool      `json:"is_active"` // 是否激活
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
