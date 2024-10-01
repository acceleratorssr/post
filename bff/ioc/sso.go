package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ssov1 "post/api/proto/gen/sso/v1"
)

func InitSSOClient(ec *clientv3.Client) ssov1.AuthServiceClient {
	type Config struct {
		Target string `json:"target"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.sso", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ec)
	if err != nil {
		panic(err)
	}
	// todo grpc 加密
	c, err := grpc.NewClient(cfg.Target,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
		grpc.WithResolvers(rs),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	return ssov1.NewAuthServiceClient(c)
}
