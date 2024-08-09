package repository

import (
	"context"
	"post/internal/domain"
	"post/internal/repository/cache"
	"post/internal/repository/dao"
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
	temp, err := a.toDomain(byID)
	temp[0].Author.Id = byID.Authorid
	//temp.Author.Name, err = a.userRepo.FindByID(ctx,temp.Author.Id)
	return temp[0], err
}

func (a *articleReaderRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.InsertReader(ctx, ToReaderEntity(art))
}

func (a *articleReaderRepository) Sync(ctx context.Context, art domain.Article) error {
	err := a.dao.SyncStatus(ctx, ToAuthorEntity(art))
	// todo 同步删除likeRepo内的数据
	if err != nil { // 防止发布出现错误
		a.cache.DeleteFirstPage(ctx, art.ID)
		a.cache.SetArticleDetail(ctx, art.ID, art)
	}

	return err
}

func (a *articleReaderRepository) toDomain(art ...dao.ArticleReader) ([]domain.Article, error) {
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
