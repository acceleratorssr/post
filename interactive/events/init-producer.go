package events

import "github.com/IBM/sarama"

type SmallMessagesProducer sarama.SyncProducer

func NewKafkaSyncProducerForSmallMessages(client sarama.Client) SmallMessagesProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionSnappy // sarama.CompressionLZ4

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func NewKafkaRecommendProducer(p SmallMessagesProducer) RecommendProducer {
	return &KafkaSyncRecommendProducer{
		producer: p,
	}
}
