package grpc_extra

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	//client *etcdv3.Client
	Port string
}

func NewServer(server *grpc.Server, port string) *Server {
	return &Server{
		Server: server,
		Port:   port,
	}
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		panic(err)
	}

	fmt.Println("grpc server run on ", s.Port)
	return s.Server.Serve(l)
}
