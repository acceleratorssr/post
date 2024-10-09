package consistent_hashing

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(ConsistentHash, &ConsistentHashPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
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
	fmt.Printf("%s -> server:%d\n", key, hashValue)
	return h.nodes[hashValue]
}

type ConsistentHashPicker struct {
	hashRing *HashRing
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

	return balancer.PickResult{
		SubConn: selectedNode,
		// 暂时没想到有什么特别需要处理的
		// 发生错误时为了保证数据一致性，所以改变节点的话可能引入复杂性
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type ConsistentHashPickerBuilder struct {
}

func (b *ConsistentHashPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	hashRing := &HashRing{}
	for con := range info.ReadySCs {
		hashRing.Add(con)
	}
	return &ConsistentHashPicker{hashRing: hashRing}
}
