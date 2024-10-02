package repository

import (
	"context"
	"post/user/domain"
	"post/user/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, userInfo *domain.UserInfo) error
}

type userRepository struct {
	dao dao.UserGormDAO
}

func NewUserRepository(dao dao.UserGormDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	return u.dao.Insert(ctx, u.toUserDao(user))
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	res, err := u.dao.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return u.toUserDomain(res), err
}

func (u *userRepository) Update(ctx context.Context, userInfo *domain.UserInfo) error {
	return u.dao.Update(ctx, (*dao.UserInfo)(userInfo))
}

func (u *userRepository) toUserDomain(user *dao.User) *domain.User {
	return &domain.User{
		ID:       user.ID,
		Nickname: user.Nickname,
		Username: user.Username,
	}
}

func (u *userRepository) toUserInfoDomain(user *dao.UserInfo) *domain.UserInfo {
	return &domain.UserInfo{
		Nickname: user.Nickname,
	}
}

func (u *userRepository) toUserDao(user *domain.User) *dao.User {
	return &dao.User{
		ID:       user.ID,
		Nickname: user.Nickname,
		Username: user.Username,
	}
}

func (u *userRepository) toUserInfoDao(user *domain.UserInfo) *dao.UserInfo {
	return &dao.UserInfo{
		Nickname: user.Nickname,
	}
}
