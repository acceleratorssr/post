package service

import (
	"context"
	"math"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/article/domain"
	"post/article/repository"
	"post/article/utils"
	"time"
)

type RankService interface {
	SetRankTopN(ctx context.Context, n int) error
	GetRankTopNBrief(ctx context.Context) ([]domain.Article, error)
}

type BatchRankService struct {
	artSvc ArticleService
	//likeSvc  service.LikeService
	likeSvc  intrv1.LikeServiceClient
	rankRepo repository.RankRepository
}

type scoreNum struct {
	art   domain.Article
	score float64
}

func NewBatchRankService(artSvc ArticleService, likeSvc intrv1.LikeServiceClient, rankRepo repository.RankRepository) RankService {
	return &BatchRankService{
		artSvc:   artSvc,
		likeSvc:  likeSvc,
		rankRepo: rankRepo,
	}
}

// GetRankTopNBrief todo test
func (svc *BatchRankService) GetRankTopNBrief(ctx context.Context) ([]domain.Article, error) {
	return svc.rankRepo.GetRankTopNBrief(ctx)
}

// SetRankTopN todo test
func (svc *BatchRankService) SetRankTopN(ctx context.Context, n int) error {
	offset := 0
	// 直接从like取出数据，然后找出topn后，返回ids即可，再根据id获取文章
	pq := utils.NewMinHeap()
	now := time.Now().UnixMilli()
	for {
		//likes, err := svc.likeSvc.GetListBatchOfLikes(ctx, "article", offset, n, now)
		likes, err := svc.likeSvc.GetListBatchOfLikes(ctx, &intrv1.GetListBatchOfLikesRequest{
			ObjType: "article",
			Offset:  int32(offset),
			Limit:   int32(n),
			Now:     now,
		})
		if err != nil {
			return err
		}

		// todo likes内有空吗
		ids, score := svc.CountRank(svc.grpc2domain(likes.Data))
		for i := 0; i < len(ids); i++ {
			// 此处是同一批次的，即当前的时间戳相同，所以一定存在该时间下的No100
			if pq.GetLen() > n && pq.GetMin().Score > score[i] {
				continue
			}
			pq.Insert(ids[i], score[i])
		}

		// 无剩余数据
		if len(likes.Data) < n {
			break
		}
		offset = offset + len(likes.Data)
	}

	res := make([]int64, pq.GetLen())
	for i := pq.GetLen() - 1; i >= 0; i-- {
		v := pq.ExtractMin()
		res[i] = v.ID
	}

	arts, err := svc.artSvc.GetPublishedByIDS(ctx, res)
	if err != nil {
		panic(err)
	}

	// 缓存每篇热帖
	err = svc.rankRepo.ReplaceTopNDetail(ctx, arts)
	if err != nil {
		panic(err)
	}

	// todo 有点粗糙缓存帖子的简要信息在一起
	err = svc.rankRepo.ReplaceTopNBrief(ctx, arts)
	if err != nil {
		panic(err)
	}

	return err
}

// CountRank 计算排行榜公式（简化为只考虑时间和点赞）：Score = (P - 1) / ((T + 2) ^ G)
func (svc *BatchRankService) CountRank(ds []domain.Like) ([]int64, []float64) {
	now := time.Now().UnixMilli()
	ids := make([]int64, 0, len(ds))
	scores := make([]float64, 0, len(ds))

	for _, d := range ds {
		ids = append(ids, d.ID)
		scores = append(scores, float64(d.LikeCount-1)/(math.Pow(float64(now-d.Ctime+2), 1.5)))
	}
	return ids, scores
}
func (svc *BatchRankService) grpc2domain(like []*intrv1.Like) []domain.Like {
	ds := make([]domain.Like, 0, len(like))
	for _, l := range like {
		ds = append(ds, domain.Like{
			ID:        l.ID,
			Ctime:     l.Ctime,
			LikeCount: l.LikeCount,
		})
	}
	return ds
}
