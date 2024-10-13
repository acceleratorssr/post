package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"post/pkg/grpc-extra"
	grpc2 "post/search/grpc"
)

func InitGRPCexServer(syncRpc *grpc2.SyncServiceServer,
	searchRpc *grpc2.SearchServiceServer) *grpc_extra.Server {
	type Config struct {
		Port string `yaml:"port"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	syncRpc.Register(server)
	searchRpc.Register(server)
	ec := grpc_extra.InitEtcdClient(cfg.Port, "search")
	return grpc_extra.NewServer(server, ec, cfg.Port)
}
