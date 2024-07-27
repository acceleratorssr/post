package repository

import (
	"context"
	"post/domain"
	"post/repository/cache"
)

type RankRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

type BatchRankCache struct {
	cache cache.ArticleCache
}

func (b *BatchRankCache) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return b.cache.SetTopN(ctx, arts)
}
