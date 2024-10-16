package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	sarama_extra "post/pkg/sarama-extra"
	"post/recommend/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()

	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewConsumers(articleConsumer *events.BatchKafkaRecommendConsumer) []sarama_extra.Consumer {
	return []sarama_extra.Consumer{
		articleConsumer,
	}
}
