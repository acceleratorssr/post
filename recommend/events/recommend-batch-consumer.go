package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/zhenghaoz/gorse/client"
	sarama_extra "post/pkg/sarama-extra"
	"time"
)

const topicRecommend = "article_recommend"

type BatchKafkaRecommendConsumer struct {
	groupID string
	client  sarama.Client
	gorse   *client.GorseClient
}

type RecommendEvent struct {
	FeedbackType string `json:"feedback_type"`
	ArticleID    uint64 `json:"article_id"`
	UserId       string `json:"uid"`
	ItemId       string `json:"item_id"`
	Timestamp    string `json:"timestamp"`
}

func NewKafkaRecommendConsumer(client sarama.Client) *BatchKafkaRecommendConsumer {
	return &BatchKafkaRecommendConsumer{
		client: client,
	}
}

func (k *BatchKafkaRecommendConsumer) Start(topic string) error {
	cg, err := sarama.NewConsumerGroupFromClient("recommend", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicRecommend},
			sarama_extra.NewConsumerBatchHandler[RecommendEvent](k.Consume, sarama_extra.WithBatchSize[RecommendEvent](100)))
		if err != nil {
			panic(err)
		}
	}()
	return err
}

func (k *BatchKafkaRecommendConsumer) Consume(msg []*sarama.ConsumerMessage, t []RecommendEvent) error {
	const maxFeedbackBatchSize = 100 // 控制长短，没测

	var baseTimeout = 2 * time.Second
	var extendedTimeout = 5 * time.Second

	ctx := context.Background()
	var fb []client.Feedback

	for _, event := range t {
		fb = append(fb, client.Feedback{
			FeedbackType: event.FeedbackType,
			UserId:       event.UserId,
			ItemId:       event.ItemId,
			Timestamp:    event.Timestamp,
		})
	}

	var timeout time.Duration
	if len(fb) > maxFeedbackBatchSize {
		timeout = extendedTimeout
	} else {
		timeout = baseTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	item, err := k.gorse.InsertFeedback(ctx, fb)
	if err != nil || item.RowAffected == 0 {
		// log
	}

	return err
}
