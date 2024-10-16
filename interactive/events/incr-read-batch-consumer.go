package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/interactive/repository"
	"post/pkg/sarama-extra"
	"time"
)

const topicRead = "article_read"

type BatchKafkaConsumer struct {
	client sarama.Client
	repo   repository.LikeRepository
}

func NewBatchKafkaConsumer(client sarama.Client,
	repo repository.LikeRepository) *BatchKafkaConsumer {
	return &BatchKafkaConsumer{
		client: client,
		repo:   repo,
	}
}
func (b *BatchKafkaConsumer) Start(topic string) error {
	cg, err := sarama.NewConsumerGroupFromClient("t", b.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicRead},
			sarama_extra.NewConsumerBatchHandler[ReadEvent](b.Consume, sarama_extra.WithDuration[ReadEvent](time.Second)))
		if err != nil {
			panic(err)
		}
	}()
	return err
}
func (b *BatchKafkaConsumer) Consume(msg []*sarama.ConsumerMessage, t []ReadEvent) error {
	artsID := make([]uint64, 0, len(t))
	for _, v := range t {
		artsID = append(artsID, v.Aid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// todo article
	err := b.repo.IncrReadCountMany(ctx, "article", artsID)
	if err != nil {
		// log
	}

	return nil
}
