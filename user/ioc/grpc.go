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

// InitGrpcSSOClient todo bug:每次启动user服务后，第一个请求总是1或者2，第二个请求及以后才为正常的3节点
func InitGrpcSSOClient() ssov1.AuthServiceClient {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	if err != nil {
		panic(err)
	}

	//// 监听 sso 服务节点
	//serviceKey := "service/sso"
	//readyChan := make(chan struct{})
	//
	//go func() {
	//	ssoCount := 0
	//	for {
	//		rch := etcdClient.Watch(context.Background(), serviceKey, etcdv3.WithPrefix())
	//
	//		for wresp := range rch {
	//			for _, ev := range wresp.Events {
	//				if ev.Type == etcdv3.EventTypePut {
	//					ssoCount++
	//				} else if ev.Type == etcdv3.EventTypeDelete {
	//					ssoCount--
	//				}
	//			}
	//
	//			// 检查当前可用的 sso 服务数量
	//			if ssoCount >= 3 {
	//				readyChan <- struct{}{}
	//				return
	//			}
	//		}
	//	}
	//}()
	//
	//// 防止节点 未就绪
	//<-readyChan

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
