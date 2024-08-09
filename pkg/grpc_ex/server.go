package grpc_ex

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":9200")
	if err != nil {
		panic(err)
	}

	fmt.Println("grpc server run on 9200")
	return s.Server.Serve(l)
}
