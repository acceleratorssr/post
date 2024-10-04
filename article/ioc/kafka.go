package ioc

import (
	_ "embed"
	"github.com/IBM/sarama"
	"post/article/events"
	"post/pkg/sarama_ex"
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

func NewKafkaSyncProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func NewKafkaConsumer(consumer *events.KafkaPublishedConsumer) []sarama_ex.Consumer {
	return []sarama_ex.Consumer{
		consumer,
	}
}

//func NewKafkaConsumer(consumer *events.KafkaConsumer) []events.Consumer {
//	return []events.Consumer{consumer}
//}

//func NewKafkaConsumer(consumer *events2.BatchKafkaConsumer) []events2.Consumer {
//	return []events2.Consumer{consumer}
//}
