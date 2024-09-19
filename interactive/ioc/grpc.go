package ioc

import (
	"google.golang.org/grpc"
	grpc2 "post/interactive/grpc"
	"post/pkg/grpc_ex"
	"post/pkg/grpc_ex/interceptors"
)

func InitGRPCexServer(intr *grpc2.LikeServiceServer) *grpc_ex.Server {
	interceptor := interceptors.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor())) // todo 可加可观测性记录，可约定grpc的header携带数据，38m
	intr.Register(server)

	port := "9201"
	return grpc_ex.NewServer(server, grpc_ex.InitEtcdClient(port, "like"), port)
}
