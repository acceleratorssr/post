package grpc

// 使用etcd作为注册中心，为grpc提供服务

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/pkg/net_ex"
	"testing"
	"time"
)

type EtcdTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = client
}

func (s *EtcdTestSuite) TestClient() {
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	// 此处etcd代表提供的resolver
	cc, err := grpc.NewClient("etcd:///service/test",
		grpc.WithResolvers(bd),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := intrv1.NewLikeServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Like(ctx, &intrv1.LikeRequest{
		ObjType: "article",
		ObjID:   1,
		Uid:     99,
	})
	require.NoError(s.T(), err)
	s.T().Log(resp)
}

func (s *EtcdTestSuite) TestServer() {
	addr := net_ex.GetOutboundIP()
	key := "service/test"

	l, err := net.Listen("tcp", ":9400")
	require.NoError(s.T(), err)

	// 以服务为维度，一个服务一个manager
	e, err := endpoints.NewManager(s.client, key)
	require.NoError(s.T(), err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	grant, err := s.client.Grant(ctx, 12)
	cancel()
	require.NoError(s.T(), err)

	// 定期调用AddEndpoint或者update即可更新addr和metadata
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	// 准备工作完成后，再注册服务
	// key 为实例的key，用instance id，或者本机IP+Port
	err = e.AddEndpoint(ctx, key+"/"+addr+":9400", endpoints.Endpoint{
		Addr: addr + ":9400",
		// metadata可传元数据
	}, etcdv3.WithLease(grant.ID))
	require.NoError(s.T(), err)

	// 续约
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		// 默认开始续约时间为ttl/3
		// 也可自己循环调用KeepAliveOnce()，自己控制间隔
		// todo （latest）写个允许控制间隔的方法，需要考虑并发情况
		alive, err := s.client.KeepAlive(ctx, grant.ID)
		cancel()
		require.NoError(s.T(), err)
		for range alive {
			s.T().Log(alive)
		}
	}()

	server := grpc.NewServer()
	intrv1.RegisterLikeServiceServer(server, &LikeServiceServer{})
	server.Serve(l)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	e.DeleteEndpoint(ctx, key+"/"+addr+":9400")
	server.GracefulStop()
}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}
