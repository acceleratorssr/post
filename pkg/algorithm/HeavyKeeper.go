package algorithm

import (
	"fmt"
	"math/rand"
)

const (
	numHashBuckets = 3    // 哈希桶数
	bucketSize     = 100  // 每个哈希桶的大小
	decayRate      = 0.95 // 衰减系数
)

type HeavyKeeper struct {
	hashBuckets [][]element // 哈希桶
}

type element struct {
	key   string
	count float64
}

func NewHeavyKeeper() *HeavyKeeper {
	hashBuckets := make([][]element, numHashBuckets)
	for i := range hashBuckets {
		hashBuckets[i] = make([]element, bucketSize)
	}
	return &HeavyKeeper{hashBuckets: hashBuckets}
}

func (hk *HeavyKeeper) Insert(key string) {
	for i := 0; i < numHashBuckets; i++ {
		index := rand.Intn(bucketSize) // 通过哈希函数随机找到位置
		bucket := hk.hashBuckets[i]
		if bucket[index].key == key {
			// 如果 key 存在，增加频率
			bucket[index].count++
		} else if bucket[index].count == 0 {
			// 如果为空，直接插入
			bucket[index] = element{key: key, count: 1}
		} else {
			// 否则按衰减概率替换
			if rand.Float64() < 1/(bucket[index].count) {
				bucket[index] = element{key: key, count: 1}
			}
		}
	}
}

func (hk *HeavyKeeper) GetHotElements() []string {
	elementFreq := make(map[string]float64)

	// 统计元素频率
	for i := 0; i < numHashBuckets; i++ {
		for _, e := range hk.hashBuckets[i] {
			if e.count > 0 {
				elementFreq[e.key] += e.count
			}
		}
	}

	fmt.Println("Element frequencies:", elementFreq)

	hotElements := []string{}
	for k, v := range elementFreq {
		if v > 100 { // 阈值
			hotElements = append(hotElements, k)
		}
	}
	return hotElements
}

func (hk *HeavyKeeper) Decay() {
	for i := 0; i < numHashBuckets; i++ {
		for j := 0; j < bucketSize; j++ {
			hk.hashBuckets[i][j].count *= decayRate
		}
	}
}
