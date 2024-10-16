package ioc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/pkg/grpc-extra"
	CHBL "post/pkg/grpc-extra/balancer/CHBL"
	"post/pkg/grpc-extra/interceptors/limit"
	grpc2 "post/user/grpc"
	"time"
)

func InitGrpcServer(user *grpc2.UserServiceServer) *grpc_extra.Server {
	interceptor := limit.NewInterceptorBuilder()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.BuildServerInterceptor())) // 第二开启
	user.Register(server)

	port := "9206"
	_ = grpc_extra.InitEtcdClient(port, "user") // 第三
	return grpc_extra.NewServer(server, port)
}

func InitGrpcSSOClient() ssov1.AuthServiceClient {
	// 可监听 sso 服务节点
	serviceKey := "service/sso"

	etcdResolver, err := CHBL.NewEtcdResolver([]string{"127.0.0.1:12379"})
	if err != nil {
		panic(err)
	}
	resolver.Register(etcdResolver) // 全局

	c, err := grpc.NewClient("etcd:///"+serviceKey,
		//grpc.WithResolvers(etcdResolver), // 局部
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "consistent_hashing_with_bounded_loads": {} } ]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	// 最先开启
	go func() {
		ticker := time.NewTicker(time.Second * 11) // 每11秒更新一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				CHBL.Fresh()
			}
		}
	}()

	return ssov1.NewAuthServiceClient(c)
}
