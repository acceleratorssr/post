package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"strconv"
)

type PublishedProducer interface {
	ProducePublishedEvent(ctx context.Context, event *PublishEvent) error
}

type KafkaSyncProducer struct {
	producer sarama.SyncProducer
}

func (k *KafkaSyncProducer) ProducePublishedEvent(ctx context.Context, event *PublishEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_published",
		Key:   sarama.StringEncoder(strconv.FormatUint(event.Article.ID, 10)),
		Value: sarama.ByteEncoder(data),
	})
	return err
}
