package grpc_ex

import (
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

	return s.Server.Serve(l)
}
