package repository

import (
	"database/sql"
	"health_app/configs"
	"log"
)

type UserSettingsRepository interface {
	Get(userID int, key string) (string, error)
	Set(userID int, key string, value string) error
}

type SQLiteUserSettingsRepository struct {
	db *sql.DB
}

func NewUserSettingsRepository() UserSettingsRepository {
	return &SQLiteUserSettingsRepository{db: configs.DB}
}

func (r *SQLiteUserSettingsRepository) Get(userID int, key string) (string, error) {
	query := `SELECT setting_value FROM user_settings WHERE user_id = ? AND setting_key = ?`
	var value string
	err := r.db.QueryRow(query, userID, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		log.Printf("[Repository] Get user_settings failed: %v", err)
		return "", err
	}
	return value, nil
}

func (r *SQLiteUserSettingsRepository) Set(userID int, key string, value string) error {
	query := `INSERT INTO user_settings (user_id, setting_key, setting_value) VALUES (?, ?, ?)
		ON CONFLICT(user_id, setting_key) DO UPDATE SET setting_value = excluded.setting_value, updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, userID, key, value)
	if err != nil {
		log.Printf("[Repository] Set user_settings failed: %v", err)
		return err
	}
	return nil
}
