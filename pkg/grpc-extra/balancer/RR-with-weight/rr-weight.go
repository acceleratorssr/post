package rrw

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	grpc_extra "post/pkg/grpc-extra"
	"sync"
)

const (
	RoundRobinWithWeight = "roundrobin-with-weight"
	weightLimit          = 100 // 防负载超不均衡
)

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(RoundRobinWithWeight,
		&WeightedPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type WeightedPicker struct {
	mutex sync.Mutex
	conns []*weightConn
}

// Pick
// 异常响应，降低权重
// 正常响应，提高权重
// 存在阈值限制
// 平滑算法
func (b *WeightedPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var totalWeight int
	maxWeight, maxi := 0, 0

	b.mutex.Lock()
	defer b.mutex.Unlock()

	maxConn := b.conns[0]
	for i, node := range b.conns { // conns 至少有一个 subConn
		totalWeight += node.weight
		node.currentWeight += node.weight

		if maxWeight < node.currentWeight {
			maxConn = node
			maxWeight = node.currentWeight
			maxi = i
		}
	}
	b.conns[maxi].currentWeight -= totalWeight

	return balancer.PickResult{
		SubConn: maxConn,
		Done: func(info balancer.DoneInfo) {
			b.mutex.Lock()
			defer b.mutex.Unlock()

			if info.Err != nil {
				maxConn.weight = max(1, maxConn.weight-1)
			} else {
				maxConn.weight = min(maxConn.weight+1, weightLimit)
			}
		},
	}, nil
}

type WeightedPickerBuilder struct {
}

func (b *WeightedPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for con, conInfo := range info.ReadySCs {
		weight := 0.0
		if conInfo.Address.Attributes.Value(grpc_extra.Weight) != nil {
			weight = conInfo.Address.Attributes.Value(grpc_extra.Weight).(float64)
		}

		conns = append(conns, &weightConn{
			SubConn:       con,
			weight:        int(weight),
			currentWeight: int(weight),
		})
	}
	return &WeightedPicker{
		conns: conns,
	}
}

type weightConn struct {
	weight        int
	currentWeight int
	balancer.SubConn
}
