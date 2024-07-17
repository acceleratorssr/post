package service

import (
	"context"
	"post/repository/cache"
)

type LikeService interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
}

type likeService struct {
	cache cache.RedisArticleCache
}

func NewLikeService(cache cache.RedisArticleCache) LikeService {
	return &likeService{
		cache: cache,
	}
}
func (l likeService) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	return l.cache.IncrReadCount(ctx, ObjType, ObjID)
}
