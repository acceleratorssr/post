package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	sarama_extra "post/pkg/sarama-extra"
	"post/recommend/events"
	"time"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()

	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll         // 设置 ack 策略
	saramaCfg.Producer.Compression = sarama.CompressionSnappy   // 设置压缩方式
	saramaCfg.Producer.Flush.Bytes = 16384                      // 设置批量发送大小，单位为字节
	saramaCfg.Producer.Flush.Frequency = 5 * time.Millisecond   // 设置批量发送频率
	saramaCfg.Producer.Flush.MaxMessages = 1000                 // 设置最大批量发送消息数
	saramaCfg.Producer.Flush.Messages = 10                      // 设置批次大小，当积累了 10 条消息时发送
	saramaCfg.Producer.Flush.Frequency = 500 * time.Millisecond // 设置批处理的最大等待时间，单位是毫秒

	saramaCfg.Consumer.Fetch.Min = 1                        // 设置 Consumer 一次拉取的最小字节数
	saramaCfg.Consumer.Fetch.Max = 1024 * 1024              // 设置 Consumer 一次拉取的最大字节数
	saramaCfg.Consumer.MaxWaitTime = 500 * time.Millisecond // 设置 Consumer 等待时间

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
