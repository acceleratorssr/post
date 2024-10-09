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

//func (ib *InterceptorBuilder) BuildClientInterceptor() grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
//
//	}
//}