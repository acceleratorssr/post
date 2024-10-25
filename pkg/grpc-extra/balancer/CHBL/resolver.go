package CHBL

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc/attributes"
	"log"
	grpc_extra "post/pkg/grpc-extra"
	"reflect"
	"time"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdResolver struct {
	client *clientv3.Client
	cc     resolver.ClientConn // 主要获取 UpdateState(State) error 进行回调
	target resolver.Target

	curAddrs map[string]bool
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
// 可通过watch机制监听服务节点的变化
func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	resp, err := r.client.Get(context.Background(), r.target.Endpoint(), clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("解析服务失败: %v", err)
	}

	freshFlag := false
	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		var node any
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			log.Printf("反序列化服务节点信息失败: %v", err)
			continue
		}
		m := node.(map[string]any)
		md := m["Metadata"].(map[string]any)

		metadata := attributes.New(grpc_extra.Weight, md[grpc_extra.Weight])
		metadata = metadata.WithValue(grpc_extra.RequestsCount, md[grpc_extra.RequestsCount])

		addr := resolver.Address{
			Addr:       m["Addr"].(string),
			Attributes: metadata,
		}
		addrs = append(addrs, addr)

		// 节点列表发生变化
		if v, ok := r.curAddrs[addr.Addr]; !ok || v == false {
			freshFlag = true
		} else {
			r.curAddrs[addr.Addr] = true
		}
	}

	if freshFlag {
		r.cc.UpdateState(resolver.State{
			Addresses: addrs, // 会被转为 Endpoints
		})
	} else {
		// 通知 picker 更新节点负载

	}
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
	go r.watchServiceNodes()
	return r, nil
}

// watchServiceNodes 启动对 etcd 中服务节点的监听
func (r *etcdResolver) watchServiceNodes() {
	watchChan := r.client.Watch(context.Background(), r.target.Endpoint(), clientv3.WithPrefix())

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT, mvccpb.DELETE: // etcdv3.EventTypePut 也行
				r.ResolveNow(resolver.ResolveNowOptions{})
			}
		}
	}
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
		client:   cli,
		curAddrs: make(map[string]bool),
	}, nil
}
