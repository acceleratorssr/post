package ioc

import (
	_ "embed"
	"github.com/IBM/sarama"
	"post/interactive/events"
	"post/interactive/repository/dao"
	"post/migrator/events/fixer"
	"post/pkg/sarama-extra"
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

func NewKafkaConsumer(consumer *events.KafkaReadConsumer,
	fix *fixer.Consumer[dao.Like]) []sarama_extra.Consumer {
	return []sarama_extra.Consumer{consumer, fix}
}

//func NewKafkaConsumer(consumer *events2.BatchKafkaConsumer) []events2.Consumer {
//	return []events2.Consumer{consumer}
//}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}
