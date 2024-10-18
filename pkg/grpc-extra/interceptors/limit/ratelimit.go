package limit

import (
	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
}

func NewInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{}
}

func (ib *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return NewTokenBucketLimit(100, 10).NewServerInterceptor()
}
