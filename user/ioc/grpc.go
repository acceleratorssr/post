package ioc

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	ecResolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/pkg/grpc-extra"
	"post/pkg/grpc-extra/interceptors/limit"
	grpc2 "post/user/grpc"
)

func InitGrpcServer(user *grpc2.UserServiceServer) *grpc_extra.Server {
	interceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor()))
	user.Register(server)

	port := "9202"
	ec := grpc_extra.InitEtcdClient(port, "user")
	return grpc_extra.NewServer(server, ec, port)
}

// InitGrpcSSOClient todo bug:每次启动user服务后，第一个请求总是1或者2，第二个请求及以后才为正常的3节点
func InitGrpcSSOClient() ssov1.AuthServiceClient {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	if err != nil {
		panic(err)
	}
	// 可监听 sso 服务节点
	serviceKey := "service/sso"

	bd, err := ecResolver.NewBuilder(etcdClient)
	c, err := grpc.NewClient("etcd:///"+serviceKey,
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "consistent_hash": {} } ]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return ssov1.NewAuthServiceClient(c)
}
