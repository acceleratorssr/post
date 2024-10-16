package service

import (
	"context"
	"post/interactive/domain"
	"post/interactive/events"
	"post/interactive/repository"
	"strconv"
	"time"
)

// LikeService //go:generate 移动到makefile里了
type LikeService interface {
	IncrReadCount(ctx context.Context, objType string, objID uint64) error
	IncrReadCountMany(ctx context.Context, objType string, objIDs []uint64) error
	Like(ctx context.Context, objType string, objID, uid uint64) error
	UnLike(ctx context.Context, objType string, objID, uid uint64) error
	Collect(ctx context.Context, objType string, objID, uid uint64) error

	UpdateReadCountMany(ctx context.Context, objType string, hmap map[uint64]int64) error

	GetListBatchOfLikes(ctx context.Context, objType string, list *domain.List) ([]domain.Like, error)
}

type likeService struct {
	repo              repository.LikeRepository
	recommendProducer events.RecommendProducer
}

func (l *likeService) UpdateReadCountMany(ctx context.Context, objType string, hmap map[uint64]int64) error {
	return l.repo.UpdateReadCountMany(ctx, objType, hmap)
}

func (l *likeService) IncrReadCountMany(ctx context.Context, objType string, objIDs []uint64) error {
	return l.repo.IncrReadCountMany(ctx, objType, objIDs)
}

func (l *likeService) GetListBatchOfLikes(ctx context.Context, objType string, list *domain.List) ([]domain.Like, error) {
	// TODO 可考虑排除老旧的数据，提高速度
	return l.repo.GetPublishedByBatch(ctx, objType, nil)
}

func (l *likeService) IncrReadCount(ctx context.Context, objType string, objID uint64) error {
	return l.repo.IncrReadCount(ctx, objType, objID)
}

// Like 没控制重复点赞
func (l *likeService) Like(ctx context.Context, objType string, objID, uid uint64) error {
	go func() {
		err := l.recommendProducer.ProduceRecommendEvent(ctx, &events.RecommendEvent{
			FeedbackType: "like",
			UserId:       strconv.FormatUint(uid, 10),
			ItemId:       objType + ":" + strconv.FormatUint(objID, 10),
			Timestamp:    time.Now().Format("2006.01.02 15:04:05"),
		})
		if err != nil {
			// log
		}
	}()

	return l.repo.IncrLikeCount(ctx, objType, objID, uid)
}

func (l *likeService) UnLike(ctx context.Context, objType string, objID, uid uint64) error {
	return l.repo.DecrLikeCount(ctx, objType, objID, uid)
}

func (l *likeService) Collect(ctx context.Context, objType string, objID, uid uint64) error {
	go func() {
		err := l.recommendProducer.ProduceRecommendEvent(ctx, &events.RecommendEvent{
			FeedbackType: "collect",
			UserId:       strconv.FormatUint(uid, 10),
			ItemId:       objType + ":" + strconv.FormatUint(objID, 10),
			Timestamp:    time.Now().Format("2006.01.02 15:04:05"),
		})
		if err != nil {
			// log
		}
	}()

	return l.repo.AddCollectionItem(ctx, objType, objID, uid)
}

func NewLikeService(repo repository.LikeRepository, rc events.RecommendProducer) LikeService {
	return &likeService{
		repo:              repo,
		recommendProducer: rc,
	}
}
