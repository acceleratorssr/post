package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userv1 "post/api/proto/gen/user/v1"
)

func InitUserClient(ec *clientv3.Client) userv1.UserServiceClient {
	type Config struct {
		Target string `json:"target"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.user", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ec)
	if err != nil {
		panic(err)
	}

	c, err := grpc.NewClient(cfg.Target,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
		grpc.WithResolvers(rs),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	return userv1.NewUserServiceClient(c)
}
