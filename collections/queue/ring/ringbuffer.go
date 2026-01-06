package ring

import (
	"errors"
	"types/collections/list"
)

// Buffer represents a circular buffer (ring buffer) data structure.
type Buffer[T any] struct {
	items    *list.List[T]
	capacity int
}

// New creates and returns a new RingBuffer with the specified capacity.
func New[T any](capacity int) (*Buffer[T], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	return &Buffer[T]{
		items:    list.New[T](),
		capacity: capacity,
	}, nil
}

// Put adds an element to the buffer. Returns true if the element was added, false if the buffer is full.
func (rb *Buffer[T]) Put(item T) bool {
	if rb.IsFull() {
		return false
	}
	rb.items.Add(item)
	return true
}

// Get removes and returns the oldest element from the buffer.
// It returns the element and a boolean indicating if the operation was successful.
func (rb *Buffer[T]) Get() (T, bool) {
	if rb.IsEmpty() {
		var zero T
		return zero, false
	}
	item, _ := rb.items.Get(0)
	rb.items.RemoveAt(0)
	return item, true
}

// Peek returns the oldest element from the buffer without removing it.
// It returns the element and a boolean indicating if the operation was successful.
func (rb *Buffer[T]) Peek() (T, bool) {
	if rb.IsEmpty() {
		var zero T
		return zero, false
	}
	return rb.items.Get(0)
}

// IsEmpty returns true if the buffer is empty, false otherwise.
func (rb *Buffer[T]) IsEmpty() bool {
	return rb.items.IsEmpty()
}

// IsFull returns true if the buffer is full, false otherwise.
func (rb *Buffer[T]) IsFull() bool {
	return rb.items.Size() == rb.capacity
}

// Size returns the number of elements currently in the buffer.
func (rb *Buffer[T]) Size() int {
	return rb.items.Size()
}

// Capacity returns the maximum number of elements the buffer can hold.
func (rb *Buffer[T]) Capacity() int {
	return rb.capacity
}

// Clear removes all elements from the buffer.
func (rb *Buffer[T]) Clear() {
	rb.items.Clear()
}
