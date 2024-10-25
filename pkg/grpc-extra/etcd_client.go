package grpc_extra

import (
	"context"
	"fmt"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"post/pkg/net-extra"
	"time"
)

const (
	RequestsCount = "requests_count"
	Weight        = "weight"
)

type Opts func(*etcdClient)

type etcdClient struct {
	client   *etcdv3.Client
	e        endpoints.Manager
	EtcdAddr []string
	Port     string
	key      string
	TTL      int64
	name     string
	ip       string

	RequestsCount int
	metadata      *ServiceMetadata // 可改进为传结构体
	needMetadata  bool

	ch chan int // 目前用于传递请求数
}

type ServiceMetadata struct {
	Weight        int `json:"weight"`
	RequestsCount int `json:"requests_count"` // 统计请求数量
}

func InitEtcdClient(port string, name string, opts ...Opts) *etcdv3.Client {
	c := etcdClient{
		Port:     port,
		EtcdAddr: []string{"127.0.0.1:12379"},
		TTL:      12,
		name:     name,
	}
	c.initIp()
	c.initService()
	c.initEtcdClient(opts...)

	return c.client
}

func (ec *etcdClient) initIp() {
	ec.ip = net_extra.GetOutboundIP()
}

func (ec *etcdClient) initService() {
	ec.key = "service"
}

func (ec *etcdClient) initEtcdClient(opts ...Opts) {
	for _, opt := range opts {
		opt(ec)
	}

	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: ec.EtcdAddr,
	})
	if err != nil {
		panic(err)
	}

	pCtx := context.Background()
	addr := ec.ip + ":" + ec.Port

	e, err := endpoints.NewManager(client, ec.key+"/"+ec.name)
	if err != nil {
		panic(err)
	}
	ec.e = e

	ctx, cancel := context.WithTimeout(pCtx, time.Second)
	grant, err := client.Grant(ctx, ec.TTL)
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithCancel(pCtx)
	defer cancel()

	ep := endpoints.Endpoint{
		Addr:     addr,
		Metadata: ec.metadata,
	}

	err = ec.e.AddEndpoint(ctx, ec.key+"/"+ec.name+"/"+addr, ep, etcdv3.WithLease(grant.ID))
	if err != nil {
		panic(err)
	}

	// 续约
	go func() {
		// todo （latest）写个允许控制间隔的方法，注意考虑并发情况
		alive, err := client.KeepAlive(pCtx, grant.ID)
		if err != nil {
			panic(err)
		}
		for range alive {
			// log
		}
	}()

	if ec.needMetadata {
		go func() {
			ec.updateMetadataPeriodically(ec.ch, grant.ID)
		}()
	}

	ec.client = client
	return
}

func (ec *etcdClient) updateMetadataPeriodically(ch chan int, id etcdv3.LeaseID) {
	ticker := time.NewTicker(time.Second * 10) // 每10秒更新一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			addr := ec.ip + ":" + ec.Port
			key := ec.key + "/" + ec.name + "/" + addr
			ec.RequestsCount = <-ch
			ctx := context.Background()

			epMap, err := ec.e.List(ctx)
			if err != nil {
				fmt.Printf("获取 endpoints 列表失败: %v\n", err)
				continue
			}
			cur, exist := epMap[key]
			if exist && cur.Metadata.(map[string]any)[RequestsCount] != nil && cur.Metadata.(map[string]any)[RequestsCount].(float64) == float64(ec.RequestsCount) { // 期间无请求，跳过更新metadata
				continue
			}

			ec.metadata.RequestsCount = ec.RequestsCount
			err = ec.e.AddEndpoint(ctx, key, endpoints.Endpoint{
				Addr:     addr,
				Metadata: ec.metadata,
			}, etcdv3.WithLease(id))
			if err != nil {
				fmt.Printf("更新 endpoint 失败: %v\n", err)
				continue
			}
		}
	}
}

func (ec *etcdClient) ShoutDown() {
	if ec.e != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := ec.e.DeleteEndpoint(ctx, ec.key+"/"+ec.name+"/"+ec.ip+":"+ec.Port)
		if err != nil {
			// log
		}
		cancel()
	}
	if ec.client == nil {
		return
	}
	err := ec.client.Close()
	if err != nil {
		//log
	}
}

func WithChannel(ch chan int) Opts {
	return func(ec *etcdClient) {
		ec.ch = ch
	}
}

func WithNeedMetadata() Opts {
	return func(ec *etcdClient) {
		ec.needMetadata = true
	}
}

func WithMetadata(weight int) Opts {
	return func(ec *etcdClient) {
		ec.metadata = &ServiceMetadata{
			Weight: weight,
		}
	}
}
