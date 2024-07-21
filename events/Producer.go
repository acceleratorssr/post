package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(ctx context.Context, event *ReadEvent, topic string) error
}

type KafkaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(p sarama.SyncProducer) Producer {
	return &KafkaSyncProducer{
		producer: p,
	}
}

// ProduceReadEvent 当重试逻辑变复杂时，使用装饰器模式
func (k *KafkaSyncProducer) ProduceReadEvent(ctx context.Context, event *ReadEvent, topic string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	ID  int64
	Uid int64
	Aid int64
}
