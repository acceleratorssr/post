package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache interface {
	SetString(ctx context.Context, key string, value string, expiration time.Duration) error
	GetString(ctx context.Context, username string) (string, error)
}

type redisCache struct {
	client redis.Cmdable
}

func (r *redisCache) GetString(ctx context.Context, username string) (string, error) {
	return r.client.Get(ctx, username).Result()
}

func (r *redisCache) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func NewRedisCache(client redis.Cmdable) RedisCache {
	return &redisCache{
		client: client,
	}
}
