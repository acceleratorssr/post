package repository

import (
	"context"
	"post/sso/domain"
	"post/sso/repository/dao"
)

type SSORepository interface {
	SaveUserInfo(ctx context.Context, user *domain.User) error
	GetInfoByUsername(ctx context.Context, username string) (*domain.User, error)
}

type ssoRepository struct {
	dao dao.SSOGormDAO
}

func (s *ssoRepository) GetInfoByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.dao.QueryByUsername(ctx, username)
	return s.ToDomain(user), err
}

func (s *ssoRepository) SaveUserInfo(ctx context.Context, user *domain.User) error {
	return s.dao.Insert(ctx, s.ToDao(user))
}

func (s *ssoRepository) ToDao(user *domain.User) *dao.User {
	return &dao.User{
		Password:   user.Password,
		Username:   user.Username,
		QrcodeURL:  user.QrcodeURL,
		TotpSecret: user.TotpSecret,
		UserAgent:  user.UserAgent,
		Nickname:   user.Nickname,
	}
}

func (s *ssoRepository) ToDomain(user *dao.User) *domain.User {
	return &domain.User{
		ID:         user.ID,
		Password:   user.Password,
		Username:   user.Username,
		QrcodeURL:  user.QrcodeURL,
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
