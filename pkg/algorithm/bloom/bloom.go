package bloom

import (
	"hash"
	"math"
	"sync"
)

type bloomFilter struct {
	bitset   []byte
	n        uint
	m        uint
	k        uint32
	hashFunc hash.Hash32

	mutex sync.Mutex
}

// n = ceil(m / (-k / log(1 - exp(log(p) / k))))		// 元素数量
// p = pow(1 - exp(-k / (m / n)), k)					// 假阳概率
// m = ceil((n * log(p)) / log(1 / pow(2, log(2))));	// bitmap 位数
// k = round((m / n) * log(2));							// 哈希函数数量

func NewBloomFilterByNP(n, p float64, hashFunc hash.Hash32) Filter {
	m := math.Ceil((n * math.Log(p)) / math.Log(1/math.Pow(2, math.Log(2))))
	k := math.Round((m / n) * 0.69)

	if k < 1 {
		k = 1
	} else if k > 30 {
		k = 30
	}

	// 防止 bitmap 太小
	if m < 64 {
		m = 64
	}

	// 控制 bitmap 的位数为 2 的 n 次幂
	bytes := (m + 7) / 8

	return &bloomFilter{
		bitset:   make([]byte, int(bytes)),
		k:        uint32(k),
		m:        uint(bytes) * 8,
		hashFunc: hashFunc,
	}
}

func (bf *bloomFilter) Set(key []byte) Filter {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// double-hashing todo 详见 blog
	h := bf.hash(key)
	delta := (h >> 17) | (h << 15)

	for i := uint32(0); i < bf.k; i++ {
		bitPos := h % uint32(bf.m)
		bf.bitset[bitPos/8] |= 1 << (bitPos % 8)
		h += delta
	}

	return bf
}

func (bf *bloomFilter) Test(key []byte) bool {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	if len(bf.bitset) < 2 {
		return false
	}

	bits := len(bf.bitset) * 8
	h := bf.hash(key)
	delta := (h >> 17) | (h << 15)

	for i := uint32(0); i < bf.k; i++ {
		bitPos := h % uint32(bits)
		if (bf.bitset[bitPos/8] & (1 << (bitPos % 8))) == 0 {
			return false
		}
		h += delta
	}
	return true
}

func (bf *bloomFilter) hash(key []byte) uint32 {
	bf.hashFunc.Reset()
	bf.hashFunc.Write(key)
	return bf.hashFunc.Sum32()
}
