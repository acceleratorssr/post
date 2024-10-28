package repository

import (
	"context"
	"post/search/domain"
	"post/search/repository/dao"
)

type articleRepository struct {
	dao dao.ArticleDAO
	//tags dao.TagDAO
}

func (a *articleRepository) DeleteArticle(ctx context.Context, id uint64) error {
	return a.dao.DeleteArticle(ctx, id)
}

func (a *articleRepository) SearchArticle(ctx context.Context, keywords []string, vector []float32, limit int) ([]domain.Article, error) {
	//ids, err := a.tags.Search(ctx, "article", keywords)
	//if err != nil {
	//	return nil, err
	//}
	arts, err := a.dao.Search(ctx, nil, keywords, vector, limit)
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

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article, vector []float32) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.ID,
		Title:   msg.Title,
		Content: msg.Content,
		Author: dao.Author{
			ID:   msg.Author.Id,
			Name: msg.Author.Name,
		},
	}, vector)
}

func NewArticleRepository(d dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: d,
	}
}
