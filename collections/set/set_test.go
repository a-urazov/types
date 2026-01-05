package set

import (
	"sort"
	"sync"
	"testing"
)

func TestNewHashSet(t *testing.T) {
	s := New[int]()
	if s.Size() != 0 {
		t.Error("New hashset should be empty")
	}
}

func TestAddSet(t *testing.T) {
	s := New[string]()
	if !s.Add("hello") {
		t.Error("Failed to add new item 'hello'")
	}
	if s.Add("hello") {
		t.Error("Should not be able to add duplicate item 'hello'")
	}
	if s.Size() != 1 {
		t.Errorf("Expected size 1, got %d", s.Size())
	}
}

func TestRemoveSet(t *testing.T) {
	s := New[int]()
	s.Add(123)
	if !s.Remove(123) {
		t.Error("Failed to remove item 123")
	}
	if s.Size() != 0 {
		t.Error("Set should be empty after removal")
	}
	if s.Remove(456) {
		t.Error("Should not be able to remove non-existent item")
	}
}

func TestContainsSet(t *testing.T) {
	s := New[string]()
	s.Add("world")
	if !s.Contains("world") {
		t.Error("Set should contain 'world'")
	}
	if s.Contains("go") {
		t.Error("Set should not contain 'go'")
	}
}

func TestClearSet(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Clear()
	if s.Size() != 0 {
		t.Error("Set should be empty after Clear")
	}
}

func TestToArraySet(t *testing.T) {
	s := New[string]()
	s.Add("a")
	s.Add("c")
	s.Add("b")

	arr := s.ToArray()
	sort.Strings(arr)

	expected := []string{"a", "b", "c"}
	if len(arr) != len(expected) {
		t.Fatalf("Expected array of length %d, got %d", len(expected), len(arr))
	}
	for i, v := range arr {
		if v != expected[i] {
			t.Errorf("Expected '%s' at index %d, got '%s'", expected[i], i, v)
		}
	}
}

func TestSetConcurrency(t *testing.T) {
	s := New[int]()
	var wg sync.WaitGroup

	// Concurrent adds
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Add(i)
			s.Add(i) // Add again to test uniqueness
		}(i)
	}
	wg.Wait()

	if s.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent adds, got %d", s.Size())
	}

	// Concurrent removes and contains
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			s.Remove(i)
		}(i)
		go func(i int) {
			defer wg.Done()
			s.Contains(i + 50)
		}(i)
	}
	wg.Wait()

	if s.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent operations, got %d", s.Size())
	}
	if s.Contains(25) {
		t.Error("Item 25 should have been removed")
	}
	if !s.Contains(75) {
		t.Error("Item 75 should still exist")
	}
}
