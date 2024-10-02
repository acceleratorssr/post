package utils

import (
	"fmt"
	"testing"
)

func TestPQ(t *testing.T) {
	h := NewMinHeap()
	h.Insert(1, 3)
	h.Insert(2, 1)
	h.Insert(3, 7)
	h.Insert(4, -10)
	h.Insert(5, 6)
	h.Insert(6, 5)
	h.Insert(7, 2)
	h.Insert(8, 4)

	fmt.Println(h.heap)
	fmt.Println(h.ExtractMin())
	fmt.Println(h.heap)
}
