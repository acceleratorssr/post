package repository

import (
	"context"
	"post/search/domain"
	"post/search/repository/dao"
)

type articleRepository struct {
	dao  dao.ArticleDAO
	tags dao.TagDAO
}

func (a *articleRepository) DeleteArticle(ctx context.Context, id uint64) error {
	return a.dao.DeleteArticle(ctx, id)
}

func (a *articleRepository) SearchArticle(ctx context.Context,
	keywords []string) ([]domain.Article, error) {
	ids, err := a.tags.Search(ctx, "article", keywords)
	if err != nil {
		return nil, err
	}
	arts, err := a.dao.Search(ctx, ids, keywords)
	if err != nil {
		return nil, err
	}

	mappedArticles := make([]domain.Article, len(arts)) // 创建目标切片
	for i, src := range arts {
		mappedArticles[i] = domain.Article{
			ID:      src.Id,
			Title:   src.Title,
			Content: src.Content,
			Tags:    src.Tags,
		}
	}

	return mappedArticles, nil
}

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.ID,
		Title:   msg.Title,
		Content: msg.Content,
	})
}

func NewArticleRepository(d dao.ArticleDAO, td dao.TagDAO) ArticleRepository {
	return &articleRepository{
		dao:  d,
		tags: td,
	}
}
