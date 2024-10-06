package utils

import "math"

type Like struct {
	ID    uint64
	Score float64
}

type Option func(*MinHeap)

func WithLimit(limit int) Option {
	return func(h *MinHeap) {
		h.limit = limit
	}
}

type MinHeap struct {
	heap  []Like
	limit int
}

func NewMinHeap(opts ...Option) *MinHeap {
	mh := &MinHeap{
		heap:  []Like{},
		limit: math.MaxInt,
	}

	for _, opt := range opts {
		opt(mh)
	}

	return mh
}

func (h *MinHeap) Insert(id uint64, val float64) {
	if len(h.heap) < h.limit {
		h.heap = append(h.heap, Like{ID: id, Score: val})
		h.heapifyUp(len(h.heap) - 1)
	} else if val > h.heap[0].Score {
		h.heap[0] = Like{ID: id, Score: val}
		h.heapifyDown(0)
	}
}

// ExtractMin 删除并返回堆中的最小元素
func (h *MinHeap) ExtractMin() *Like {
	if len(h.heap) == 0 {
		return nil
	}
	minn := h.heap[0]
	h.heap[0] = h.heap[len(h.heap)-1]
	h.heap = h.heap[:len(h.heap)-1]
	h.heapifyDown(0)
	return &minn
}

func (h *MinHeap) GetMin() *Like {
	if len(h.heap) == 0 {
		return nil
	}
	return &h.heap[0]
}

func (h *MinHeap) GetLen() int {
	return len(h.heap)
}

// heapifyUp 维护堆的性质从下往上
func (h *MinHeap) heapifyUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if h.heap[parent].Score > h.heap[index].Score {
			h.heap[parent], h.heap[index] = h.heap[index], h.heap[parent]
			index = parent
		} else {
			break
		}
	}
}

// heapifyDown 维护堆的性质从上往下
func (h *MinHeap) heapifyDown(index int) {
	lastIndex := len(h.heap) - 1
	for {
		leftChild := 2*index + 1
		rightChild := 2*index + 2
		smallest := index

		if leftChild <= lastIndex && h.heap[leftChild].Score < h.heap[smallest].Score {
			smallest = leftChild
		}
		if rightChild <= lastIndex && h.heap[rightChild].Score < h.heap[smallest].Score {
			smallest = rightChild
		}
		if smallest != index {
			h.heap[index], h.heap[smallest] = h.heap[smallest], h.heap[index]
			index = smallest
		} else {
			break
		}
	}
}
