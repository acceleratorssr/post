package consistent_hashing

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(ConsistentHash,
		&ConsistentHashPickerBuilder{},
		base.Config{HealthCheck: true},
	)
}

func init() {
	balancer.Register(newBuilder())
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
