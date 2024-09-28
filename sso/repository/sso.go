package repository

import (
	"context"
	"post/sso/domain"
	"post/sso/repository/dao"
)

type SSORepository interface {
	SaveUserInfo(ctx context.Context, user *domain.User, utime int64) error
	GetInfoByUsername(ctx context.Context, username string) (*domain.User, error)
	GetTotpSecret(ctx context.Context, username string) (string, error)
	FindUsername(ctx context.Context, username string) bool
}

type ssoRepository struct {
	dao dao.SSOGormDAO
}

func (s *ssoRepository) FindUsername(ctx context.Context, username string) bool {
	return s.dao.UsernameExistOrNot(ctx, username)
}

func (s *ssoRepository) GetTotpSecret(ctx context.Context, username string) (string, error) {
	return s.dao.QueryTotpSecret(ctx, username)
}

func (s *ssoRepository) GetInfoByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.dao.Query(ctx, username)
	return s.ToDomain(user), err
}

func (s *ssoRepository) SaveUserInfo(ctx context.Context, user *domain.User, utime int64) error {
	return s.dao.Upsert(ctx, s.ToDao(user, utime))
}

func (s *ssoRepository) ToDao(user *domain.User, now int64) *dao.User {
	return &dao.User{
		Password:   user.Password,
		Username:   user.Username,
		TotpSecret: user.TotpSecret,
		UserAgent:  user.UserAgent,
		Nickname:   user.Nickname,
		Utime:      now,
		Ctime:      now,
	}
}

func (s *ssoRepository) ToDomain(user *dao.User) *domain.User {
	return &domain.User{
		Password:   user.Password,
		Username:   user.Username,
		TotpSecret: user.TotpSecret,
		UserAgent:  user.UserAgent,
		Nickname:   user.Nickname,
	}
}

func NewSSOGormRepository(db dao.SSOGormDAO) SSORepository {
	return &ssoRepository{
		dao: db,
	}
}
