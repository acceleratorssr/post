package events

import "post/pkg/sarama_ex"

func NewKafkaConsumer(consumer *KafkaPublishedConsumer) []sarama_ex.Consumer {
	return []sarama_ex.Consumer{
		consumer,
	}
}
