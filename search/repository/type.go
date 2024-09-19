package repository

import (
	"context"
	"post/search/domain"
)

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}
