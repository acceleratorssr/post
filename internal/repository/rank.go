package repository

import (
	"context"
	"post/internal/domain"
	"post/internal/repository/cache"
)

type RankRepository interface {
	ReplaceTopNDetail(ctx context.Context, arts []domain.Article) error
	ReplaceTopNBrief(ctx context.Context, arts []domain.Article) error
	GetRankTopNBrief(ctx context.Context) ([]domain.Article, error)
}

type BatchRankCache struct {
	redisCache cache.RankCache
	localCache cache.RankCache
}

func NewBatchRankCache(cache cache.RankCache, localCache *cache.LocalCacheForRank) RankRepository {
	return &BatchRankCache{
		redisCache: cache,
		localCache: localCache,
	}
}

func (b *BatchRankCache) GetRankTopNBrief(ctx context.Context) ([]domain.Article, error) {
	brief, err := b.localCache.GetTopNBrief(ctx)
	if err == nil {
		return brief, nil
	}

	brief, err = b.redisCache.GetTopNBrief(ctx)
	if err == nil {
		go b.localCache.SetTopNBrief(ctx, brief)
	}

	return brief, err
}

func (b *BatchRankCache) ReplaceTopNDetail(ctx context.Context, arts []domain.Article) error {
	go b.localCache.SetTopN(ctx, arts)
	return b.redisCache.SetTopN(ctx, arts)
}

func (b *BatchRankCache) ReplaceTopNBrief(ctx context.Context, arts []domain.Article) error {
	go b.localCache.SetTopNBrief(ctx, arts)
	return b.redisCache.SetTopNBrief(ctx, arts)
}
