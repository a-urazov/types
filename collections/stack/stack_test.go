package stack

import (
	"sync"
	"testing"
)

func TestNewStack(t *testing.T) {
	s := New[int]()
	if !s.IsEmpty() || s.Size() != 0 {
		t.Error("New stack should be empty")
	}
}

func TestPush(t *testing.T) {
	s := New[string]()
	s.Push("first")
	s.Push("second")
	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}
	val, _ := s.Peek()
	if val != "second" {
		t.Errorf("Expected to peek 'second', got '%s'", val)
	}
}

func TestPop(t *testing.T) {
	s := New[int]()
	s.Push(10)
	s.Push(20)

	val, ok := s.Pop()
	if !ok || val != 20 {
		t.Errorf("Expected to pop 20, got %d", val)
	}
	if s.Size() != 1 {
		t.Errorf("Expected size 1 after pop, got %d", s.Size())
	}

	val, ok = s.Pop()
	if !ok || val != 10 {
		t.Errorf("Expected to pop 10, got %d", val)
	}
	if !s.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}

	_, ok = s.Pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}
}

func TestPeekStack(t *testing.T) {
	s := New[float64]()
	s.Push(3.14)

	val, ok := s.Peek()
	if !ok || val != 3.14 {
		t.Errorf("Peek should return 3.14, got %v", val)
	}
	if s.Size() != 1 {
		t.Error("Peek should not remove the item from the stack")
	}

	s.Pop()
	_, ok = s.Peek()
	if ok {
		t.Error("Peek on empty stack should return false")
	}
}

func TestClearStack(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Clear()
	if !s.IsEmpty() {
		t.Error("Stack should be empty after Clear")
	}
}

func TestStackConcurrency(t *testing.T) {
	s := New[int]()
	var wg sync.WaitGroup

	// Concurrent pushes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Push(i)
		}(i)
	}
	wg.Wait()

	if s.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent pushes, got %d", s.Size())
	}

	// Concurrent pops
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Pop()
		}()
	}
	wg.Wait()

	if s.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent pops, got %d", s.Size())
	}
}
