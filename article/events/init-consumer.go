package events

import "post/pkg/sarama-extra"

func NewKafkaConsumer(consumer *KafkaPublishedConsumer) []sarama_extra.Consumer {
	return []sarama_extra.Consumer{
		consumer,
	}
}
