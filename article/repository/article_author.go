package repository

import (
	"context"
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
	node  dao.UniqueID
}

func NewArticleAuthorRepository(dao dao.ArticleDao, cache cache.ArticleCache, node dao.UniqueID) ArticleAuthorRepository {
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

// ListSelf
// 顺序性：两次刷新缓存时，如果中间用户保存数据，且第一次缓存先于第二次缓存构建完成，则会出现数据不一致的情况
// 故需要通过 kafka ，而不是 goroutine 构建缓存保证顺序性；
func (a *articleAuthorRepository) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	res, err := a.dao.GetListByAuthor(ctx, uid, a.toDAOlist(list))
	if err != nil {
		return nil, err
	}
	data, err := a.toDomain(res...)

	// todo kafka producer 写入缓存

	return data, err
}

func (a *articleAuthorRepository) Create(ctx context.Context, art *domain.Article) error {
	art.ID = a.node.Generate()

	return a.dao.Insert(ctx, a.ToAuthorEntity(art))
}

func (a *articleAuthorRepository) Update(ctx context.Context, art *domain.Article) error {
	return a.dao.UpdateByID(ctx, a.ToAuthorEntity(art))
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

func (a *articleAuthorRepository) ToAuthorEntity(art *domain.Article) *dao.ArticleAuthor {
	return &dao.ArticleAuthor{
		SnowID:   art.ID,
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
	}
}
