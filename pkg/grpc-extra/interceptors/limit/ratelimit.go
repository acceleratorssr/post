package limit

import (
	"google.golang.org/grpc"
)

type Opts func(*InterceptorBuilder)

type InterceptorBuilder struct {
	interval int64
	capacity int
}

func NewInterceptorBuilder(opts ...Opts) *InterceptorBuilder {
	ib := &InterceptorBuilder{
		interval: 100,
		capacity: 10,
	}
	for _, opt := range opts {
		opt(ib)
	}
	return ib
}

func (ib *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return NewTokenBucketLimit(ib.interval, ib.capacity).NewServerInterceptor()
}

func WithInterval(interval int64) Opts {
	return func(ib *InterceptorBuilder) {
		ib.interval = interval
	}
}

func WithCapacity(capacity int) Opts {
	return func(ib *InterceptorBuilder) {
		ib.capacity = capacity
	}
}
