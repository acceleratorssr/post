package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/pkg/logger"
	"post/pkg/sarama-extra"
	"post/search/domain"
	"post/search/service"
	"time"
)

const topicSyncArticle = "article_published"

type ArticleConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.Logger
}

func NewArticleConsumer(client sarama.Client,
	l logger.Logger,
	svc service.SyncService) *ArticleConsumer {
	return &ArticleConsumer{
		syncSvc: svc,
		client:  client,
		l:       l,
	}
}

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  Author `json:"author"`
	Ctime   int64  `json:"ctime"`
	Utime   int64  `json:"utime"`
}

type ArticleEvent struct {
	Article   *Article `json:"article"`
	OnlyCache bool     `json:"only_cache"`
	Uid       uint64   `json:"uid"`
	Delete    bool     `json:"delete"`
}

func (a *ArticleConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("es",
		a.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncArticle},
			sarama_extra.NewHandler[ArticleEvent](a.Consume))
		if err != nil {
			panic(err)
		}
	}()
	return err
}

func (a *ArticleConsumer) Consume(sg *sarama.ConsumerMessage,
	evt ArticleEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if evt.Delete {
		return a.syncSvc.DeleteArticle(ctx, evt.Article.ID)
	}
	return a.syncSvc.InputArticle(ctx, a.toDomain(evt))
}

func (a *ArticleConsumer) toDomain(article ArticleEvent) domain.Article {
	return domain.Article{
		ID:      article.Article.ID,
		Title:   article.Article.Title,
		Content: article.Article.Content,
	}
}
