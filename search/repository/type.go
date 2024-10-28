package repository

import (
	"context"
	"post/search/domain"
)

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article, vector []float32) error
	SearchArticle(ctx context.Context, keywords []string, vector []float32, limit int) ([]domain.Article, error)
	DeleteArticle(ctx context.Context, id uint64) error
}
