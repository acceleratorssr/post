package ioc

import (
	_ "embed"
	"github.com/IBM/sarama"
	"post/article/bridge"
)

//go:embed kafka.yaml
var addr string

func InitKafka() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	addrs := make([]string, 1)
	addrs[0] = addr

	client, err := sarama.NewClient(addrs, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

// NewKafkaSyncProducerForLargeMessages
// 使用 ZSTD 压缩算法的生产者
func NewKafkaSyncProducerForLargeMessages(client sarama.Client) bridge.LargeMessagesProducer {
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
func NewKafkaSyncProducerForSmallMessages(client sarama.Client) bridge.SmallMessagesProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionSnappy // sarama.CompressionLZ4

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}
