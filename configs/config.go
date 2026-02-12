package configs

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type HealthIndicator struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type HealthCategory struct {
	Name       string            `json:"name"`
	Indicators []HealthIndicator `json:"indicators"`
}

var DefaultHealthConfig = []HealthCategory{
	{
		Name: "生化筛查",
		Indicators: []HealthIndicator{
			{Name: "肌酐", Unit: "μmol/L"},
			{Name: "尿素", Unit: "mmol/L"},
			{Name: "钾", Unit: "mmol/L"},
			{Name: "空腹血糖", Unit: "mmol/L"},
		},
	},
	{
		Name: "血常规",
		Indicators: []HealthIndicator{
			{Name: "红细胞", Unit: "×10^12/L"},
			{Name: "血红蛋白", Unit: "g/L"},
			{Name: "红细胞压积", Unit: "%"},
		},
	},
	{
		Name: "尿常规",
		Indicators: []HealthIndicator{
			{Name: "隐血", Unit: "HPF"},
			{Name: "蛋白质", Unit: "g/L"},
			{Name: "红细胞", Unit: "HPF"},
			{Name: "非鳞状上皮细胞", Unit: ""},
		},
	},
	{
		Name: "尿蛋白、尿素、肌酐测定",
		Indicators: []HealthIndicator{
			{Name: "尿蛋白", Unit: "mg/dL"},
			{Name: "尿蛋白肌酐比值", Unit: "mg/g"},
		},
	},
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./health_app.db")
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)             // 最大打开连接数
	DB.SetMaxIdleConns(25)             // 最大空闲连接数
	DB.SetConnMaxLifetime(time.Hour)   // 连接最大生命周期（1小时）
	DB.SetConnMaxIdleTime(time.Minute) // 连接最大空闲时间（1分钟）

	// 测试连接
	if err = DB.Ping(); err != nil {
		log.Fatal("无法ping数据库:", err)
	}

	// 创建 users 表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		age INTEGER,
		gender TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("[DB] users 表创建成功或已存在")

	// 创建 health_records 表
	createHealthRecordsTable := `
	CREATE TABLE IF NOT EXISTS health_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		category TEXT NOT NULL,
		record_date TEXT NOT NULL,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	_, err = DB.Exec(createHealthRecordsTable)
	if err != nil {
		log.Fatalf("Error creating health_records table: %v", err)
	}
	log.Println("[DB] health_records 表创建成功或已存在")

	// 创建 test_results 表（与 health_records 关联）
	createTestResultsTable := `
	CREATE TABLE IF NOT EXISTS test_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		record_id INTEGER NOT NULL,
		indicator TEXT NOT NULL,
		value REAL NOT NULL,
		unit TEXT NOT NULL,
		reference TEXT,
		abnormal BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (record_id) REFERENCES health_records(id)
	);`
	_, err = DB.Exec(createTestResultsTable)
	if err != nil {
		log.Fatalf("Error creating test_results table: %v", err)
	}
	log.Println("[DB] test_results 表创建成功或已存在")

	// 创建 reminders 表
	createRemindersTable := `
	CREATE TABLE IF NOT EXISTS reminders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		next_date TEXT NOT NULL,
		frequency TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(createRemindersTable)
	if err != nil {
		log.Fatalf("Error creating reminders table: %v", err)
	}
	log.Println("[DB] reminders 表创建成功或已存在")

	// 初始化健康指标配置
	InitHealthConfig()
}

// 初始化健康指标配置
func InitHealthConfig() {
	// 这里可以添加任何需要的初始化逻辑
	// 目前只是预定义了配置，没有需要在数据库中存储的内容
}
