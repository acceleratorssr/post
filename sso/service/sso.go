package service

import (
	"context"
	"post/sso/domain"
	"post/sso/repository"
)

type AuthUserService interface {
	GetInfoByUsername(ctx context.Context, username string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User, utime int64) error
	GetTotpSecret(ctx context.Context, username string) (string, error)
}

type ssoService struct {
	repo repository.SSORepository
}

func (a *ssoService) GetTotpSecret(ctx context.Context, username string) (string, error) {
	return a.repo.GetTotpSecret(ctx, username)
}

func (a *ssoService) SaveUser(ctx context.Context, user *domain.User, utime int64) error {
	return a.repo.SaveUserInfo(ctx, user, utime)
}

func (a *ssoService) GetInfoByUsername(ctx context.Context, username string) (*domain.User, error) {
	return a.repo.GetInfoByUsername(ctx, username)
}

func NewAuthService(repo repository.SSORepository) AuthUserService {
	return &ssoService{
		repo: repo,
	}
}
