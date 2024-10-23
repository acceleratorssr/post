package bloom

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"hash"
	"hash/crc32"
	"hash/fnv"
	mr "math/rand"
	"testing"
	"time"

	"github.com/spaolacci/murmur3"
)

// 	i := uint32(100)
//  n1 := make([]byte, 4)
//  binary.BigEndian.PutUint32(n1, i)
//  filter.Add(n1)

func TestBloomThird(t *testing.T) {
	s := mr.New(mr.NewSource(time.Now().UnixNano()))

	keySet := make(map[string]struct{})
	var keys []string

	for len(keys) < 1000 {
		key, err := GenerateRandomString(10, s)
		if err != nil {
			t.Fatalf("Failed to generate random string: %v", err)
		}
		if _, exists := keySet[key]; !exists {
			keySet[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	testCases := []struct {
		name     string
		hashFunc func() hash.Hash32
	}{
		{"FNV-1a", fnv.New32},
		{"CRC32", crc32.NewIEEE},
		{"Murmur3", func() hash.Hash32 { return murmur3.New32() }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := bloom.NewWithEstimates(uint(len(keys)), 1e-5)

			for _, key := range keys {
				filter.AddString(key)
			}

			for _, key := range keys {
				if !filter.Test([]byte(key)) {
					t.Errorf("Expected true for key '%s', got false", key)
				}
			}

			falsePositiveCount := 0
			tests := 10000
			for i := 0; i < tests; i++ {
				randomKey := fmt.Sprintf("key%d", mr.Intn(20000)+10000)
				if filter.Test([]byte(randomKey)) {
					falsePositiveCount++
				}
			}

			falsePositiveRate := float64(falsePositiveCount) / float64(tests)
			t.Logf("Hash Function: %s, False Positive Rate: %f%%", tc.name, falsePositiveRate*100)
		})
	}
}
