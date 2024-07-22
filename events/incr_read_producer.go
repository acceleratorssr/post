package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(ctx context.Context, event *ReadEvent) error
	ProduceReadEventMany(ctx context.Context, event *ReadEventMany) error
}

type KafkaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(p sarama.SyncProducer) Producer {
	return &KafkaSyncProducer{
		producer: p,
	}
}

func (k *KafkaSyncProducer) ProduceReadEventMany(ctx context.Context, event *ReadEventMany) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_read",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

// ProduceReadEvent 当重试逻辑变复杂时，使用装饰器模式
func (k *KafkaSyncProducer) ProduceReadEvent(ctx context.Context, event *ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_read",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	ID  int64
	Uid int64
	Aid int64
}

type ReadEventMany struct {
	ID  int64
	Uid []int64
	Aid []int64
}
