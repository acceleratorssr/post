package ioc

import (
	"google.golang.org/grpc"
	grpc2 "post/interactive/grpc"
	"post/pkg/grpc-extra"
	"post/pkg/grpc-extra/interceptors/limit"
)

func InitGRPCexServer(intr *grpc2.LikeServiceServer) *grpc_extra.Server {
	interceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor())) // todo 可加可观测性记录，可约定grpc的header携带数据，38m
	intr.Register(server)

	port := "9201"
	return grpc_extra.NewServer(server, grpc_extra.InitEtcdClient(port, "like"), port)
}
