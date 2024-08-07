package main

import (
	"google.golang.org/grpc"
	"net"
	intrv1 "post/api/proto/gen/intr/v1"
	grpc2 "post/interactive/grpc"
)

func main() {
	server := grpc.NewServer()
	intrSvc := &grpc2.LikeServiceServer{}
	intrv1.RegisterLikeServiceServer(server, intrSvc)

	l, err := net.Listen("tcp", ":9200")
	if err != nil {
		panic(err)
	}

	err = server.Serve(l)
}
