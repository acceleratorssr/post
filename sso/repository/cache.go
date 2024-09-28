package repository

import (
	"context"
	"post/sso/repository/cache"
	"time"
)

type SSOCache interface {
	SetString(ctx context.Context, key, value string, expiration time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
}

type ssoCache struct {
	cache cache.RedisCache
}

func (s *ssoCache) GetString(ctx context.Context, key string) (string, error) {
	return s.cache.GetString(ctx, key)
}

func (s *ssoCache) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	return s.cache.SetString(ctx, key, value, expiration)
}

func NewSSOCache(cache cache.RedisCache) SSOCache {
	return &ssoCache{
		cache: cache,
	}
}
