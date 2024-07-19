package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLikeCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisLikeCache) IncrReadCnt(ctx context.Context, ObjType string, ObjID int64) error {
	key := r.getKey(ObjType, ObjID)
	return r.client.Incr(ctx, key).Err()
}
