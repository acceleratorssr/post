package ioc

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/weightedroundrobin" // 加权轮询
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "post/api/proto/gen/intr/v1"
	grpc2 "post/article/grpc"
	"post/pkg/grpc_ex"
	"post/pkg/grpc_ex/interceptors/limit"
)

func InitArticleService(article *grpc2.ArticleServiceServer) *grpc_ex.Server {
	limitInterceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(limitInterceptor.BuildServerInterceptor()))
	article.RegisterServer(server)

	port := "9204"
	return grpc_ex.NewServer(server, grpc_ex.InitEtcdClient(port, "article"), port)
}

func InitLikeClient() intrv1.LikeServiceClient {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(etcdClient)
	c, err := grpc.NewClient("etcd:///service/like",
		// 注：传入的是json串, https://github.com/grpc/grpc/blob/master/doc/service_config.md
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`), // rrpick uint32
		grpc.WithResolvers(bd),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewLikeServiceClient(c)

	return remote
}
