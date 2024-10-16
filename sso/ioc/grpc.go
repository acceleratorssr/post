package ioc

import (
	"google.golang.org/grpc"
	"post/pkg/grpc-extra"
	"post/pkg/grpc-extra/interceptors/limit"
	loadcount "post/pkg/grpc-extra/interceptors/load-count"
	grpc2 "post/sso/grpc"
)

func InitGrpcSSOServer(sso *grpc2.AuthServiceServer) *grpc_extra.Server {
	limitInterceptor := limit.NewInterceptorBuilder()
	loadCount := loadcount.NewLoadCount()
	ch := make(chan int, 10)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(limitInterceptor.BuildServerInterceptor(), loadCount.LoadCountInterceptor(ch)))
	sso.RegisterServer(server)

	port := "9205"
	_ = grpc_extra.InitEtcdClient(port, "sso", grpc_extra.WithChannel(ch))

	return grpc_extra.NewServer(server, port)
}
