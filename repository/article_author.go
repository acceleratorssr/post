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
	cache cache.RedisArticleCache
}

func NewArticleAuthorRepository(dao dao.ArticleDao) ArticleAuthorRepository {
	return &articleAuthorRepository{dao: dao}
}

func (a *articleAuthorRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return a.toDomain(res), nil
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
	data, err := a.toDomainMany(res)

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

// TODO 有没有更方便的转换方式
func (a *articleAuthorRepository) toDomainMany(art []dao.ArticleAuthor) ([]domain.Article, error) {
	domainA := make([]domain.Article, 0)
	for i, _ := range art {
		domainA = append(domainA, a.toDomain(art[i]))

	}
	return domainA, nil
}

func (a *articleAuthorRepository) toDomain(art dao.ArticleAuthor) domain.Article {
	return domain.Article{
		ID:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.StatusType(art.Status),
		Ctime:   time.UnixMilli(art.Ctime),
		Utime:   time.UnixMilli(art.Utime),
	}
}
