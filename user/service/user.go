package service

import (
	"context"
	"post/user/domain"
	"post/user/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserInfoByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, userInfo *domain.UserInfo) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepository: repo,
	}
}

// CreateUser 创建普通用户
func (u *userService) CreateUser(ctx context.Context, user *domain.User) error {
	return u.userRepository.Create(ctx, user)
}

// GetUserInfoByUsername 根据用户名获取用户信息
func (u *userService) GetUserInfoByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.userRepository.GetByUsername(ctx, username)
}

// UpdateUser 更新用户密码或昵称
func (u *userService) UpdateUser(ctx context.Context, userInfo *domain.UserInfo) error {
	return u.userRepository.Update(ctx, userInfo)
}
