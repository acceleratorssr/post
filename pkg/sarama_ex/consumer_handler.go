package sarama_ex

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

// todo 添加prometheus监控

// Handler 实现consumerGroupHandler接口
type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{fn: fn}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			// 数据有问题
		}

		err = h.fn(msg, t)
		if err != nil {
			// 记录消费失败消息，重试fn
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
