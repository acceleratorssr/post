package bloom

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type Filter interface {
	Set(key []byte) Filter
	Test(data []byte) bool
}

type bloomThirdWrapper struct {
	b *bloom.BloomFilter
}

func (b *bloomThirdWrapper) Set(key []byte) Filter {
	bf := b.b.Add(key)
	return &bloomThirdWrapper{bf}
}

func (b *bloomThirdWrapper) Test(data []byte) bool {
	return b.b.Test(data)
}

func NewBloomFilterThirdByNP(n uint, p float64) Filter {
	bf := bloom.NewWithEstimates(n, p)
	return &bloomThirdWrapper{bf}
}
