package fixer

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"post/migrator"
	"post/migrator/events"
	"post/migrator/fixer"
	"post/pkg/sarama_ex"
	"time"
)

type Consumer[T migrator.Entity] struct {
	client       sarama.Client
	dependBase   *fixer.Fixer[T]
	dependTarget *fixer.Fixer[T]
	topic        string
}

func NewConsumer[T migrator.Entity](client sarama.Client, base, target *gorm.DB, topic string) *Consumer[T] {
	return &Consumer[T]{
		client:       client,
		dependBase:   fixer.NewFixerDependBase[T](base, target),
		dependTarget: fixer.NewFixerDependTarget[T](base, target),
		topic:        topic,
	}
}

// Start 这边就是自己启动 goroutine 了
func (r *Consumer[T]) Start(topic string) error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator-fix",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{r.topic},
			sarama_ex.NewHandler[events.InconsistentEvent](r.Consume))
		if err != nil {
			// log
		}
	}()
	return err
}

func (r *Consumer[T]) Consume(msg *sarama.ConsumerMessage, t events.InconsistentEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch t.Direction {
	case "base":
		return r.dependBase.Fix(ctx, t)
	case "target":
		return r.dependTarget.Fix(ctx, t)
	}
	return errors.New("未知的校验方向")
}
