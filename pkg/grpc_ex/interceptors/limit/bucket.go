package limit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type TokenBucketLimit struct {
	intervalPerSecond int64
	buckets           chan struct{}
	closeCh           chan struct{}
}

func NewTokenBucketLimit(interval int64, capacity int) *TokenBucketLimit {
	return &TokenBucketLimit{
		intervalPerSecond: interval,
		buckets:           make(chan struct{}, capacity),
	}
}

func (l *TokenBucketLimit) NewServerInterceptor() grpc.UnaryServerInterceptor {
	interval := time.Second / time.Duration(l.intervalPerSecond)
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-l.closeCh:
				return
			case <-ticker.C:
				select {
				case l.buckets <- struct{}{}:

				default:

				}
			}
		}
	}()
	// 此处可调整限流的粒度
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		select {
		case <-l.buckets:
			return handler(ctx, req)
		default: // 并发越高，越不能阻塞，如果不高，此处也可以阻塞等待，直到超时
			return nil, errors.New("too many requests")
		}
	}
}

// Close 仅可调用一次
func (l *TokenBucketLimit) Close() {
	close(l.closeCh)
}
