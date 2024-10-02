package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	articlev1 "post/api/proto/gen/article/v1"
)

func InitArticleClient(ec *clientv3.Client) articlev1.ArticleServiceClient {
	type Config struct {
		Target string `json:"target"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.article", &cfg)
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

	return articlev1.NewArticleServiceClient(c)
}
