package CHBL

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"reflect"
)

// balancer 注册阶段调用
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(CHBL,
		&ConsistentHashPickerBuilder{
			subConn: make(map[balancer.SubConn]struct{}),
		},
		base.Config{HealthCheck: true},
	)
}

func init() {
	balancer.Register(newBuilder())
}

type ConsistentHashPickerBuilder struct {
	HashRing *HashRing
	subConn  map[balancer.SubConn]struct{}
}

func (b *ConsistentHashPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if b.HashRing == nil {
		b.HashRing = NewHashRing(WithC(2)) // 不是负载敏感，所以平衡参数选高点，减少数据迁移量
		for con := range info.ReadySCs {
			conPtr := reflect.ValueOf(con).Pointer()
			b.HashRing.AddNode(con, uint32(conPtr))
			b.subConn[con] = struct{}{}
		}
	} else {
		avgLoad := 0
		currentNodes := make(map[balancer.SubConn]struct{})
		for subConn := range info.ReadySCs {
			currentNodes[subConn] = struct{}{}
		}

		for oldOne := range b.subConn {
			conPtr := reflect.ValueOf(oldOne).Pointer()
			if _, exists := currentNodes[oldOne]; !exists {
				b.HashRing.RemoveNode(uint32(conPtr))
				delete(b.subConn, oldOne)
				delete(currentNodes, oldOne)
				fmt.Println("remove node:", oldOne)
			} else {
				cnt := info.ReadySCs[oldOne].Address.Attributes.Value("request_count").(int) // 真实节点为负载获取的粒度
				avgLoad += cnt
				b.HashRing.load[uint32(conPtr)] = float64(cnt)
			}
		}

		for newOne := range currentNodes {
			if _, exists := b.subConn[newOne]; !exists {
				conPtr := reflect.ValueOf(newOne).Pointer()
				b.HashRing.AddNode(newOne, uint32(conPtr))

				b.subConn[newOne] = struct{}{}
				fmt.Println("add node:", newOne)

				cnt := info.ReadySCs[newOne].Address.Attributes.Value("request_count").(int)
				avgLoad += cnt
				b.HashRing.load[uint32(conPtr)] = float64(cnt)
				fmt.Printf("server: %s, load:%d \n", info.ReadySCs[newOne].Address.Addr, cnt)
			}
		}

		b.HashRing.avgLoad = max(avgLoad/len(b.subConn), 1000)
	}

	fmt.Printf("picker build :%d nodes.\n", len(b.HashRing.virtualNodes))
	return &ConsistentHashPicker{
		hashRing: b.HashRing,
	}
}
