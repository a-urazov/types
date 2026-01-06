package circularlist

import (
	"sync"
)

// Node represents a node in the circular doubly linked list
type Node[T any] struct {
	Value T
	next  *Node[T]
	prev  *Node[T]
}

// CircularList represents a thread-safe circular doubly linked list
type CircularList[T any] struct {
	head  *Node[T]
	tail  *Node[T]
	size  int
	mutex sync.RWMutex
}

// New creates a new empty CircularList
func New[T any]() *CircularList[T] {
	return &CircularList[T]{}
}

// Add adds an element to the end of the list
func (cl *CircularList[T]) Add(value T) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	newNode := &Node[T]{Value: value}

	if cl.size == 0 {
		// First element
		newNode.next = newNode
		newNode.prev = newNode
		cl.head = newNode
		cl.tail = newNode
	} else {
		// Add to the end
		newNode.prev = cl.tail
		newNode.next = cl.head
		cl.head.prev = newNode
		cl.tail.next = newNode
		cl.tail = newNode
	}

	cl.size++
}

// AddAt adds an element at the specified index
func (cl *CircularList[T]) AddAt(index int, value T) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if index < 0 || index > cl.size {
		return false
	}

	if index == cl.size {
		// Add to the end
		cl.addAtEnd(value)
		return true
	}

	if index == 0 {
		// Add to the beginning
		cl.addAtBeginning(value)
		return true
	}

	// Add in the middle
	current := cl.getNodeAt(index)
	if current == nil {
		return false
	}

	newNode := &Node[T]{Value: value}
	newNode.prev = current.prev
	newNode.next = current
	current.prev.next = newNode
	current.prev = newNode

	cl.size++
	return true
}

// addAtEnd adds an element to the end of the list (internal method)
func (cl *CircularList[T]) addAtEnd(value T) {
	newNode := &Node[T]{Value: value}

	if cl.size == 0 {
		// First element
		newNode.next = newNode
		newNode.prev = newNode
		cl.head = newNode
		cl.tail = newNode
	} else {
		// Add to the end
		newNode.prev = cl.tail
		newNode.next = cl.head
		cl.head.prev = newNode
		cl.tail.next = newNode
		cl.tail = newNode
	}

	cl.size++
}

// addAtBeginning adds an element to the beginning of the list (internal method)
func (cl *CircularList[T]) addAtBeginning(value T) {
	newNode := &Node[T]{Value: value}

	if cl.size == 0 {
		// First element
		newNode.next = newNode
		newNode.prev = newNode
		cl.head = newNode
		cl.tail = newNode
	} else {
		// Add to the beginning
		newNode.next = cl.head
		newNode.prev = cl.tail
		cl.head.prev = newNode
		cl.tail.next = newNode
		cl.head = newNode
	}

	cl.size++
}

// Remove removes the first occurrence of the specified element
func (cl *CircularList[T]) Remove(value T) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if cl.size == 0 {
		return false
	}

	current := cl.head
	for i := 0; i < cl.size; i++ {
		if any(current.Value) == any(value) {
			return cl.removeNode(current)
		}
		current = current.next
	}

	return false
}

// RemoveAt removes the element at the specified index
func (cl *CircularList[T]) RemoveAt(index int) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if index < 0 || index >= cl.size {
		return false
	}

	nodeToRemove := cl.getNodeAt(index)
	if nodeToRemove == nil {
		return false
	}

	return cl.removeNode(nodeToRemove)
}

// removeNode removes the specified node (internal method)
func (cl *CircularList[T]) removeNode(node *Node[T]) bool {
	if cl.size == 0 {
		return false
	}

	if cl.size == 1 {
		// Only one element
		cl.head = nil
		cl.tail = nil
	} else {
		// Update neighbors
		node.prev.next = node.next
		node.next.prev = node.prev

		// Update head/tail if necessary
		if node == cl.head {
			cl.head = node.next
		}
		if node == cl.tail {
			cl.tail = node.prev
		}
	}

	cl.size--
	return true
}

// Get returns the element at the specified index
func (cl *CircularList[T]) Get(index int) (T, bool) {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	if index < 0 || index >= cl.size {
		var zero T
		return zero, false
	}

	node := cl.getNodeAt(index)
	if node == nil {
		var zero T
		return zero, false
	}

	return node.Value, true
}

