package repository

import (
	"context"
	"post/search/domain"
)

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article) error
	SearchArticle(ctx context.Context, keywords []string) ([]domain.Article, error)
}
