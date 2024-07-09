package service

import (
	"context"
	"post/domain"
	"post/repository"
)

// mockgen -source=D:\桌面\pkg\post\service\article.go -package=svcmocks -destination=D:\桌面\pkg\post\service\mock\article_mock.go

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
}

// todo 换成双repo
type articleService struct {
	author repository.ArticleAuthorRepository
	reader repository.ArticleReaderRepository
}

func NewArticleService(author repository.ArticleAuthorRepository,
	reader repository.ArticleReaderRepository) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
	}
}

// Save Author表保存
// 此时reader的对应数据一定是不存在或者是过期状态
// 草稿态
func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.TypeSaved
	if art.ID != 0 {
		return art.ID, a.author.Update(ctx, art)
	}
	return a.author.Create(ctx, art)
}

// Publish
// 草稿态 => 发布态
func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	var err error
	art.Status = domain.TypePublished
	if art.ID != 0 {
		err = a.author.Update(ctx, art)
	} else {
		// 制作库
		art.ID, err = a.author.Create(ctx, art)
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}
	// 线上库
	var id int64
	for i := 0; i < 5; i++ { //无脑重试
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
	}
	return id, err
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	art.Status = domain.TypeWithdraw
	return a.reader.Sync(ctx, art)
}
