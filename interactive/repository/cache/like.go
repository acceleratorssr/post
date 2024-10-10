package cache

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"post/interactive/domain"
	"strconv"
)

//go:embed cnt.lua
var luaIncrCnt string

type ArticleLikeCache interface {
	IncrReadCount(ctx context.Context, objType string, objID uint64) error
	IncrCount(ctx context.Context, objType string, objID uint64) error
	DecrCount(ctx context.Context, objType string, objID uint64) error
	GetCount(ctx context.Context, objType string, objID uint64) (int64, error)

	GetCountByPrefix(ctx context.Context, prefix string) (map[string]int64, error)
}

type RedisArticleLikeCache struct {
	client redis.Cmdable
}

func (r *RedisArticleLikeCache) GetCountByPrefix(ctx context.Context, prefix string) (map[string]int64, error) {
	result := make(map[string]int64)
	var cursor uint64
	var err error
	pattern := prefix + "*"

	for {
		var batch []string
		batch, cursor, err = r.client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range batch {
			value, err := r.client.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			count, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				continue
			}

			result[key] = count
		}

		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func (r *RedisArticleLikeCache) GetCount(ctx context.Context, objType string, objID uint64) (int64, error) {
	//return r.client.HMGet(ctx, r.keyIncrLikeCount(ObjType, ObjID))
	return r.client.HGet(ctx, domain.KeyIncrLikeCount(objType, objID), "like_cnt").Int64()
}

func (r *RedisArticleLikeCache) DecrCount(ctx context.Context, objType string, objID uint64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{domain.KeyIncrLikeCount(objType, objID)}, "like_cnt", -1).Err()
}

func (r *RedisArticleLikeCache) IncrCount(ctx context.Context, objType string, objID uint64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{domain.KeyIncrLikeCount(objType, objID)}, "like_cnt", 1).Err()
}

func (r *RedisArticleLikeCache) IncrReadCount(ctx context.Context, objType string, objID uint64) error {
	err := r.client.Eval(ctx, luaIncrCnt,
		[]string{domain.KeyIncrReadCount(objType, objID)}, "read_cnt", 1).Err()
	return err
}

func NewRedisArticleLikeCache(client redis.Cmdable) ArticleLikeCache {
	return &RedisArticleLikeCache{
		client: client,
	}
}
