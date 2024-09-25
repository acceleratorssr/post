package ioc

import (
	"google.golang.org/grpc"
	"post/pkg/grpc_ex"
	"post/pkg/grpc_ex/interceptors"
	grpc2 "post/user/grpc"
)

func InitGrpcServer(user *grpc2.UserServiceServer) *grpc_ex.Server {
	interceptor := interceptors.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor()))
	user.Register(server)

	port := "9202"
	return grpc_ex.NewServer(server, grpc_ex.InitEtcdClient(port, "user"), port)
}
