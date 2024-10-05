package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type ReadProducer interface {
	ProduceReadEvent(ctx context.Context, event *ReadEvent) error
	ProduceReadEventMany(ctx context.Context, event *ReadEventMany) error
}

type KafkaReadSyncProducer struct {
	producer sarama.SyncProducer
}

func (k *KafkaReadSyncProducer) ProduceReadEventMany(ctx context.Context, event *ReadEventMany) error {
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

// ProduceReadEvent 当重试逻辑变复杂时，考虑使用装饰器模式
func (k *KafkaReadSyncProducer) ProduceReadEvent(ctx context.Context, event *ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 同步无err，则发送成功
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_read",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	ID  uint64
	Uid uint64
	Aid uint64
}

type ReadEventMany struct {
	ID  uint64
	Uid []uint64
	Aid []uint64
}
