package ioc

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
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
	return grpc_extra.NewServer(server, grpc_extra.InitEtcdClient(port, "user"), port)
}

func InitGrpcSSOClient() ssov1.AuthServiceClient {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(etcdClient)
	c, err := grpc.NewClient("etcd:///service/sso",
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "consistent_hash": {} } ]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return ssov1.NewAuthServiceClient(c)
}
