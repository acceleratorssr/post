package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type SaramaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewSaramaProducer(producer sarama.SyncProducer, topic string) *SaramaProducer {
	return &SaramaProducer{
		producer: producer,
		topic:    topic,
	}
}

func (s *SaramaProducer) InconsistentEventProducer(ctx context.Context, event *InconsistentEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = s.producer.SendMessage(msg)
	return err
}
