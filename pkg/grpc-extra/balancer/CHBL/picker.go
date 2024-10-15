package CHBL

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"hash/fnv"
	"math"
	"reflect"
	"sort"
	"sync"
)

const CHBL = "consistent_hashing_with_bounded_loads"

type grpcClientReq[req any, resp any] func(context.Context, req, ...grpc.CallOption) (resp, error)

type Opt func(*HashRing)

// RegisterWithKey key为用户指定hash的特征key
func RegisterWithKey[request any, response any](ctx context.Context, key string, req request, client grpcClientReq[request, response]) (response, error) {
	ctx = context.WithValue(ctx, "hash_key", key)

	return client(ctx, req)
}

type HashRing struct {
	mutex         sync.Mutex
	virtualNodes  []uint32                    // 节点的哈希值
	virtualToReal map[uint32]balancer.SubConn // 记录虚拟节点到真实节点的映射
	virtualNode   int                         // 每个物理节点的虚拟节点数量
	load          map[uint32]float64          // 真实节点的负载，目前简单期间仅存请求数量作为标准							-- 有界一致性哈希算法
	avgLoad       int                         // 记录节点的平均允许负载 											 	-- 有界一致性哈希算法
	c             float64                     // 平衡参数，取1.25~2，默认为1.25										-- 有界一致性哈希算法
}

func (h *HashRing) AddNode(node balancer.SubConn, id uint32) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var hash uint32
	for i := 0; i < h.virtualNode; i++ {
		hash = h.HashKey(fmt.Sprintf("%d#%d", id, i))
		h.virtualNodes = append(h.virtualNodes, hash)
		h.virtualToReal[hash] = node
	}
	sort.Slice(h.virtualNodes, func(i, j int) bool { return h.virtualNodes[i] < h.virtualNodes[j] })
}

func (h *HashRing) RemoveNode(id uint32) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualNode; i++ {
		hash := h.HashKey(fmt.Sprintf("%d#%d", id, i))

		for j := 0; j < len(h.virtualNodes); j++ {
			if h.virtualNodes[j] == hash {
				h.virtualNodes = append(h.virtualNodes[:j], h.virtualNodes[j+1:]...)
				break
			}
		}

		delete(h.virtualToReal, hash)
		delete(h.load, hash)
	}

	sort.Slice(h.virtualNodes, func(i, j int) bool { return h.virtualNodes[i] < h.virtualNodes[j] })
}

func (h *HashRing) HashKey(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}

func (h *HashRing) GetNode(key string) balancer.SubConn {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(h.virtualNodes) == 0 {
		return nil
	}

	idx := 0
	hash := h.HashKey(key)
	idx = sort.Search(len(h.virtualNodes), func(i int) bool { return h.virtualNodes[i] >= hash })
	// 如果没有找到合适的节点，返回第一个节点（回到环的起点）
	if idx == len(h.virtualNodes) {
		idx = 0
	}

	reqIdx := idx
	realNode := h.virtualToReal[h.virtualNodes[reqIdx]]

	for {
		// 当前虚拟节点的负载大于平均负载*c，则选择下一个真实节点的虚拟节点
		conPtr := reflect.ValueOf(h.virtualToReal[h.virtualNodes[reqIdx]]).Pointer()
		if h.load[uint32(conPtr)] < float64(h.avgLoad)*h.c {
			fmt.Printf("%s -> server:%d real server:%v len:%d load:%f\n", key, reqIdx, h.virtualToReal[h.virtualNodes[reqIdx]], len(h.virtualNodes), h.load[uint32(conPtr)])
			return h.virtualToReal[h.virtualNodes[reqIdx]]
		}

		fmt.Printf("load:%f avg:%d\n", h.load[uint32(conPtr)], h.avgLoad)
		fmt.Printf("overload: %v\n", h.virtualToReal[h.virtualNodes[reqIdx]])

		// 获取下一个真实节点
		for h.virtualToReal[h.virtualNodes[reqIdx]] == realNode {
			reqIdx++
			reqIdx %= len(h.virtualNodes)
		}

		// 如果所有节点都过载
		// 节点全过载，返回请求哈希后的默认节点连接
		if reqIdx == idx {
			return h.virtualToReal[h.virtualNodes[idx]]
		}
	}
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

func NewHashRing(opts ...Opt) *HashRing {
	HR := &HashRing{
		virtualNodes:  make([]uint32, 0),
		virtualToReal: make(map[uint32]balancer.SubConn),
		virtualNode:   100,
		load:          make(map[uint32]float64),
		avgLoad:       math.MaxInt,
		c:             1.25,
	}
	for _, opt := range opts {
		opt(HR)
	}
	return HR
}

func WithVirtualNode(virtualNode int) Opt {
	return func(h *HashRing) {
		h.virtualNode = virtualNode
	}
}

func WithAvgLoad(avgLoad int) Opt {
	return func(h *HashRing) {
		h.avgLoad = avgLoad
	}
}

func WithC(c float64) Opt {
	return func(h *HashRing) {
		h.c = c
	}
}
