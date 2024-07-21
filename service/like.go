package service

import (
	"context"
	"post/repository"
)

type LikeService interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	Like(ctx context.Context, ObjType string, ObjID, uid int64) error
	UnLike(ctx context.Context, ObjType string, ObjID, uid int64) error
	Collect(ctx context.Context, ObjType string, ObjID, uid int64) error
}

type likeService struct {
	repo repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) LikeService {
	return &likeService{
		repo: repo,
	}
}
func (l *likeService) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	return l.repo.IncrReadCount(ctx, ObjType, ObjID)
}

func (l *likeService) Like(ctx context.Context, ObjType string, ObjID, uid int64) error {
	return l.repo.IncrLikeCount(ctx, ObjType, ObjID, uid)
}

func (l *likeService) UnLike(ctx context.Context, ObjType string, ObjID, uid int64) error {
	return l.repo.DecrLikeCount(ctx, ObjType, ObjID, uid)
}

func (l *likeService) Collect(ctx context.Context, ObjType string, ObjID, uid int64) error {
	return l.repo.AddCollectionItem(ctx, ObjType, ObjID, uid)
}
