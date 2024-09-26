package dao

import (
	"context"
	"gorm.io/gorm"
)

type SSOGormDAO interface {
	Insert(ctx context.Context, user *User) error
	QueryByUsername(ctx context.Context, username string) (*User, error)
}

type ssoGormDAO struct {
	db *gorm.DB
}

func (s *ssoGormDAO) QueryByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	return user, s.db.WithContext(ctx).Where("username = ?", username).First(user).Error
}

func (s *ssoGormDAO) Insert(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func NewSSOGormDAO(db *gorm.DB) SSOGormDAO {
	return &ssoGormDAO{
		db: db,
	}
}
