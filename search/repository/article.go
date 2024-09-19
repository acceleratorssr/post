package repository

import (
	"context"
	"post/search/domain"
	"post/search/repository/dao"
)

type articleRepository struct {
	dao dao.ArticleDAO
}

func (a *articleRepository) SearchArticle(ctx context.Context,
	uid int64,
	keywords []string) ([]domain.Article, error) {
	arts, err := a.dao.Search(ctx, nil, keywords)
	if err != nil {
		return nil, err
	}

	mappedArticles := make([]domain.Article, len(arts)) // 创建目标切片
	for i, src := range arts {
		mappedArticles[i] = domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
			Tags:    src.Tags,
		}
	}

	return mappedArticles, nil
}

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  msg.Status,
		Content: msg.Content,
	})
}

func NewArticleRepository(d dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: d,
	}
}
