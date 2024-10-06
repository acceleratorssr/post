package repository

import (
	"context"
	"post/article/domain"
	"post/article/repository/cache"
	"post/article/repository/dao"
)

type ArticleReaderRepository interface {
	// Save 包含新建或者更新
	Save(ctx context.Context, art *domain.Article) error
	Withdraw(ctx context.Context, aid, uid uint64) error
	GetPublishedByID(ctx context.Context, id uint64) (*domain.Article, error)
	GetPublishedByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error)
	ListPublished(ctx context.Context, list *domain.List) ([]domain.Article, error)
	ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error)
}

type articleReaderRepository struct {
	dao   dao.ArticleDao
	cache cache.RedisArticleCache
}

func (a *articleReaderRepository) GetPublishedByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error) {
	arts, err := a.dao.GetPublishedByIDs(ctx, aids)
	if err != nil {
		return nil, err
	}

	return a.toDomain(arts...)
}

func (a *articleReaderRepository) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	arts, err := a.dao.ListByID(ctx, uid, a.toDAOlist(list))
	if err != nil {
		return nil, err
	}

	return a.toDomain(arts...)
}

func (a *articleReaderRepository) ListPublished(ctx context.Context, list *domain.List) ([]domain.Article, error) {
	if list.LastValue == 0 && list.Limit <= 100 {
		//return a.cache.ListPublished(ctx, offset, limit, keyword, desc)
	}
	arts, err := a.dao.ListPublished(ctx, a.toDAOlist(list))
	if err != nil {
		return nil, err
	}

	return a.toDomain(arts...)
}

func (a *articleReaderRepository) GetPublishedByID(ctx context.Context, id uint64) (*domain.Article, error) {
	// 注意，如果帖子数据在OSS上时，需要从前端直接获取，因为考虑到内容不算太敏感
	art, err := a.dao.GetPublishedByID(ctx, id)
	if err != nil {
		return &domain.Article{}, err
	}
	temp, err := a.toDomain(*art)
	temp[0].Author.Id = art.Authorid

	return &temp[0], err
}

func (a *articleReaderRepository) Save(ctx context.Context, art *domain.Article) error {
	return a.dao.UpsertReader(ctx, &a.ToReaderEntity([]domain.Article{*art}...)[0])
}

func (a *articleReaderRepository) Withdraw(ctx context.Context, aid, uid uint64) error {
	err := a.dao.DeleteReader(ctx, aid, uid)

	return err
}

func (a *articleReaderRepository) toDomain(art ...dao.ArticleReader) ([]domain.Article, error) {
	domainA := make([]domain.Article, 0, len(art))
	for i, _ := range art {
		domainA = append(domainA, domain.Article{
			ID:      uint64(art[i].SnowID),
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

func (a *articleReaderRepository) toDAOlist(list *domain.List) *dao.List {
	return &dao.List{
		Desc:      list.Desc,
		LastValue: list.LastValue,
		Limit:     list.Limit,
		OrderBy:   list.OrderBy,
	}
}

func (a *articleReaderRepository) ToReaderEntity(art ...domain.Article) []dao.ArticleReader {
	var ans = make([]dao.ArticleReader, 0, len(art))
	for i, _ := range art {
		ans[i] = dao.ArticleReader{
			SnowID:   int64(art[i].ID),
			Title:    art[i].Title,
			Content:  art[i].Content,
			Authorid: art[i].Author.Id,
		}
	}
	return ans
}

func NewArticleReaderRepository(dao dao.ArticleDao) ArticleReaderRepository {
	return &articleReaderRepository{
		dao: dao}
}
