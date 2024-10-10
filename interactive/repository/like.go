package repository

import (
	"context"
	"post/interactive/domain"
	"post/interactive/repository/cache"
	"post/interactive/repository/dao"
)

type LikeRepository interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID uint64) error
	IncrReadCountMany(ctx context.Context, ObjType string, ObjIDs []uint64) error

	IncrLikeCount(ctx context.Context, ObjType string, ObjID, uid uint64) error
	DecrLikeCount(ctx context.Context, ObjType string, ObjID, uid uint64) error
	AddCollectionItem(ctx context.Context, ObjType string, ObjID, uid uint64) error

	UpdateReadCountMany(ctx context.Context, objType string, hmap map[uint64]int64) error

	GetPublishedByBatch(ctx context.Context, ObjType string, list *domain.List) ([]domain.Like, error)
}

type likeRepository struct {
	dao   dao.ArticleLikeDao
	cache cache.ArticleLikeCache
}

func (l *likeRepository) UpdateReadCountMany(ctx context.Context, objType string, hmap map[uint64]int64) error {
	return l.dao.UpdateReadCountMany(ctx, objType, hmap)
}

func (l *likeRepository) GetPublishedByBatch(ctx context.Context, ObjType string, list *domain.List) ([]domain.Like, error) {
	likes, err := l.dao.GetLikeByBatch(ctx, ObjType, list.Limit, list.LastValue, list.OrderBy, list.Desc)
	if err != nil {
		// log
		return nil, err
	}
	return l.toDomain(likes...)
}

func (l *likeRepository) IncrReadCountMany(ctx context.Context, ObjType string, ObjIDs []uint64) error {
	return l.dao.IncrReadCountMany(ctx, ObjType, ObjIDs)
}

func (l *likeRepository) IncrReadCount(ctx context.Context, ObjType string, ObjID uint64) error {
	go func() {
		err := l.cache.IncrReadCount(ctx, ObjType, ObjID)
		if err != nil {
			// log
			return
		}
	}()
	return l.dao.IncrReadCount(ctx, ObjType, ObjID)
}

func (l *likeRepository) AddCollectionItem(ctx context.Context, ObjType string, ObjID, uid uint64) error {
	// 收藏夹访问次数比较低频，不加缓存
	return l.dao.InsertCollection(ctx, ObjType, ObjID, uid)
}

func (l *likeRepository) IncrLikeCount(ctx context.Context, ObjType string, ObjID, uid uint64) error {
	go func() {
		err := l.dao.InSertLike(ctx, ObjType, ObjID, uid)
		if err != nil {
			// log
		}
	}()

	return l.cache.IncrCount(ctx, ObjType, ObjID)
}

func (l *likeRepository) DecrLikeCount(ctx context.Context, ObjType string, ObjID, uid uint64) error {
	go func() {
		err := l.dao.DeleteLike(ctx, ObjType, ObjID, uid)
		if err != nil {
			// log
		}
	}()

	return l.cache.DecrCount(ctx, ObjType, ObjID)
}

func (l *likeRepository) toDomain(art ...dao.Like) ([]domain.Like, error) {
	domainL := make([]domain.Like, 0, len(art))
	for i, _ := range art {
		domainL = append(domainL, domain.Like{
			ID:        art[i].ID,
			LikeCount: art[i].LikeCount,
			Ctime:     art[i].Ctime,
		})
	}
	return domainL, nil
}

func NewLikeRepository(dao dao.ArticleLikeDao, cache cache.ArticleLikeCache) LikeRepository {
	return &likeRepository{
		dao:   dao,
		cache: cache,
	}
}
