package repository

import (
	"context"
	"post/domain"
	"post/repository/dao"
)

type ArticleReaderRepository interface {
	// Save 包含新建或者更新
	Save(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) error
}

type articleReaderRepository struct {
	dao dao.ArticleDao
}

func NewArticleReaderRepository(dao dao.ArticleDao) ArticleReaderRepository {
	return &articleReaderRepository{dao: dao}
}

func (a *articleReaderRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, ToEntity(art))
}

func (a *articleReaderRepository) Sync(ctx context.Context, art domain.Article) error {
	return a.dao.SyncStatus(ctx, ToEntity(art))
}
