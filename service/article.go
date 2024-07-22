package service

import (
	"context"
	"post/domain"
	"post/events"
	"post/repository"
)

// mockgen -source=D:\桌面\pkg\post\service\incr_read_producer.go -package=svcmocks -destination=D:\桌面\pkg\post\service\mock\article_mock.go

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error)
	GetAuthorModelsByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedByID(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	author   repository.ArticleAuthorRepository
	reader   repository.ArticleReaderRepository
	producer events.Producer

	ch chan events.ReadEvent
}

func NewArticleService(author repository.ArticleAuthorRepository,
	reader repository.ArticleReaderRepository,
	producer events.Producer) ArticleService {
	// producer也可以批处理，要多写一个consumer
	//ch := make(chan events.ReadEvent, 100)
	//go func() {
	//	for {
	//		uids := make([]int64, 0, 100)
	//		aids := make([]int64, 0, 100)
	//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//		for i := 0; i < 100; i++ {
	//			select {
	//			case e, ok := <-ch:
	//				if !ok {
	//					break
	//				}
	//				uids = append(uids, e.Uid)
	//				aids = append(aids, e.Aid)
	//			case <-ctx.Done():
	//				break
	//			}
	//		}
	//		cancel()
	//		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//		producer.ProduceReadEventMany(ctx, &events.ReadEventMany{
	//			Aid: aids,
	//			Uid: uids,
	//		})
	//		cancel()
	//	}
	//}()
	return &articleService{
		author:   author,
		reader:   reader,
		producer: producer,
		//ch:       ch,
	}
}

func (svc *articleService) GetPublishedByID(ctx context.Context, id, uid int64) (domain.Article, error) {
	art, err := svc.reader.GetPublishedByID(ctx, id)
	if err == nil {
		go func() {
			svc.producer.ProduceReadEvent(ctx, &events.ReadEvent{
				Uid: uid,
				Aid: id,
			})
		}()

		//// 批量
		//go func() {
		//	svc.ch <- events.ReadEvent{
		//		Uid: uid,
		//		Aid: id,
		//	}
		//}()
	}
	return art, err
}

func (svc *articleService) GetAuthorModelsByID(ctx context.Context, id int64) (domain.Article, error) {
	return svc.author.GetByID(ctx, id)
}

func (svc *articleService) List(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error) {
	return svc.author.List(ctx, uid, limit, offset)
}

// Save Author表保存
// 此时reader的对应数据一定是不存在或者是过期状态
// 草稿态
func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.TypeSaved
	if art.ID != 0 {
		return art.ID, svc.author.Update(ctx, art)
	}
	return svc.author.Create(ctx, art)
}

// Publish
// 草稿态 => 发布态
func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	var err error
	art.Status = domain.TypePublished
	if art.ID != 0 {
		err = svc.author.Update(ctx, art)
	} else {
		// 制作库
		art.ID, err = svc.author.Create(ctx, art)
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
		id, err = svc.reader.Save(ctx, art)
		if err == nil {
			break
		}
	}
	return id, err
}

func (svc *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	art.Status = domain.TypeWithdraw
	return svc.reader.Sync(ctx, art)
}
