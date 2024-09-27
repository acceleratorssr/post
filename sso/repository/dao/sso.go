package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SSOGormDAO interface {
	Upsert(ctx context.Context, user *User) error
	QueryByUsername(ctx context.Context, username string) (*User, error)
	QueryTotpSecretByUsername(ctx context.Context, username string) (string, error)
}

type ssoGormDAO struct {
	db *gorm.DB
}

func (s *ssoGormDAO) QueryTotpSecretByUsername(ctx context.Context, username string) (string, error) {
	user := &User{}
	return user.TotpSecret, s.db.WithContext(ctx).Select("totp_secret").Where("username = ?", username).First(&user).Error
}

func (s *ssoGormDAO) QueryByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	return user, s.db.WithContext(ctx).Where("username = ?", username).First(user).Error
}

func (s *ssoGormDAO) Upsert(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"nickname": user.Nickname,
			"utime":    user.Utime,
		}),
	}).Create(user).Error
}

func NewSSOGormDAO(db *gorm.DB) SSOGormDAO {
	return &ssoGormDAO{
		db: db,
	}
}
