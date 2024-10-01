package ioc

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/weightedroundrobin" // 加权轮询
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/interactive/service"
	"post/internal/web/client"
)

func InitIntrGRPCClient(svc service.LikeService) intrv1.LikeServiceClient {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(etcdClient)
	c, err := grpc.NewClient("etcd:///service/like",
		// 注：传入的是json串, https://github.com/grpc/grpc/blob/master/doc/service_config.md
		//grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`), // rrpick uint32
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "weighted_round_robin": {} } ]}`), // 对应名字在balancer里
		grpc.WithResolvers(bd),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewLikeServiceClient(c)

	local := client.NewLikeServiceAdapter(svc)

	g := client.NewGreyScaleServiceAdapter(local, remote)
	g.UpdateThreshold(0) // 调整流量比例
	return g
}
