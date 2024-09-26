package ioc

import (
	"google.golang.org/grpc"
	"post/pkg/grpc_ex"
	"post/pkg/grpc_ex/interceptors/limit"
	grpc2 "post/sso/grpc"
)

func InitGrpcSSOServer(sso *grpc2.AuthServiceServer) *grpc_ex.Server {
	interceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor()))
	sso.RegisterServer(server)

	port := "9203"
	return grpc_ex.NewServer(server, grpc_ex.InitEtcdClient(port, "sso"), port)
}
