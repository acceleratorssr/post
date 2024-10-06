package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/article/domain"
	"post/article/repository/cache"
	"post/article/repository/dao"
	"post/pkg/sarama_ex"
	"strconv"
	"time"
)

type KafkaPublishedConsumer struct {
	groupID string
	client  sarama.Client
	dao     dao.ArticleDao
	cache   cache.ArticleCache
}

func NewKafkaPublishedConsumer(client sarama.Client,
	dao dao.ArticleDao, cache cache.ArticleCache) *KafkaPublishedConsumer {
	return &KafkaPublishedConsumer{
		client: client,
		dao:    dao,
		cache:  cache,
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
	var err error

	if t.OnlyCache {
		// 仅构建 用户 专属的批量文章缓存，hash[uid]: aid-article
		k.cache.SetListDetailByHashKey(ctx, []domain.Article{*k.toDomain(t)}, strconv.FormatUint(t.Uid, 10))
	} else {
		// 持久化发布的文章
		err = k.dao.UpsertReader(ctx, k.toDAO(t))
		if err != nil {
			// log
			return err
		}
		// 缓存发布文章，string
		err = k.cache.SetArticleDetail(ctx, k.toDomain(t))
		if err != nil {
			// log
			return err
		}
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

func (k *KafkaPublishedConsumer) toDAO(msg PublishEvent) *dao.ArticleReader {
	return &dao.ArticleReader{
		ID:       msg.Article.ID,
		Title:    msg.Article.Title,
		Content:  msg.Article.Content,
		Authorid: msg.Article.Author.Id,
		Ctime:    msg.Article.Ctime,
		Utime:    msg.Article.Utime,
	}
}
