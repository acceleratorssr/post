package events

import "post/article/bridge"

func NewKafkaPublishProducer(p bridge.LargeMessagesProducer) PublishedProducer {
	return &KafkaSyncProducer{
		producer: p,
	}
}

func NewKafkaReadProducer(p bridge.SmallMessagesProducer) ReadProducer {
	return &KafkaReadSyncProducer{
		producer: p,
	}
}
