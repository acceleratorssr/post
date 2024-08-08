package ioc

import (
	"google.golang.org/grpc"
	grpc2 "post/interactive/grpc"
	"post/pkg/grpc_ex"
)

func InitGRPCexServer(intr *grpc2.LikeServiceServer) *grpc_ex.Server {
	server := grpc.NewServer()
	intr.Register(server)
	return &grpc_ex.Server{
		Addr:   ":9200",
		Server: server,
	}
}
