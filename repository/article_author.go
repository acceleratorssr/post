package repository

import (
	"context"
	"post/domain"
	"post/repository/dao"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type articleAuthorRepository struct {
	dao dao.ArticleDao
}

func NewArticleAuthorRepository(dao dao.ArticleDao) ArticleAuthorRepository {
	return &articleAuthorRepository{dao: dao}
}

func (a *articleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, ToEntity(art))
}

func (a *articleAuthorRepository) Update(ctx context.Context, art domain.Article) error {
	return a.dao.UpdateByID(ctx, ToEntity(art))
}
