package service

import (
	"context"
	"post/interactive/domain"
	"post/interactive/repository"
)

// LikeService //go:generate 移动到makefile里了
type LikeService interface {
	IncrReadCount(ctx context.Context, objType string, objID uint64) error
	Like(ctx context.Context, objType string, objID, uid uint64) error
	UnLike(ctx context.Context, objType string, objID, uid uint64) error
	Collect(ctx context.Context, objType string, objID, uid uint64) error

	GetListBatchOfLikes(ctx context.Context, objType string, offset, limit int, now int64) ([]domain.Like, error)
}

type likeService struct {
	repo repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) LikeService {
	return &likeService{
		repo: repo,
	}
}

func (l *likeService) GetListBatchOfLikes(ctx context.Context, objType string, offset, limit int, now int64) ([]domain.Like, error) {
	// TODO 可考虑排除老旧的数据，提高速度
	return l.repo.GetListAllOfLikes(ctx, objType, offset, limit, now)
}

func (l *likeService) IncrReadCount(ctx context.Context, objType string, objID uint64) error {
	return l.repo.IncrReadCount(ctx, objType, objID)
}

// Like 没控制重复点赞
func (l *likeService) Like(ctx context.Context, objType string, objID, uid uint64) error {
	return l.repo.IncrLikeCount(ctx, objType, objID, uid)
}

func (l *likeService) UnLike(ctx context.Context, objType string, objID, uid uint64) error {
	return l.repo.DecrLikeCount(ctx, objType, objID, uid)
}

func (l *likeService) Collect(ctx context.Context, objType string, objID, uid uint64) error {
	return l.repo.AddCollectionItem(ctx, objType, objID, uid)
}
