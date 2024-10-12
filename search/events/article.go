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

type ArticleEvent struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Status  int32  `json:"status"`
	Content string `json:"content"`
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
	return a.syncSvc.InputArticle(ctx, a.toDomain(evt))
}

func (a *ArticleConsumer) toDomain(article ArticleEvent) domain.Article {
	return domain.Article{
		ID:      article.Id,
		Title:   article.Title,
		Content: article.Content,
	}
}
