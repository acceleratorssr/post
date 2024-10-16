package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	sarama_extra "post/pkg/sarama-extra"
	"post/search/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true

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

// NewConsumers 依旧是所有的 Consumer 在这里注册一下
func NewConsumers(articleConsumer *events.ArticleConsumer) []sarama_extra.Consumer {
	return []sarama_extra.Consumer{
		articleConsumer,
	}
}
