package repository

import (
	"context"
	"post/sso/repository/cache"
	"time"
)

type SSOCache interface {
	SetUsernameAndKey(ctx context.Context, username, TotpSecret string, expiration time.Duration) error
	GetUsernameAndKey(ctx context.Context, username string) (string, error)
}

type ssoCache struct {
	cache cache.RedisCache
}

func (s *ssoCache) GetUsernameAndKey(ctx context.Context, username string) (string, error) {
	return s.cache.GetString(ctx, username)
}

func (s *ssoCache) SetUsernameAndKey(ctx context.Context, username, TotpSecret string, expiration time.Duration) error {
	return s.cache.SetString(ctx, username, TotpSecret, expiration)
}

func NewSSOCache(cache cache.RedisCache) SSOCache {
	return &ssoCache{
		cache: cache,
	}
}
