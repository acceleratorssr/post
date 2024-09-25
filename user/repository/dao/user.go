package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserGormDAO interface {
	Insert(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type userGormDAO struct {
	db *gorm.DB
}

func NewUserGormDAO(db *gorm.DB) UserGormDAO {
	return &userGormDAO{
		db: db,
	}
}

func (u userGormDAO) Insert(ctx context.Context, user *User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u userGormDAO) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user *User
	return user, u.db.WithContext(ctx).Where("username = ?", username).First(user).Error
}

func (u userGormDAO) Update(ctx context.Context, user *User) error {
	return u.db.WithContext(ctx).Updates(user).Error
}
