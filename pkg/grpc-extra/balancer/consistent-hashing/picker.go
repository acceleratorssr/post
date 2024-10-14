package consistent_hashing

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"hash/fnv"
	"sort"
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
	mutex       sync.Mutex
	hashNodes   []uint32                    // 节点的哈希值
	nodeMap     map[uint32]balancer.SubConn // 哈希值对应的实际节点
	virtualNode int                         // 每个物理节点的虚拟节点数量
}

func NewHashRing(virtualNode int) *HashRing {
	return &HashRing{
		hashNodes:   make([]uint32, 0),
		nodeMap:     make(map[uint32]balancer.SubConn),
		virtualNode: virtualNode,
	}
}

func (h *HashRing) AddNode(node balancer.SubConn, id int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualNode; i++ {
		hash := h.hashKey(fmt.Sprintf("%d#%d", id, i))
		h.hashNodes = append(h.hashNodes, hash)
		h.nodeMap[hash] = node
	}
	sort.Slice(h.hashNodes, func(i, j int) bool { return h.hashNodes[i] < h.hashNodes[j] })
}

func (h *HashRing) RemoveNode(id int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualNode; i++ {
		hash := h.hashKey(fmt.Sprintf("%d#%d", id, i))

		for j := 0; j < len(h.hashNodes); j++ {
			if h.hashNodes[j] == hash {
				h.hashNodes = append(h.hashNodes[:j], h.hashNodes[j+1:]...)
				break
			}
		}

		delete(h.nodeMap, hash)
	}

	sort.Slice(h.hashNodes, func(i, j int) bool { return h.hashNodes[i] < h.hashNodes[j] })
}

func (h *HashRing) hashKey(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}

func (h *HashRing) GetNode(key string) balancer.SubConn {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(h.hashNodes) == 0 {
		return nil
	}
	hash := h.hashKey(key)
	idx := sort.Search(len(h.hashNodes), func(i int) bool { return h.hashNodes[i] >= hash })

	// 如果没有找到合适的节点，返回第一个节点（回到环的起点）
	if idx == len(h.hashNodes) {
		idx = 0
	}

	// todo 上k8s收集换掉
	fmt.Printf("%s -> server:%d len:%d\n", key, idx, len(h.hashNodes))
	return h.nodeMap[h.hashNodes[idx]]
}

type ConsistentHashPicker struct {
	hashRing *HashRing
	mutex    sync.Mutex
}

func (p *ConsistentHashPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 如果没设置key，则类似于 pick_first，连接环的第一个节点
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
