package repository

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"post/article/domain"
	"post/article/repository/cache"
	"post/article/repository/dao"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art *domain.Article) error
	Update(ctx context.Context, art *domain.Article) error
	ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error)
	GetByID(ctx context.Context, aid, uid uint64) (*domain.Article, error)
}

type articleAuthorRepository struct {
	dao   dao.ArticleDao
	cache cache.ArticleCache
	node  *snowflake.Node
}

func NewArticleAuthorRepository(dao dao.ArticleDao, cache cache.ArticleCache, node *snowflake.Node) ArticleAuthorRepository {
	return &articleAuthorRepository{
		dao:   dao,
		cache: cache,
		node:  node,
	}
}

func (a *articleAuthorRepository) GetByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	res, err := a.dao.GetAuthorByID(ctx, aid, uid)
	if err != nil {
		return &domain.Article{}, err
	}

	art, err := a.toDomain(*res)
	if err != nil {
		return &domain.Article{}, err
	}
	return &art[0], err
}
func (a *articleAuthorRepository) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	if list.LastValue == 0 && list.Limit <= 100 {
		data, err := a.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return data, nil
		}
	}

	res, err := a.dao.GetListByAuthor(ctx, uid, a.toDAOlist(list))
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

func (a *articleAuthorRepository) Create(ctx context.Context, art *domain.Article) error {
	art.ID = uint64(a.node.Generate().Int64())
	defer a.cache.DeleteFirstPage(ctx, art.Author.Id)
	return a.dao.Insert(ctx, ToAuthorEntity(art))
}

func (a *articleAuthorRepository) Update(ctx context.Context, art *domain.Article) error {
	defer a.cache.DeleteFirstPage(ctx, art.Author.Id)
	return a.dao.UpdateByID(ctx, ToAuthorEntity(art))
}

func (a *articleAuthorRepository) toDAOlist(list *domain.List) *dao.List {
	return &dao.List{
		Desc:      list.Desc,
		LastValue: list.LastValue,
		Limit:     list.Limit,
		OrderBy:   list.OrderBy,
	}
}

func (a *articleAuthorRepository) toDomain(art ...dao.ArticleAuthor) ([]domain.Article, error) {
	domainA := make([]domain.Article, 0, len(art))
	for i, _ := range art {
		domainA = append(domainA, domain.Article{
			ID:      art[i].SnowID,
			Title:   art[i].Title,
			Content: art[i].Content,
			Ctime:   art[i].Ctime,
			Utime:   art[i].Utime,
			Author: domain.Author{
				Id: art[i].Authorid,
			},
		})
	}
	return domainA, nil
}
