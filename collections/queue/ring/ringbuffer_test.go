package ring

import (
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	rb, err := New[int](5)
	if err != nil {
		t.Errorf("New() should not return an error, got %v", err)
	}
	if rb == nil {
		t.Error("New() should not return nil")
	}
	if rb.Capacity() != 5 {
		t.Errorf("New ring buffer capacity should be 5, got %d", rb.Capacity())
	}
	if rb.Size() != 0 {
		t.Errorf("New ring buffer size should be 0, got %d", rb.Size())
	}
	if !rb.IsEmpty() {
		t.Error("New ring buffer should be empty")
	}
	if rb.IsFull() {
		t.Error("New ring buffer should not be full")
	}

	// Test error case
	_, err = New[int](0)
	if err == nil {
		t.Error("New() should return an error for capacity <= 0")
	}
}

func TestPut(t *testing.T) {
	rb, _ := New[int](3)

	if !rb.Put(1) {
		t.Error("Put() should return true when buffer is not full")
	}
	if rb.Size() != 1 {
		t.Errorf("Ring buffer size should be 1, got %d", rb.Size())
	}
	if rb.IsEmpty() {
		t.Error("Ring buffer should not be empty after Put")
	}

	if !rb.Put(2) {
		t.Error("Put() should return true when buffer is not full")
	}
	if rb.Size() != 2 {
		t.Errorf("Ring buffer size should be 2, got %d", rb.Size())
	}

	if !rb.Put(3) {
		t.Error("Put() should return true when buffer is not full")
	}
	if rb.Size() != 3 {
		t.Errorf("Ring buffer size should be 3, got %d", rb.Size())
	}
	if !rb.IsFull() {
		t.Error("Ring buffer should be full after putting 3 elements in capacity 3")
	}

	if rb.Put(4) {
		t.Error("Put() should return false when buffer is full")
	}
	if rb.Size() != 3 {
		t.Errorf("Ring buffer size should still be 3, got %d", rb.Size())
	}
}

func TestGet(t *testing.T) {
	rb, _ := New[int](3)

	// Test Get on empty buffer
	_, ok := rb.Get()
	if ok {
		t.Error("Get() on empty buffer should return false")
	}

	rb.Put(1)
	rb.Put(2)
	rb.Put(3)

	val, ok := rb.Get()
	if !ok || val != 1 {
		t.Errorf("Get() should return 1, got %v", val)
	}
	if rb.Size() != 2 {
		t.Errorf("Ring buffer size should be 2, got %d", rb.Size())
	}

	val, ok = rb.Get()
	if !ok || val != 2 {
		t.Errorf("Get() should return 2, got %v", val)
	}
	if rb.Size() != 1 {
		t.Errorf("Ring buffer size should be 1, got %d", rb.Size())
	}

	// Add more items to test wrap-around behavior
	rb.Put(4)
	rb.Put(5)
	if rb.Size() != 3 {
		t.Errorf("Ring buffer size should be 3, got %d", rb.Size())
	}

	// Remove remaining items
	val, ok = rb.Get()
	if !ok || val != 3 {
		t.Errorf("Get() should return 3, got %v", val)
	}
	val, ok = rb.Get()
	if !ok || val != 4 {
		t.Errorf("Get() should return 4, got %v", val)
	}
	val, ok = rb.Get()
	if !ok || val != 5 {
		t.Errorf("Get() should return 5, got %v", val)
	}

	if rb.Size() != 0 {
		t.Errorf("Ring buffer size should be 0, got %d", rb.Size())
	}
	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty after getting all elements")
	}

	// Test Get after buffer wraps around and becomes empty again
	_, ok = rb.Get()
	if ok {
		t.Error("Get() on empty buffer should return false")
	}
}

func TestPeek(t *testing.T) {
	rb, _ := New[int](3)

	// Test Peek on empty buffer
	_, ok := rb.Peek()
	if ok {
		t.Error("Peek() on empty buffer should return false")
	}

	rb.Put(1)
	rb.Put(2)

	val, ok := rb.Peek()
	if !ok || val != 1 {
		t.Errorf("Peek() should return 1, got %v", val)
	}
	if rb.Size() != 2 {
		t.Errorf("Peek() should not change size, got %d", rb.Size())
	}

	// Get one item, peek should return the next
	rb.Get()
	val, ok = rb.Peek()
	if !ok || val != 2 {
		t.Errorf("Peek() should return 2 after getting first item, got %v", val)
	}
}

func TestIsFull(t *testing.T) {
	rb, _ := New[int](2)

	if rb.IsFull() {
		t.Error("Ring buffer should not be full initially")
	}

	rb.Put(1)
	if rb.IsFull() {
		t.Error("Ring buffer should not be full after putting 1 element in capacity 2")
	}

	rb.Put(2)
	if !rb.IsFull() {
		t.Error("Ring buffer should be full after putting 2 elements in capacity 2")
	}
}

func TestSize(t *testing.T) {
	rb, _ := New[int](3)

	if rb.Size() != 0 {
		t.Errorf("New ring buffer size should be 0, got %d", rb.Size())
	}

	rb.Put(1)
	if rb.Size() != 1 {
		t.Errorf("Ring buffer size should be 1, got %d", rb.Size())
	}

	rb.Put(2)
	if rb.Size() != 2 {
		t.Errorf("Ring buffer size should be 2, got %d", rb.Size())
	}

	rb.Get()
	if rb.Size() != 1 {
		t.Errorf("Ring buffer size should be 1, got %d", rb.Size())
	}

	rb.Get()
	if rb.Size() != 0 {
		t.Errorf("Ring buffer size should be 0, got %d", rb.Size())
	}
}

func TestCapacity(t *testing.T) {
	rb, _ := New[int](5)

	if rb.Capacity() != 5 {
		t.Errorf("Ring buffer capacity should be 5, got %d", rb.Capacity())
	}
}

func TestClear(t *testing.T) {
	rb, _ := New[int](3)

	rb.Put(1)
	rb.Put(2)
	rb.Put(3)

	rb.Clear()

	if rb.Size() != 0 {
		t.Errorf("Ring buffer size should be 0 after Clear, got %d", rb.Size())
	}
	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty after Clear")
	}
	if rb.IsFull() {
		t.Error("Ring buffer should not be full after Clear")
	}

	// Ensure positions are reset
	rb.Put(1)
	if rb.Size() != 1 {
		t.Errorf("Ring buffer size should be 1 after putting after Clear, got %d", rb.Size())
	}
}
