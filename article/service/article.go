package service

import (
	"context"
	"post/article/domain"
	"post/article/events"
	"post/article/repository"
	"strconv"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, art *domain.Article) error
	Publish(ctx context.Context, art *domain.Article) error
	Withdraw(ctx context.Context, aid, uid uint64) error

	ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error)

	GetAuthorModelsByID(ctx context.Context, aid, uid uint64) (*domain.Article, error)

	GetPublishedByID(ctx context.Context, id, uid uint64) (*domain.Article, error)
	ListPublished(ctx context.Context, list *domain.List, uid uint64) ([]domain.Article, error)

	GetArtByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error)
}

type articleService struct {
	author            repository.ArticleAuthorRepository
	reader            repository.ArticleReaderRepository
	readProducer      events.ReadProducer
	publishedProducer events.PublishedProducer
	recommendProducer events.RecommendProducer

	ch chan events.ReadEvent
}

func (svc *articleService) GetArtByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error) {
	ds, err := svc.reader.GetPublishedByIDs(ctx, aids)
	if err != nil {
		return nil, err
	}
	go func() {
		now := time.Now().Format("2006.01.02 15:04:05")
		for _, art := range ds {
			err := svc.recommendProducer.ProduceRecommendEvent(ctx, &events.RecommendEvent{
				FeedbackType: "read",
				UserId:       "system",
				ItemId:       "article:" + strconv.FormatUint(art.ID, 10),
				Timestamp:    now,
			})
			if err != nil {
				// log
			}
		}
	}()
	return ds, err
}

func (svc *articleService) ListPublished(ctx context.Context, list *domain.List, uid uint64) ([]domain.Article, error) {
	published, err := svc.reader.ListPublished(ctx, list, uid)
	if err != nil {
		return nil, err
	}

	return published, nil
}

func (svc *articleService) GetPublishedByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	art, err := svc.reader.GetPublishedByID(ctx, aid, uid)
	if err == nil { // 增加阅读数
		go func() {
			svc.readProducer.ProduceReadEvent(ctx, &events.ReadEvent{
				Uid: uid,
				Aid: aid,
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

func (svc *articleService) GetAuthorModelsByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	art, err := svc.author.GetByID(ctx, aid, uid)
	go func() {
		err := svc.recommendProducer.ProduceRecommendEvent(ctx, &events.RecommendEvent{
			FeedbackType: "read",
			UserId:       strconv.FormatUint(uid, 10),
			ItemId:       "article:" + strconv.FormatUint(art.ID, 10),
			Timestamp:    time.Now().Format("2006.01.02 15:04:05"),
		})
		if err != nil {
			// log
		}
	}()
	return art, err
}

// ListSelf 仅缓存符合 LastValue 的前 limit 条数据
func (svc *articleService) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	var authorList []domain.Article
	aList, err := svc.author.ListSelf(ctx, uid, list)
	authorList = aList
	if err != nil {
		return nil, err
	}

	return authorList, nil
}

// Save Author表保存
// 此时reader的对应数据一定是不存在或者是过期状态
// 草稿态
func (svc *articleService) Save(ctx context.Context, art *domain.Article) error {
	if art.ID != 0 {
		return svc.author.Update(ctx, art)
	}
	return svc.author.Create(ctx, art)
}

// Publish 保存并发布
// 草稿态 => 发布态
func (svc *articleService) Publish(ctx context.Context, art *domain.Article) error {
	var err error
	if art.ID != 0 {
		err = svc.author.Update(ctx, art)
	} else {
		// 新建编辑后，立刻发布
		err = svc.author.Create(ctx, art)
	}
	if err != nil {
		return err
	}

	// 发布文章到kafka上，供其他服务异步消费
	go func() {
		e := svc.publishedProducer.ProducePublishedEvent(ctx, &events.PublishEvent{
			Article:   svc.toMQ(art),
			OnlyCache: false,
			Uid:       ctx.Value("uid").(uint64),
			Delete:    false,
		})
		if e != nil {
			// todo 记录mysql任务表，扫表重发
		}
	}()
	return err
}

func (svc *articleService) Withdraw(ctx context.Context, aid, uid uint64) error {
	err := svc.reader.Withdraw(ctx, aid, uid)

	go func() {
		e := svc.publishedProducer.ProducePublishedEvent(ctx, &events.PublishEvent{
			Article: &events.Article{
				ID: aid,
			},
			OnlyCache: false,
			Uid:       uid,
			Delete:    true,
		})
		if e != nil {
		}
	}()

	return err
}

func (svc *articleService) toMQ(art *domain.Article) *events.Article {
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

func NewArticleService(author repository.ArticleAuthorRepository,
	reader repository.ArticleReaderRepository,
	producer events.ReadProducer,
	publishedProducer events.PublishedProducer,
	recommendProducer events.RecommendProducer) ArticleService {
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
		author:            author,
		reader:            reader,
		readProducer:      producer,
		publishedProducer: publishedProducer,
		recommendProducer: recommendProducer,
	}
}
