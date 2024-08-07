package startup

import (
	_ "embed"
	"github.com/IBM/sarama"
)

//go:embed kafka.yaml
var addr string

func InitKafka() sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{addr}, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}
