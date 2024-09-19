package grpc_ex

import (
	"context"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"post/pkg/net_ex"
	"time"
)

type etcdClient struct {
	client   *etcdv3.Client
	e        endpoints.Manager
	EtcdAddr []string
	Port     string
	key      string
	TTL      int64
	name     string
	ip       string
}

func InitEtcdClient(port string, name string) *etcdv3.Client {
	c := etcdClient{
		Port:     port,
		EtcdAddr: []string{"127.0.0.1:12379"},
		TTL:      12,
		name:     name,
	}
	c.initIp()
	c.initService()
	c.initEtcdClient()

	return c.client
}

func (ec *etcdClient) initIp() {
	ec.ip = net_ex.GetOutboundIP()
}

func (ec *etcdClient) initService() {
	ec.key = "service"
}
func (ec *etcdClient) initEtcdClient() {
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
	err = ec.e.AddEndpoint(ctx, ec.key+"/"+ec.name+"/"+addr, endpoints.Endpoint{
		Addr: addr,
		// metadata可传元数据
	}, etcdv3.WithLease(grant.ID))
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

	ec.client = client
	return
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
