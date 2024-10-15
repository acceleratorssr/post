package CHBL

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc/attributes"
	"log"
	"reflect"
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

var R resolver.Resolver // 丑陋的包变量，但是为了手动刷新metadata，暂时没有好办法

func Fresh() {
	R.ResolveNow(resolver.ResolveNowOptions{})
}

// ResolveNow
// 调用时机：
// 初次建立连接
// 客户端重新连接：连接中断时，gRPC 客户端会调用 ResolveNow 尝试重新解析
// resolver.Resolver.ResolveNow() 手动触发解析过程
func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
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

	r.cc.UpdateState(resolver.State{Addresses: addrs}) // todo 考虑频繁调用的性能损失？突然有点慌
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
	R = r
	return r, nil
}

func (r *etcdResolver) byReflect(kvs []*mvccpb.KeyValue, addrs []resolver.Address) {
	for _, kv := range kvs {
		var node any
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			fmt.Printf("反序列化服务节点信息失败: %v", err)
			continue
		}

		addr := resolver.Address{}
		if !r.copyFields(node, addr) {
			fmt.Printf("metadata 传入的不是 struct")
		}

		addrs = append(addrs, addr)
	}
}

func (r *etcdResolver) copyFields(src interface{}, dest interface{}) bool {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest).Elem()

	if srcValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		destField := destValue.FieldByName(field.Name)

		if destField.IsValid() && destField.CanSet() {
			destField.Set(srcValue.Field(i))
		}
	}

	return true
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
