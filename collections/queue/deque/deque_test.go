package deque

import (
	"testing"
)

func TestNewDeque(t *testing.T) {
	d := New[int]()
	if d == nil {
		t.Error("New() should not return nil")
	}
	if !d.IsEmpty() {
		t.Error("New deque should be empty")
	}
	if d.Size() != 0 {
		t.Errorf("New deque size should be 0, got %d", d.Size())
	}
}

func TestPushFront(t *testing.T) {
	d := New[int]()
	d.PushFront(1)
	if d.IsEmpty() {
		t.Error("Deque should not be empty after PushFront")
	}
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}
	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() should return 1, got %v", val)
	}

	d.PushFront(2)
	if d.Size() != 2 {
		t.Errorf("Deque size should be 2, got %d", d.Size())
	}
	val, ok = d.Front()
	if !ok || val != 2 {
		t.Errorf("Front() should return 2, got %v", val)
	}
}

func TestPushBack(t *testing.T) {
	d := New[int]()
	d.PushBack(1)
	if d.IsEmpty() {
		t.Error("Deque should not be empty after PushBack")
	}
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}
	val, ok := d.Back()
	if !ok || val != 1 {
		t.Errorf("Back() should return 1, got %v", val)
	}

	d.PushBack(2)
	if d.Size() != 2 {
		t.Errorf("Deque size should be 2, got %d", d.Size())
	}
	val, ok = d.Back()
	if !ok || val != 2 {
		t.Errorf("Back() should return 2, got %v", val)
	}
}

func TestPopFront(t *testing.T) {
	d := New[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	val, ok := d.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() should return 1, got %v", val)
	}
	if d.Size() != 2 {
		t.Errorf("Deque size should be 2, got %d", d.Size())
	}

	val, ok = d.PopFront()
	if !ok || val != 2 {
		t.Errorf("PopFront() should return 2, got %v", val)
	}
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}

	val, ok = d.PopFront()
	if !ok || val != 3 {
		t.Errorf("PopFront() should return 3, got %v", val)
	}
	if d.Size() != 0 {
		t.Errorf("Deque size should be 0, got %d", d.Size())
	}
	if !d.IsEmpty() {
		t.Error("Deque should be empty after popping all elements")
	}

	_, ok = d.PopFront()
	if ok {
		t.Error("PopFront() on empty deque should return false")
	}
}

func TestPopBack(t *testing.T) {
	d := New[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	val, ok := d.PopBack()
	if !ok || val != 3 {
		t.Errorf("PopBack() should return 3, got %v", val)
	}
	if d.Size() != 2 {
		t.Errorf("Deque size should be 2, got %d", d.Size())
	}

	val, ok = d.PopBack()
	if !ok || val != 2 {
		t.Errorf("PopBack() should return 2, got %v", val)
	}
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}

	val, ok = d.PopBack()
	if !ok || val != 1 {
		t.Errorf("PopBack() should return 1, got %v", val)
	}
	if d.Size() != 0 {
		t.Errorf("Deque size should be 0, got %d", d.Size())
	}
	if !d.IsEmpty() {
		t.Error("Deque should be empty after popping all elements")
	}

	_, ok = d.PopBack()
	if ok {
		t.Error("PopBack() on empty deque should return false")
	}
}

func TestFront(t *testing.T) {
	d := New[int]()
	_, ok := d.Front()
	if ok {
		t.Error("Front() on empty deque should return false")
	}

	d.PushBack(1)
	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() should return 1, got %v", val)
	}
	d.PushBack(2)
	val, ok = d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() should still return 1, got %v", val)
	}
}

func TestBack(t *testing.T) {
	d := New[int]()
	_, ok := d.Back()
	if ok {
		t.Error("Back() on empty deque should return false")
	}

	d.PushFront(1)
	val, ok := d.Back()
	if !ok || val != 1 {
		t.Errorf("Back() should return 1, got %v", val)
	}
	d.PushFront(2)
	val, ok = d.Back()
	if !ok || val != 1 {
		t.Errorf("Back() should still return 1, got %v", val)
	}
}

func TestIsEmpty(t *testing.T) {
	d := New[int]()
	if !d.IsEmpty() {
		t.Error("New deque should be empty")
	}
	d.PushFront(1)
	if d.IsEmpty() {
		t.Error("Deque should not be empty after PushFront")
	}
	d.PopFront()
	if !d.IsEmpty() {
		t.Error("Deque should be empty after PopFront")
	}
}

func TestSize(t *testing.T) {
	d := New[int]()
	if d.Size() != 0 {
		t.Errorf("New deque size should be 0, got %d", d.Size())
	}
	d.PushFront(1)
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}
	d.PushBack(2)
	if d.Size() != 2 {
		t.Errorf("Deque size should be 2, got %d", d.Size())
	}
	d.PopFront()
	if d.Size() != 1 {
		t.Errorf("Deque size should be 1, got %d", d.Size())
	}
	d.PopBack()
	if d.Size() != 0 {
		t.Errorf("Deque size should be 0, got %d", d.Size())
	}
}
