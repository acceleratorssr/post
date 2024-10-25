package load_count

import (
	"context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type LoadCount struct {
	RequestCount        atomic.Int32
	TotalProcessingTime time.Duration
}

type Opt func(*LoadCount)

func NewLoadCount(opt ...Opt) *LoadCount {
	l := &LoadCount{}
	for _, o := range opt {
		o(l)
	}
	return l
}

func (l *LoadCount) LoadCountInterceptor(ch chan int) grpc.UnaryServerInterceptor {
	go func() {
		ticker := time.NewTicker(time.Second * 10) // 每10秒更新一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ch <- int(l.RequestCount.Load())
			}
		}
	}()
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		l.RequestCount.Add(1)
		return handler(ctx, req)
	}
}
