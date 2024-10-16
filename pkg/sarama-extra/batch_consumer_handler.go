package sarama_extra

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

const (
	DefaultBatchSize = 50
	DefaultDuration  = 10 * time.Second
)

type Opt[T any] func(*ConsumerBatchHandler[T])

type ConsumerBatchHandler[T any] struct {
	fn        func(msg []*sarama.ConsumerMessage, t []T) error
	batchSize int
	duration  time.Duration
}

// NewConsumerBatchHandler
// 默认50条消息，10秒超时
func NewConsumerBatchHandler[T any](fn func(msg []*sarama.ConsumerMessage, t []T) error, opts ...Opt[T]) *ConsumerBatchHandler[T] {
	c := &ConsumerBatchHandler[T]{
		fn:        fn,
		batchSize: DefaultBatchSize,
		duration:  DefaultDuration,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *ConsumerBatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerBatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerBatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.duration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, c.batchSize)
		ts := make([]T, 0, c.batchSize)

		// 该批次的获取时间超过 c.duration 时，直接提交
		for i := 0; i < c.batchSize && !done; i++ {
			select {
			case msg, ok := <-claim.Messages():
				if !ok {
					cancel()
					return nil
				}

				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					// log
					session.MarkMessage(msg, "")
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			case <-ctx.Done():
				done = true
			}

		}

		//if len(msgs) == 0 {
		//	continue
		//}
		err := c.fn(msgs, ts)
		if err == nil {
			// log
			for _, msg := range msgs {
				session.MarkMessage(msg, "")
			}
		} else {
			// 重试
		}
		cancel()
	}

}

func WithBatchSize[T any](size int) Opt[T] {
	return func(c *ConsumerBatchHandler[T]) {
		c.batchSize = size
	}
}

func WithDuration[T any](d time.Duration) Opt[T] {
	return func(c *ConsumerBatchHandler[T]) {
		c.duration = d
	}
}
