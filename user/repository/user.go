package repository

import (
	"context"
	"post/user/domain"
	"post/user/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
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
	return u.dao.Insert(ctx, u.toDao(user))
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	res, err := u.dao.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return u.toDomain(res), err
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) error {
	return u.dao.Update(ctx, u.toDao(user))
}

func (u *userRepository) toDomain(user *dao.User) *domain.User {
	return &domain.User{
		ID:          user.ID,
		Nickname:    user.Nickname,
		Password:    user.Password,
		Permissions: user.Permissions,
		Username:    user.Username,
	}
}

func (u *userRepository) toDao(user *domain.User) *dao.User {
	return &dao.User{
		ID:          user.ID,
		Nickname:    user.Nickname,
		Password:    user.Password,
		Permissions: user.Permissions,
		Username:    user.Username,
	}
}
