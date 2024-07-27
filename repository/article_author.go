package repository

import (
	"context"
	"post/domain"
	"post/repository/cache"
	"post/repository/dao"
	"time"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
}

type articleAuthorRepository struct {
	dao   dao.ArticleDao
	cache cache.ArticleCache
}

func NewArticleAuthorRepository(dao dao.ArticleDao, cache cache.ArticleCache) ArticleAuthorRepository {
	return &articleAuthorRepository{
		dao:   dao,
		cache: cache,
	}
}

func (a *articleAuthorRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	art, err := a.toDomain(res)
	if err != nil {
		return domain.Article{}, err
	}
	return art[0], err
}
func (a *articleAuthorRepository) List(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := a.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return data, nil
		}
	}

	res, err := a.dao.GetListByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	data, err := a.toDomain(res...)

	// 回源
	go func() {
		err = a.cache.SetFirstPage(ctx, uid, data)
		if err != nil {
			//log
		}
	}()

	return data, err
}

func (a *articleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer a.cache.DeleteFirstPage(ctx, art.Author.Id)
	return a.dao.Insert(ctx, ToEntity(art))
}

func (a *articleAuthorRepository) Update(ctx context.Context, art domain.Article) error {
	defer a.cache.DeleteFirstPage(ctx, art.Author.Id)
	return a.dao.UpdateByID(ctx, ToEntity(art))
}

func (a *articleAuthorRepository) toDomain(art ...dao.ArticleAuthor) ([]domain.Article, error) {
	domainA := make([]domain.Article, 0, len(art))
	for i, _ := range art {
		domainA = append(domainA, domain.Article{
			ID:      art[i].Id,
			Title:   art[i].Title,
			Content: art[i].Content,
			Status:  domain.StatusType(art[i].Status),
			Ctime:   time.UnixMilli(art[i].Ctime),
			Utime:   time.UnixMilli(art[i].Utime),
		})
	}
	return domainA, nil
}
