package repository

import (
	"context"
	"post/domain"
	"post/repository/cache"
	"post/repository/dao"
	"time"
)

type ArticleReaderRepository interface {
	// Save 包含新建或者更新
	Save(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) error
	GetPublishedByID(ctx context.Context, id int64) (domain.Article, error)
}

type articleReaderRepository struct {
	dao   dao.ArticleDao
	cache cache.RedisArticleCache
	// userRepo
}

func NewArticleReaderRepository(dao dao.ArticleDao) ArticleReaderRepository {
	return &articleReaderRepository{dao: dao}
}

func (a *articleReaderRepository) GetPublishedByID(ctx context.Context, id int64) (domain.Article, error) {
	// 注意，如果帖子数据在OSS上时，需要从前端直接获取，因为考虑到内容不算太敏感
	byID, err := a.dao.GetPublishedByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	temp := a.toDomain(byID)
	temp.Author.Id = byID.Authorid
	//temp.Author.Name, err = a.userRepo.FindByID(ctx,temp.Author.Id)
	return temp, err
}

func (a *articleReaderRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, ToEntity(art))
}

func (a *articleReaderRepository) Sync(ctx context.Context, art domain.Article) error {
	err := a.dao.SyncStatus(ctx, ToEntity(art))
	if err != nil { // 防止发布出现错误
		a.cache.DeleteFirstPage(ctx, art.ID)
		a.cache.SetArticleDetail(ctx, art.ID, art)
	}

	return err
}

func (a *articleReaderRepository) toDomain(art dao.ArticleReader) domain.Article {
	return domain.Article{
		ID:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.StatusType(art.Status),
		Ctime:   time.UnixMilli(art.Ctime),
		Utime:   time.UnixMilli(art.Utime),
	}
}
