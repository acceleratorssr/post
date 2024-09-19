package events

import (
	"context"
	"github.com/IBM/sarama"
	"post/pkg/logger"
	"post/pkg/sarama_ex"
	"post/search/service"
	"time"
)

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.Logger
}

func (a *SyncDataEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data",
		a.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncArticle},
			sarama_ex.NewHandler[SyncDataEvent](a.Consume))
		if err != nil {
			a.l.Error("发生错误，退出消费循环 ", logger.Error(err))
		}
	}()
	return err
}

func (a *SyncDataEventConsumer) Consume(sg *sarama.ConsumerMessage,
	evt SyncDataEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return a.svc.InputAny(ctx, evt.IndexName, evt.DocID, evt.Data)
}
