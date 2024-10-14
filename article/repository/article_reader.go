package repository

import (
	"context"
	"math"
	"post/article/domain"
	"post/article/events"
	"post/article/repository/cache"
	"post/article/repository/dao"
	"strconv"
)

type ArticleReaderRepository interface {
	// Save 包含新建或者更新
	Save(ctx context.Context, art *domain.Article) error
	Withdraw(ctx context.Context, aid, uid uint64) error
	GetPublishedByID(ctx context.Context, aid, uid uint64) (*domain.Article, error)
	GetPublishedByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error)
	ListPublished(ctx context.Context, list *domain.List, uid uint64) ([]domain.Article, error)
	ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error)
}

// todo 加入并检测singleflight的优化
type articleReaderRepository struct {
	dao               dao.ArticleDao
	cache             cache.ArticleCache
	publishedProducer events.PublishedProducer
}

func (a *articleReaderRepository) GetPublishedByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error) {
	arts, err := a.dao.GetPublishedByIDs(ctx, aids)
	if err != nil {
		return nil, err
	}

	return a.toDomain(arts...), nil
}

func (a *articleReaderRepository) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	arts, err := a.dao.ListByID(ctx, uid, a.toDAOlist(list))
	if err != nil {
		return nil, err
	}

	return a.toDomain(arts...), nil
}

func (a *articleReaderRepository) ListPublished(ctx context.Context, list *domain.List, uid uint64) ([]domain.Article, error) {
	arts, err := a.dao.ListPublished(ctx, a.toDAOlist(list))
	if err != nil {
		return nil, err
	}

	res := a.toDomain(arts...)

	if list.LastValue == math.MaxInt && list.Limit <= 100 {
		go func() {
			for i := 0; i < len(res); i++ {
				e := a.publishedProducer.ProducePublishedEvent(ctx, &events.PublishEvent{
					Article:   a.toMQ(res[i]),
					OnlyCache: true,
					Uid:       uid,
				})
				if e != nil {
					// todo 记录mysql任务表，扫表重发
				}
			}
		}()
	}

	return res, nil
}

func (a *articleReaderRepository) GetPublishedByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	// 注意，如果帖子数据在OSS上时，需要从前端直接获取，因为考虑到内容不算太敏感
	ans, err := a.cache.GetListDetailByHashKey(ctx, aid, strconv.FormatUint(uid, 10))
	if err != nil {
		ans, err = a.cache.GetArticleDetail(ctx, aid)
	}
	if err == nil {
		return ans, nil
	}

	art, err := a.dao.GetPublishedByID(ctx, aid)
	if err != nil {
		return &domain.Article{}, err
	}
	temp := a.toDomain(*art)
	temp[0].Author.Id = art.Authorid

	return &temp[0], nil
}

func (a *articleReaderRepository) Save(ctx context.Context, art *domain.Article) error {
	return a.dao.UpsertReader(ctx, &a.ToReaderEntity([]domain.Article{*art}...)[0])
}

func (a *articleReaderRepository) Withdraw(ctx context.Context, aid, uid uint64) error {
	err := a.dao.DeleteReader(ctx, aid, uid)
	if err != nil {
		// log
		return err
	}

	err = a.cache.DeleteListInfo(ctx, aid)
	if err != nil {
		// log
	}
	return err
}

func (a *articleReaderRepository) toDomain(art ...dao.ArticleReader) []domain.Article {
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
	return domainA
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

func (a *articleReaderRepository) toMQ(art domain.Article) *events.Article {
	return &events.Article{
		ID:      art.ID,
		Title:   art.Title,
		Content: art.Content,
		Author: events.Author{
			Id:   art.Author.Id,
			Name: art.Author.Name,
		},
		Utime: art.Utime,
		Ctime: art.Ctime,
	}
}

func NewArticleReaderRepository(dao dao.ArticleDao,
	cache cache.ArticleCache,
	publishedProducer events.PublishedProducer) ArticleReaderRepository {
	return &articleReaderRepository{
		dao:               dao,
		cache:             cache,
		publishedProducer: publishedProducer,
	}
}
