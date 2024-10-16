package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc_extra "post/pkg/grpc-extra"
	loadcount "post/pkg/grpc-extra/interceptors/load-count"
	grpc2 "post/recommend/grpc"
)

func InitGrpcRecommendServer(rc *grpc2.RecommendServiceServer) *grpc_extra.Server {
	type Config struct {
		Port string `yaml:"port"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}

	loadCount := loadcount.NewLoadCount()
	ch := make(chan int, 10)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(loadCount.LoadCountInterceptor(ch)))
	rc.RegisterServer(server)

	_ = grpc_extra.InitEtcdClient(cfg.Port, "recommend", grpc_extra.WithChannel(ch))

	return grpc_extra.NewServer(server, cfg.Port)
}
