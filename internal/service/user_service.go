package service

import (
	"errors"
	"health_app/internal/models"
	"health_app/internal/repository"
)

// UserService 用户服务
type UserService struct {
	repo repository.UserRepository
}

// NewUserService 创建新的用户服务实例
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(user *models.User) error {
	// 验证必需字段
	if user.Name == "" || user.Password == "" {
		return errors.New("username and password are required")
	}

	return s.repo.Create(user)
}

// LoginUser 用户登录
func (s *UserService) LoginUser(name, password string) (*models.User, error) {
	// 验证输入参数
	if name == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	return s.repo.FindByNameAndPassword(name, password)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.FindByID(id)
}
