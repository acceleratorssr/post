package events

import (
	"github.com/IBM/sarama"
)

type LargeMessagesProducer sarama.SyncProducer

type SmallMessagesProducer sarama.SyncProducer

// NewKafkaSyncProducerForLargeMessages
// 使用 ZSTD 压缩算法的生产者
func NewKafkaSyncProducerForLargeMessages(client sarama.Client) LargeMessagesProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionZSTD // sarama.CompressionGZIP

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

// NewKafkaSyncProducerForSmallMessages
// 使用 Snappy 压缩算法的消费者
func NewKafkaSyncProducerForSmallMessages(client sarama.Client) SmallMessagesProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionSnappy // sarama.CompressionLZ4

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func NewKafkaPublishProducer(p LargeMessagesProducer) PublishedProducer {
	return &KafkaPublishedSyncProducer{
		producer: p,
	}
}

func NewKafkaReadProducer(p SmallMessagesProducer) ReadProducer {
	return &KafkaReadSyncProducer{
		producer: p,
	}
}

func NewKafkaRecommendProducer(p SmallMessagesProducer) RecommendProducer {
	return &KafkaSyncRecommendProducer{
		producer: p,
	}
}
