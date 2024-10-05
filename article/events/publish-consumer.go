package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/article/domain"
	"post/article/repository"
	"post/pkg/sarama_ex"
	"time"
)

type KafkaPublishedConsumer struct {
	groupID string
	client  sarama.Client
	repo    repository.ArticleReaderRepository
}

func NewKafkaPublishedConsumer(client sarama.Client,
	repo repository.ArticleReaderRepository) *KafkaPublishedConsumer {
	return &KafkaPublishedConsumer{
		client: client,
		repo:   repo,
	}
}

func (k *KafkaPublishedConsumer) Start(topic string) error {
	cg, err := sarama.NewConsumerGroupFromClient("t", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"article_published"},
			sarama_ex.NewHandler[PublishEvent](k.Consume))
		if err != nil {
			panic(err)
		}
	}()
	return err
}

func (k *KafkaPublishedConsumer) Consume(msg *sarama.ConsumerMessage, t PublishEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := k.repo.Save(ctx, k.toDomain(t))
	if err != nil {
		return err
	} else {
		// 异步写入缓存
	}
	return err
}

func (k *KafkaPublishedConsumer) toDomain(msg PublishEvent) *domain.Article {
	return &domain.Article{
		ID:      msg.Article.ID,
		Title:   msg.Article.Title,
		Content: msg.Article.Content,
		Author: domain.Author{
			Id:   msg.Article.Author.Id,
			Name: msg.Article.Author.Name,
		},
		Ctime: msg.Article.Ctime,
		Utime: msg.Article.Utime,
	}
}
