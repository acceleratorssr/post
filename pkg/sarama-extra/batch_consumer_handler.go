package sarama_ex

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

const batchSize = 10

type ConsumerBatchHandler[T any] struct {
	fn func(msg []*sarama.ConsumerMessage, t []T) error
	// todo option 设置如下两字段
	batchSize int
	duration  time.Duration
}

func NewConsumerBatchHandler[T any](fn func(msg []*sarama.ConsumerMessage, t []T) error) *ConsumerBatchHandler[T] {
	return &ConsumerBatchHandler[T]{
		fn:        fn,
		batchSize: batchSize,
		duration:  100 * time.Second,
	}
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
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)

		// 超时控制，保证提交
		for i := 0; i < batchSize && !done; i++ {
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
