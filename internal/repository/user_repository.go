package repository

import (
	"database/sql"
	"health_app/configs"
	"health_app/internal/models"
	"log"
)

// UserRepository 定义用户的数据访问接口
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByNameAndPassword(name, password string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

// SQLiteUserRepository SQLite实现的用户仓库
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建新的用户仓库实例
func NewUserRepository() UserRepository {
	return &SQLiteUserRepository{
		db: configs.DB,
	}
}

// Create 创建新用户
func (r *SQLiteUserRepository) Create(user *models.User) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("[Repository] 创建用户失败: %v", err)
		return err
	}
	return nil
}

// FindByID 根据ID查找用户
func (r *SQLiteUserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, password FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用户不存在，返回nil而不是错误
		}
		log.Printf("[Repository] 查询用户失败: %v", err)
		return nil, err
	}
	return &user, nil
}

// FindByNameAndPassword 根据用户名和密码查找用户
func (r *SQLiteUserRepository) FindByNameAndPassword(name, password string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, password FROM users WHERE name = ? AND password = ?"
	err := r.db.QueryRow(query, name, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用户不存在，返回nil而不是错误
		}
		log.Printf("[Repository] 查询用户失败: %v", err)
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *SQLiteUserRepository) Update(user *models.User) error {
	query := "UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("[Repository] 更新用户失败: %v", err)
		return err
	}
	return nil
}

// Delete 删除用户
func (r *SQLiteUserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("[Repository] 删除用户失败: %v", err)
		return err
	}
	return nil
}
