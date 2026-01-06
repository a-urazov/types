package segmenttree

import (
	"sync"
)

// SegmentTree represents a thread-safe segment tree data structure
type SegmentTree[T any] struct {
	data     []T
	tree     []T
	size     int
	mergeFn  func(a, b T) T
	mutex    sync.RWMutex
}

// New creates a new SegmentTree from the given data slice and merge function
// The merge function should be associative (e.g., sum, min, max)
func New[T any](data []T, mergeFn func(a, b T) T) *SegmentTree[T] {
	if len(data) == 0 {
		return &SegmentTree[T]{
			data:    []T{},
			tree:    []T{},
			size:    0,
			mergeFn: mergeFn,
		}
	}

	size := len(data)
	// Create a copy of the input data
	copiedData := make([]T, size)
	copy(copiedData, data)

	// Calculate tree size: next power of 2 >= size, then 2 * that - 1
	treeSize := 1
	for treeSize < size {
		treeSize <<= 1
	}
	treeSize = 2*treeSize - 1

	tree := make([]T, treeSize)

	st := &SegmentTree[T]{
		data:    copiedData,
		tree:    tree,
		size:    size,
		mergeFn: mergeFn,
	}

	// Build the tree
	st.build(0, 0, size-1)
	return st
}

// build recursively builds the segment tree
func (st *SegmentTree[T]) build(node, start, end int) {
	if start == end {
		st.tree[node] = st.data[start]
		return
	}

	mid := (start + end) / 2
	leftChild := 2*node + 1
	rightChild := 2*node + 2

	st.build(leftChild, start, mid)
	st.build(rightChild, mid+1, end)

	st.tree[node] = st.mergeFn(st.tree[leftChild], st.tree[rightChild])
}

// Query returns the result of applying the merge function to the range [left, right]
// Returns zero value and false if the range is invalid
func (st *SegmentTree[T]) Query(left, right int) (T, bool) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	if left < 0 || right >= st.size || left > right {
		var zero T
		return zero, false
	}

	return st.queryRange(0, 0, st.size-1, left, right)
}

// queryRange recursively queries the segment tree
func (st *SegmentTree[T]) queryRange(node, nodeStart, nodeEnd, queryStart, queryEnd int) (T, bool) {
	// Complete overlap
	if queryStart <= nodeStart && nodeEnd <= queryEnd {
		return st.tree[node], true
	}

	// No overlap
	if nodeEnd < queryStart || nodeStart > queryEnd {
		var zero T
		return zero, false
	}

	// Partial overlap
	mid := (nodeStart + nodeEnd) / 2
	leftChild := 2*node + 1
	rightChild := 2*node + 2

	leftResult, leftOk := st.queryRange(leftChild, nodeStart, mid, queryStart, queryEnd)
	rightResult, rightOk := st.queryRange(rightChild, mid+1, nodeEnd, queryStart, queryEnd)

	if leftOk && rightOk {
		return st.mergeFn(leftResult, rightResult), true
	}
	if leftOk {
		return leftResult, true
	}
	if rightOk {
		return rightResult, true
	}

	var zero T
	return zero, false
}

// Update updates the value at the given index
// Returns true if the update was successful, false if the index is invalid
func (st *SegmentTree[T]) Update(index int, value T) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if index < 0 || index >= st.size {
		return false
	}

	st.data[index] = value
	st.updateRange(0, 0, st.size-1, index, value)
	return true
}

// updateRange recursively updates the segment tree
func (st *SegmentTree[T]) updateRange(node, start, end, index int, value T) {
	if start == end {
		st.tree[node] = value
		return
	}

	mid := (start + end) / 2
	leftChild := 2*node + 1
	rightChild := 2*node + 2

	if index <= mid {
		st.updateRange(leftChild, start, mid, index, value)
	} else {
		st.updateRange(rightChild, mid+1, end, index, value)
	}

	st.tree[node] = st.mergeFn(st.tree[leftChild], st.tree[rightChild])
}

// Size returns the number of elements in the segment tree
func (st *SegmentTree[T]) Size() int {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	return st.size
}

// IsEmpty returns true if the segment tree is empty
func (st *SegmentTree[T]) IsEmpty() bool {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	return st.size == 0
}

// Get returns the value at the given index
// Returns zero value and false if the index is invalid
func (st *SegmentTree[T]) Get(index int) (T, bool) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	if index < 0 || index >= st.size {
		var zero T
		return zero, false
	}

	return st.data[index], true
}

// ToSlice returns a copy of the underlying data array
func (st *SegmentTree[T]) ToSlice() []T {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	result := make([]T, st.size)
	copy(result, st.data)
	return result
}

// ForEach iterates over all elements in the underlying data array
func (st *SegmentTree[T]) ForEach(fn func(index int, value T)) {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	for i, value := range st.data {
		fn(i, value)
	}
}

// Clear removes all elements from the segment tree
func (st *SegmentTree[T]) Clear() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.data = []T{}
	st.tree = []T{}
	st.size = 0
}

// Resize changes the size of the segment tree
// New elements are initialized with zero values
// Existing elements are preserved up to the new size
func (st *SegmentTree[T]) Resize(newSize int) {
	if newSize < 0 {
		newSize = 0
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	if newSize == st.size {
		return
	}

	oldData := st.data
	newData := make([]T, newSize)

	// Copy existing data
	copyLen := st.size
	if copyLen > newSize {
		copyLen = newSize
	}
	copy(newData, oldData[:copyLen])

	st.data = newData
	st.size = newSize

	if newSize == 0 {
		st.tree = []T{}
	} else {
		// Rebuild the tree
		treeSize := 1
		for treeSize < newSize {
			treeSize <<= 1
		}
		treeSize = 2*treeSize - 1
		st.tree = make([]T, treeSize)
		st.build(0, 0, newSize-1)
	}
}

// Sum creates a segment tree optimized for sum queries
func Sum(data []int) *SegmentTree[int] {
	return New(data, func(a, b int) int { return a + b })
}

// Min creates a segment tree optimized for minimum queries
func Min(data []int) *SegmentTree[int] {
	return New(data, func(a, b int) int {
		if a < b {
			return a
		}
		return b
	})
}

// Max creates a segment tree optimized for maximum queries
func Max(data []int) *SegmentTree[int] {
	return New(data, func(a, b int) int {
		if a > b {
			return a
		}
		return b
	})
}