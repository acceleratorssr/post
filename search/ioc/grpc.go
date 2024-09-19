package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"post/pkg/grpc_ex"
	grpc2 "post/search/grpc"
)

func InitGRPCxServer(syncRpc *grpc2.SyncServiceServer,
	searchRpc *grpc2.SearchServiceServer) *grpc_ex.Server {
	type Config struct {
		Port string `yaml:"port"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	syncRpc.Register(server)
	searchRpc.Register(server)
	return grpc_ex.NewServer(server, grpc_ex.InitEtcdClient(cfg.Port, "search"), cfg.Port)
}
