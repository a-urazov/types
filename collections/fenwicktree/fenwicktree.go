package fenwicktree

import (
	"sync"
)

// FenwickTree represents a thread-safe Fenwick tree (Binary Indexed Tree)
type FenwickTree struct {
	tree  []int
	size  int
	mutex sync.RWMutex
}

// New creates a new FenwickTree of the given size
func New(size int) *FenwickTree {
	if size < 0 {
		size = 0
	}
	return &FenwickTree{
		tree: make([]int, size+1),
		size: size,
	}
}

// FromSlice creates a new FenwickTree from the given data slice
func FromSlice(data []int) *FenwickTree {
	size := len(data)
	ft := New(size)
	for i, val := range data {
		ft.Add(i, val)
	}
	return ft
}

// Add adds the given value to the element at the specified index
func (ft *FenwickTree) Add(index, value int) {
	if index < 0 || index >= ft.size {
		return
	}

	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	// Fenwick tree uses 1-based indexing internally
	index++
	for index <= ft.size {
		ft.tree[index] += value
		// Move to the next responsible node
		index += index & -index
	}
}

// Query returns the sum of elements from the beginning up to the specified index (inclusive)
func (ft *FenwickTree) Query(index int) int {
	if index < 0 {
		return 0
	}
	if index >= ft.size {
		index = ft.size - 1
	}

	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	sum := 0
	// Fenwick tree uses 1-based indexing internally
	index++
	for index > 0 {
		sum += ft.tree[index]
		// Move to the parent node
		index -= index & -index
	}
	return sum
}

// QueryRange returns the sum of elements in the range [left, right] (inclusive)
func (ft *FenwickTree) QueryRange(left, right int) int {
	if left > right {
		return 0
	}
	if left < 0 {
		left = 0
	}
	if right >= ft.size {
		right = ft.size - 1
	}

	// Sum(left, right) = Sum(0, right) - Sum(0, left-1)
	sumRight := ft.Query(right)
	sumLeft := ft.Query(left - 1)
	return sumRight - sumLeft
}

// Get returns the value of the element at the specified index
// This is O(log n) as it's calculated from prefix sums
func (ft *FenwickTree) Get(index int) int {
	if index < 0 || index >= ft.size {
		return 0
	}
	return ft.QueryRange(index, index)
}

// Set sets the value of the element at the specified index
// This replaces the existing value
func (ft *FenwickTree) Set(index, value int) {
	if index < 0 || index >= ft.size {
		return
	}

	// Get current value
	currentValue := ft.Get(index)
	// Calculate difference and add it
	diff := value - currentValue
	ft.Add(index, diff)
}

// Size returns the number of elements in the Fenwick tree
func (ft *FenwickTree) Size() int {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	return ft.size
}

// IsEmpty returns true if the Fenwick tree is empty (all elements are zero)
func (ft *FenwickTree) IsEmpty() bool {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	for _, val := range ft.tree {
		if val != 0 {
			return false
		}
	}
	return true
}

// Clear resets the Fenwick tree to all zeros
func (ft *FenwickTree) Clear() {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	for i := range ft.tree {
		ft.tree[i] = 0
	}
}

// ToSlice returns a slice representing the current state of the data
// This is an O(n log n) operation
func (ft *FenwickTree) ToSlice() []int {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	data := make([]int, ft.size)
	for i := 0; i < ft.size; i++ {
		data[i] = ft.QueryRange(i, i)
	}
	return data
}

// ForEach iterates over all elements in the Fenwick tree
// This is an O(n log n) operation
func (ft *FenwickTree) ForEach(fn func(index, value int)) {
	slice := ft.ToSlice()
	for i, val := range slice {
		fn(i, val)
	}
}

// Resize changes the size of the Fenwick tree
// New elements are initialized with zero values
// Existing elements are preserved up to the new size
func (ft *FenwickTree) Resize(newSize int) {
	if newSize < 0 {
		newSize = 0
	}

	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	if newSize == ft.size {
		return
	}

	// Get old data without causing deadlock
	oldData := make([]int, ft.size)
	for i := 0; i < ft.size; i++ {
		oldData[i] = ft.queryRangeInternal(i, i)
	}

	newData := make([]int, newSize)
	copyLen := ft.size
	if copyLen > newSize {
		copyLen = newSize
	}
	copy(newData, oldData[:copyLen])

	ft.tree = make([]int, newSize+1)
	ft.size = newSize
	for i, val := range newData {
		ft.addInternal(i, val)
	}
}

// queryRangeInternal is the non-locking version of QueryRange
func (ft *FenwickTree) queryRangeInternal(left, right int) int {
	if left > right {
		return 0
	}
	if left < 0 {
		left = 0
	}
	if right >= ft.size {
		right = ft.size - 1
	}

	sumRight := ft.queryInternal(right)
	sumLeft := ft.queryInternal(left - 1)
	return sumRight - sumLeft
}

// queryInternal is the non-locking version of Query
func (ft *FenwickTree) queryInternal(index int) int {
	if index < 0 {
		return 0
	}
	if index >= ft.size {
		index = ft.size - 1
	}

	sum := 0
	index++
	for index > 0 {
		sum += ft.tree[index]
		index -= index & -index
	}
	return sum
}

// addInternal is the non-locking version of Add
func (ft *FenwickTree) addInternal(index, value int) {
	index++
	for index <= ft.size {
		ft.tree[index] += value
		index += index & -index
	}
}
