package linked

import (
	"sync"
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	l := New[int]()
	if l.Size() != 0 {
		t.Error("New linked list should be empty")
	}
}

func TestAddFirst(t *testing.T) {
	l := New[string]()
	l.AddFirst("b")
	l.AddFirst("a")
	if l.Size() != 2 {
		t.Errorf("Expected size 2, got %d", l.Size())
	}
	first, _ := l.First()
	if first != "a" {
		t.Errorf("Expected first element to be 'a', got '%s'", first)
	}
	last, _ := l.Last()
	if last != "b" {
		t.Errorf("Expected last element to be 'b', got '%s'", last)
	}
}

func TestAddLast(t *testing.T) {
	l := New[int]()
	l.AddLast(1)
	l.AddLast(2)
	first, _ := l.First()
	if first != 1 {
		t.Errorf("Expected first element to be 1, got %d", first)
	}
	last, _ := l.Last()
	if last != 2 {
		t.Errorf("Expected last element to be 2, got %d", last)
	}
}

func TestRemoveFirst(t *testing.T) {
	l := New[int]()
	l.AddLast(1)
	l.AddLast(2)

	if !l.RemoveFirst() {
		t.Error("Failed to remove first element")
	}
	if l.Size() != 1 {
		t.Errorf("Expected size 1, got %d", l.Size())
	}
	first, _ := l.First()
	if first != 2 {
		t.Errorf("New first element should be 2, got %d", first)
	}

	l.RemoveFirst()
	if l.Size() != 0 {
		t.Error("List should be empty")
	}
	if l.RemoveFirst() {
		t.Error("Should not be able to remove from empty list")
	}
}

func TestRemoveLast(t *testing.T) {
	l := New[string]()
	l.AddLast("a")
	l.AddLast("b")

	if !l.RemoveLast() {
		t.Error("Failed to remove last element")
	}
	if l.Size() != 1 {
		t.Errorf("Expected size 1, got %d", l.Size())
	}
	last, _ := l.Last()
	if last != "a" {
		t.Errorf("New last element should be 'a', got '%s'", last)
	}
}

func TestClearLinkedList(t *testing.T) {
	l := New[int]()
	l.AddLast(100)
	l.Clear()
	if l.Size() != 0 {
		t.Error("List should be empty after Clear")
	}
	if _, ok := l.First(); ok {
		t.Error("First should return false for empty list")
	}
}

func TestLinkedListConcurrency(t *testing.T) {
	l := New[int]()
	var wg sync.WaitGroup

	// Concurrent AddFirst and AddLast
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			l.AddFirst(i)
		}(i)
		go func(i int) {
			defer wg.Done()
			l.AddLast(i + 100)
		}(i)
	}
	wg.Wait()

	if l.Size() != 200 {
		t.Errorf("Expected size 200, got %d", l.Size())
	}

	// Concurrent RemoveFirst and RemoveLast
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			l.RemoveFirst()
		}()
		go func() {
			defer wg.Done()
			l.RemoveLast()
		}()
	}
	wg.Wait()

	if l.Size() != 100 {
		t.Errorf("Expected size 100 after removals, got %d", l.Size())
	}
}
