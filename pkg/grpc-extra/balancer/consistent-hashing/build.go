package consistent_hashing

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"reflect"
)

// balancer 注册阶段调用
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(ConsistentHash,
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
		b.HashRing = NewHashRing(100)
		for con := range info.ReadySCs {
			conPtr := reflect.ValueOf(con).Pointer()
			b.HashRing.AddNode(con, int(conPtr))
		}
	} else {
		currentNodes := make(map[balancer.SubConn]struct{})
		for subConn := range info.ReadySCs {
			currentNodes[subConn] = struct{}{}
		}

		for oldOne := range b.subConn {
			if _, exists := currentNodes[oldOne]; !exists {
				conPtr := reflect.ValueOf(oldOne).Pointer()
				b.HashRing.RemoveNode(int(conPtr))
				delete(b.subConn, oldOne)
				fmt.Println("remove node:", oldOne)
			}
		}

		for newOne := range currentNodes {
			if _, exists := b.subConn[newOne]; !exists {
				conPtr := reflect.ValueOf(newOne).Pointer()
				b.HashRing.AddNode(newOne, int(conPtr))
				b.subConn[newOne] = struct{}{}
				fmt.Println("add node:", newOne)
			}
		}
	}

	fmt.Printf("picker build :%d nodes.\n", len(b.HashRing.hashNodes))
	return &ConsistentHashPicker{
		hashRing: b.HashRing,
	}
}
