package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/interactive/repository"
	"post/pkg/sarama-extra"
	"time"
)

//// Consumer 有其他需求时，记得实现该接口即可，wire会填充进去的
//type Consumer interface {
//	Start(topic string) error
//}

type KafkaReadConsumer struct {
	groupID string
	client  sarama.Client
	repo    repository.LikeRepository
}

func NewKafkaIncrReadConsumer(client sarama.Client,
	repo repository.LikeRepository) *KafkaReadConsumer { // todo 为什么不返回接口
	return &KafkaReadConsumer{
		client: client,
		repo:   repo,
	}
}
func (k *KafkaReadConsumer) Start(topic string) error {
	// todo k.groupID
	cg, err := sarama.NewConsumerGroupFromClient("t", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			//[]string{topic},
			[]string{"article_read"},
			sarama_extra.NewHandler[ReadEvent](k.Consume))
		if err != nil {
			panic(err)
		}
	}()
	return err
}

// Consume TODO 非幂等，如果要保证幂等，初步设想需要多一张表记录用户的访问记录，producer前校验该用户是否已经访问过当前页面
func (k *KafkaReadConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// TODO 这里的objType从哪拿
	return k.repo.IncrReadCount(ctx, "article", t.Aid)
}
