package ioc

import (
	_ "embed"
	"github.com/IBM/sarama"
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
