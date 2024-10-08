package ioc

import (
	"google.golang.org/grpc"
	"post/pkg/grpc-extra"
	"post/pkg/grpc-extra/interceptors/limit"
	grpc2 "post/sso/grpc"
)

func InitGrpcSSOServer(sso *grpc2.AuthServiceServer) *grpc_extra.Server {
	limitInterceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(limitInterceptor.BuildServerInterceptor()))
	sso.RegisterServer(server)

	port := "9203"
	return grpc_extra.NewServer(server, grpc_extra.InitEtcdClient(port, "sso"), port)
}