// Set sets the element at the specified index
func (cl *CircularList[T]) Set(index int, value T) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if index < 0 || index >= cl.size {
		return false
	}

	node := cl.getNodeAt(index)
	if node == nil {
		return false
	}

	node.Value = value
	return true
}

// getNodeAt returns the node at the specified index (internal method)
func (cl *CircularList[T]) getNodeAt(index int) *Node[T] {
	if index < 0 || index >= cl.size {
		return nil
	}

	// Optimize for common cases
	if index == 0 {
		return cl.head
	}
	if index == cl.size-1 {
		return cl.tail
	}

	// Choose the shorter path (from head or tail)
	var current *Node[T]
	if index < cl.size/2 {
		// Go from head
		current = cl.head
		for i := 0; i < index; i++ {
			current = current.next
		}
	} else {
		// Go from tail
		current = cl.tail
		for i := cl.size - 1; i > index; i-- {
			current = current.prev
		}
	}

	return current
}

// Size returns the number of elements in the list
func (cl *CircularList[T]) Size() int {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	return cl.size
}

// IsEmpty returns true if the list is empty
func (cl *CircularList[T]) IsEmpty() bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	return cl.size == 0
}

// Clear removes all elements from the list
func (cl *CircularList[T]) Clear() {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.head = nil
	cl.tail = nil
	cl.size = 0
}

// Contains checks if the specified element exists in the list
func (cl *CircularList[T]) Contains(value T) bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.head
	for i := 0; i < cl.size; i++ {
		if any(current.Value) == any(value) {
			return true
		}
		current = current.next
	}

	return false
}

// IndexOf returns the index of the first occurrence of the specified element
func (cl *CircularList[T]) IndexOf(value T) int {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.head
	for i := 0; i < cl.size; i++ {
		if any(current.Value) == any(value) {
			return i
		}
		current = current.next
	}

	return -1
}

// ToSlice returns a slice containing all elements in the list
func (cl *CircularList[T]) ToSlice() []T {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	if cl.size == 0 {
		return []T{}
	}

	result := make([]T, cl.size)
	current := cl.head
	for i := 0; i < cl.size; i++ {
		result[i] = current.Value
		current = current.next
	}

	return result
}

// ForEach iterates over all elements in the list
func (cl *CircularList[T]) ForEach(fn func(index int, value T)) {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.head
	for i := 0; i < cl.size; i++ {
		fn(i, current.Value)
		current = current.next
	}
}

// ReverseForEach iterates over all elements in the list in reverse order
func (cl *CircularList[T]) ReverseForEach(fn func(index int, value T)) {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.tail
	for i := cl.size - 1; i >= 0; i-- {
		fn(i, current.Value)
		current = current.prev
	}
}

// GetNext returns the next element after the specified value in the circular list
func (cl *CircularList[T]) GetNext(value T) (T, bool) {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.head
	for i := 0; i < cl.size; i++ {
		if any(current.Value) == any(value) {
			next := current.next
			return next.Value, true
		}
		current = current.next
	}

	var zero T
	return zero, false
}

// GetPrev returns the previous element before the specified value in the circular list
func (cl *CircularList[T]) GetPrev(value T) (T, bool) {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()

	current := cl.head
	for i := 0; i < cl.size; i++ {
		if any(current.Value) == any(value) {
			prev := current.prev
			return prev.Value, true
		}
		current = current.next
	}

	var zero T
	return zero, false
}

// RotateLeft rotates the list to the left by n positions
func (cl *CircularList[T]) RotateLeft(n int) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if cl.size <= 1 || n == 0 {
		return
	}

	// Normalize rotation amount
	n = n % cl.size
	if n < 0 {
		n = cl.size + n
	}

	// Find the new head by moving right (since moving left means going to the next element)
	current := cl.head
	for i := 0; i < n; i++ {
		current = current.next
	}

	// Update head and tail
	cl.head = current
	cl.tail = current.prev
}

// RotateRight rotates the list to the right by n positions
func (cl *CircularList[T]) RotateRight(n int) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	if cl.size <= 1 || n == 0 {
		return
	}

	// Normalize rotation amount
	n = n % cl.size
	if n < 0 {
		n = cl.size + n
	}

	// Find the new head by moving left (since moving right means going to the previous element)
	current := cl.head
	for i := 0; i < cl.size-n; i++ { // Move left by n positions means moving right by (size-n)
		current = current.prev
	}

	// Update head and tail
	cl.head = current
	cl.tail = current.prev
}
