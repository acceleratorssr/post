package repository

import (
	"context"
	"post/repository/cache"
	"post/repository/dao"
)

type LikeRepository interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	IncrLikeCount(ctx context.Context, ObjType string, ObjID, uid int64) error
	DecrLikeCount(ctx context.Context, ObjType string, ObjID, uid int64) error
}

type likeRepository struct {
	dao   dao.ArticleLikeDao
	cache cache.RedisArticleCache
}

func NewLikeRepository(dao dao.ArticleLikeDao, cache cache.RedisArticleCache) LikeRepository {
	return &likeRepository{
		dao:   dao,
		cache: cache,
	}
}
func (l *likeRepository) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	go l.cache.IncrReadCount(ctx, ObjType, ObjID)
	return l.dao.IncrReadCount(ctx, ObjType, ObjID)
}

func (l *likeRepository) IncrLikeCount(ctx context.Context, ObjType string, ObjID, uid int64) error {
	err := l.dao.InSertLike(ctx, ObjType, ObjID, uid)
	if err != nil {
		return err
	}

	return l.cache.IncrLikeCount(ctx, ObjType, ObjID)
}

func (l *likeRepository) DecrLikeCount(ctx context.Context, ObjType string, ObjID, uid int64) error {
	err := l.dao.DeleteLike(ctx, ObjType, ObjID, uid)
	if err != nil {
		return err
	}

	return l.cache.DecrLikeCount(ctx, ObjType, ObjID)
}
