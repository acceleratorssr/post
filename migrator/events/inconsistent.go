package events

import "context"

type InconsistentProducer interface {
	InconsistentEventProducer(ctx context.Context, event *InconsistentEvent) error
}

type InconsistentEvent struct {
	ID        int64
	Direction string // 确定参考表
	Type      string // 确定事件类型
}

const (
	InconsistentEventTypeNoEquals = "NoEquals"
	InconsistentEventTypeNoExist  = "NoExist"
)
