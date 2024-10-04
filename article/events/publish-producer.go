package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type PublishedProducer interface {
	ProducePublishedEvent(ctx context.Context, event *PublishEvent) error
}

type KafkaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaPublishProducer(p sarama.SyncProducer) PublishedProducer {
	return &KafkaSyncProducer{
		producer: p,
	}
}

func (k *KafkaSyncProducer) ProducePublishedEvent(ctx context.Context, event *PublishEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_published",
		Value: sarama.ByteEncoder(data),
	})
	return err
}
