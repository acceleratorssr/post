package cache

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"strconv"
)

//go:embed cnt.lua
var luaIncrCnt string

type ArticleLikeCache interface {
	IncrReadCount(ctx context.Context, objType string, objID uint64) error
	IncrLikeCount(ctx context.Context, objType string, objID uint64) error
	DecrLikeCount(ctx context.Context, objType string, objID uint64) error
	GetLikeCount(ctx context.Context, objType string, objID uint64) (int64, error)
}

type RedisArticleLikeCache struct {
	client redis.Cmdable
}

func NewRedisArticleLikeCache(client redis.Cmdable) ArticleLikeCache {
	return &RedisArticleLikeCache{
		client: client,
	}
}

func (r *RedisArticleLikeCache) GetLikeCount(ctx context.Context, objType string, objID uint64) (int64, error) {
	//return r.client.HMGet(ctx, r.keyIncrLikeCount(ObjType, ObjID))
	return r.client.HGet(ctx, r.keyIncrLikeCount(objType, objID), "like_cnt").Int64()
}

func (r *RedisArticleLikeCache) DecrLikeCount(ctx context.Context, objType string, objID uint64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrLikeCount(objType, objID)}, "like_cnt", -1).Err()
}

func (r *RedisArticleLikeCache) IncrLikeCount(ctx context.Context, objType string, objID uint64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrLikeCount(objType, objID)}, "like_cnt", 1).Err()
}

func (r *RedisArticleLikeCache) IncrReadCount(ctx context.Context, objType string, objID uint64) error {
	err := r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrReadCount(objType, objID)}, "read_cnt", 1).Err()
	return err
}

func (r *RedisArticleLikeCache) keyIncrReadCount(ObjType string, ObjID uint64) string {
	return "article_incr_read_count:" + ObjType + ":" + strconv.FormatUint(ObjID, 10)
}

func (r *RedisArticleLikeCache) keyIncrLikeCount(ObjType string, ObjID uint64) string {
	return "article_incr_Like_count:" + ObjType + ":" + strconv.FormatUint(ObjID, 10)
}
