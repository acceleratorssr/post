package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserGormDAO interface {
	Insert(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, userInfo *UserInfo) error
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
	user := &User{}
	return user, u.db.WithContext(ctx).Select("NickName").Where("username = ?", username).First(user).Error
}

func (u userGormDAO) Update(ctx context.Context, userInfo *UserInfo) error {
	return u.db.Model(&User{}).WithContext(ctx).Updates(map[string]string{
		"nickname": userInfo.Nickname,
	}).Error
}
