package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type RecommendProducer interface {
	ProduceRecommendEvent(ctx context.Context, event *RecommendEvent) error
}

type KafkaSyncRecommendProducer struct {
	producer sarama.SyncProducer
}

func (k *KafkaSyncRecommendProducer) ProduceRecommendEvent(ctx context.Context, event *RecommendEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_recommend",
		Value: sarama.ByteEncoder(data),
	})
	return err
}
