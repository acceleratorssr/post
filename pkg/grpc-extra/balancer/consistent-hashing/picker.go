package consistent_hashing

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"hash/fnv"
	"sync"
)

const ConsistentHash = "consistent_hash"

type grpcClientReq[req any, resp any] func(context.Context, req, ...grpc.CallOption) (resp, error)

// RegisterWithKey key为用户指定hash的特征key
func RegisterWithKey[request any, response any](ctx context.Context, key string, req request, client grpcClientReq[request, response]) (response, error) {
	ctx = context.WithValue(ctx, "hash_key", key)

	return client(ctx, req)
}

type HashRing struct {
	mutex sync.Mutex
	nodes []balancer.SubConn
}

func (h *HashRing) Add(node balancer.SubConn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.nodes = append(h.nodes, node)
}

func (h *HashRing) GetNode(key string) balancer.SubConn {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(h.nodes) == 0 {
		return nil
	}

	hash := fnv.New32a()
	hash.Write([]byte(key))
	hashValue := hash.Sum32() % uint32(len(h.nodes))
	// todo 上k8s收集换掉
	fmt.Printf("%s -> server:%d \n", key, hashValue)
	fmt.Println(len(h.nodes))
	return h.nodes[hashValue]
}

type ConsistentHashPicker struct {
	hashRing *HashRing
	//connectionPool *ConnectionPool
}

func (p *ConsistentHashPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.hashRing.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// 如果没设置key，则和pick_first等效
	key, _ := info.Ctx.Value("hash_key").(string)

	selectedNode := p.hashRing.GetNode(key)
	if selectedNode == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	//conn, err := p.connectionPool.GetConnection(selectedNode.)
	//if err != nil {
	//	return balancer.PickResult{}, err
	//}

	return balancer.PickResult{
		SubConn: selectedNode,
		// 暂时没想到有什么特别需要处理的
		// 发生错误时为了保证数据一致性，所以改变节点的话可能引入复杂性
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

//// ConnectionPool 用于流式gRPC建立连接
//type ConnectionPool struct {
//	mu          sync.Mutex
//	connections map[string]*grpc.ClientConn
//}
//
//func NewConnectionPool() *ConnectionPool {
//	return &ConnectionPool{
//		connections: make(map[string]*grpc.ClientConn),
//	}
//}
//
//func (pool *ConnectionPool) GetConnection(address string) (*grpc.ClientConn, error) {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()
//
//	if conn, exists := pool.connections[address]; exists {
//		return conn, nil
//	}
//	conn, err := grpc.NewClient("etcd:///service/sso",
//		grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		return nil, err
//	}
//
//	pool.connections[address] = conn
//	return conn, nil
//}
//
//func (pool *ConnectionPool) CloseConnections() {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()
//
//	for _, conn := range pool.connections {
//		conn.Close()
//	}
//	pool.connections = make(map[string]*grpc.ClientConn)
//}
