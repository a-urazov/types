package ring

import (
	"errors"
)

// Buffer represents a circular buffer (ring buffer) data structure.
type Buffer[T any] struct {
	data     []T
	size     int
	readPos  int
	writePos int
	count    int
}

// New creates and returns a new RingBuffer with the specified capacity.
func New[T any](capacity int) (*Buffer[T], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	return &Buffer[T]{
		data: make([]T, capacity),
		size: capacity,
	}, nil
}

// Put adds an element to the buffer. Returns true if the element was added, false if the buffer is full.
func (rb *Buffer[T]) Put(item T) bool {
	if rb.IsFull() {
		return false
	}
	rb.data[rb.writePos] = item
	rb.writePos = (rb.writePos + 1) % rb.size
	rb.count++
	return true
}

// Get removes and returns the oldest element from the buffer.
// It returns the element and a boolean indicating if the operation was successful.
func (rb *Buffer[T]) Get() (T, bool) {
	if rb.IsEmpty() {
		var zero T
		return zero, false
	}
	item := rb.data[rb.readPos]
	rb.readPos = (rb.readPos + 1) % rb.size
	rb.count--
	return item, true
}

// Peek returns the oldest element from the buffer without removing it.
// It returns the element and a boolean indicating if the operation was successful.
func (rb *Buffer[T]) Peek() (T, bool) {
	if rb.IsEmpty() {
		var zero T
		return zero, false
	}
	return rb.data[rb.readPos], true
}

// IsEmpty returns true if the buffer is empty, false otherwise.
func (rb *Buffer[T]) IsEmpty() bool {
	return rb.count == 0
}

// IsFull returns true if the buffer is full, false otherwise.
func (rb *Buffer[T]) IsFull() bool {
	return rb.count == rb.size
}

// Size returns the number of elements currently in the buffer.
func (rb *Buffer[T]) Size() int {
	return rb.count
}

// Capacity returns the maximum number of elements the buffer can hold.
func (rb *Buffer[T]) Capacity() int {
	return rb.size
}

// Clear removes all elements from the buffer.
func (rb *Buffer[T]) Clear() {
	rb.readPos = 0
	rb.writePos = 0
	rb.count = 0
}
