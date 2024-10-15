package CHBL

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/attributes"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdResolver struct {
	client *clientv3.Client
	cc     resolver.ClientConn
	target resolver.Target
}

type NodeValue struct {
	Op       int    `json:"op"`
	Addr     string `json:"addr"`
	Metadata int    `json:"metadata"` // RequestCount
}

// ResolveNow
// 调用时机：
// 初次建立连接
// 客户端重新连接：连接中断时，gRPC 客户端会调用 ResolveNow 尝试重新解析
// resolver.Resolver.ResolveNow() 手动触发解析过程
func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	fmt.Println(r.target.Endpoint())
	resp, err := r.client.Get(context.Background(), r.target.Endpoint(), clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("解析服务失败: %v", err)
	}

	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		var node NodeValue
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			log.Printf("反序列化服务节点信息失败: %v", err)
			continue
		}

		// 暂时只存个请求次数
		addr := resolver.Address{
			Addr:       node.Addr,
			Attributes: attributes.New("request_count", node.Metadata),
		}

		addrs = append(addrs, addr)
	}

	r.cc.UpdateState(resolver.State{Addresses: addrs}) // todo 考虑频繁调用的性能损失？
}

// Close 命名服务关闭了，故关闭本地客户端
func (r *etcdResolver) Close() {
	_ = r.client.Close()
}

func (r *etcdResolver) Scheme() string {
	return "etcd"
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	r.target = target
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func NewEtcdResolver(etcdEndpoints []string) (resolver.Builder, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &etcdResolver{
		client: cli,
	}, nil
}
