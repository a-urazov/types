package rbtree

import (
	"testing"
)

func TestRBTree_SetGet(t *testing.T) {
	tree := New[int, string]()

	tree.Set(10, "ten")
	tree.Set(5, "five")
	tree.Set(15, "fifteen")

	if val, ok := tree.Get(10); !ok || val != "ten" {
		t.Errorf("Expected 'ten', got %v, %v", val, ok)
	}

	if val, ok := tree.Get(5); !ok || val != "five" {
		t.Errorf("Expected 'five', got %v, %v", val, ok)
	}

	if val, ok := tree.Get(15); !ok || val != "fifteen" {
		t.Errorf("Expected 'fifteen', got %v, %v", val, ok)
	}

	if val, ok := tree.Get(99); ok || val != "" {
		t.Errorf("Expected empty string and false, got %v, %v", val, ok)
	}
}

func TestRBTree_Contains(t *testing.T) {
	tree := New[int, string]()

	tree.Set(10, "ten")
	tree.Set(5, "five")

	if !tree.Contains(10) {
		t.Errorf("Expected key 10 to be contained")
	}

	if !tree.Contains(5) {
		t.Errorf("Expected key 5 to be contained")
	}

	if tree.Contains(99) {
		t.Errorf("Expected key 99 to not be contained")
	}
}

func TestRBTree_Size(t *testing.T) {
	tree := New[int, string]()

	if tree.Size() != 0 {
		t.Errorf("Expected size 0, got %d", tree.Size())
	}

	tree.Set(10, "ten")
	if tree.Size() != 1 {
		t.Errorf("Expected size 1, got %d", tree.Size())
	}

	tree.Set(5, "five")
	if tree.Size() != 2 {
		t.Errorf("Expected size 2, got %d", tree.Size())
	}

	tree.Set(10, "ten-updated") // Update, not insert
	if tree.Size() != 2 {
		t.Errorf("Expected size 2, got %d", tree.Size())
	}
}

func TestRBTree_IsEmpty(t *testing.T) {
	tree := New[int, string]()

	if !tree.IsEmpty() {
		t.Errorf("Expected tree to be empty")
	}

	tree.Set(10, "ten")
	if tree.IsEmpty() {
		t.Errorf("Expected tree to not be empty")
	}

	// Since deletion is not fully implemented, we can't test the empty-after-delete case
}

func TestRBTree_Delete(t *testing.T) {
	tree := New[int, string]()

	tree.Set(10, "ten")
	tree.Set(5, "five")
	tree.Set(15, "fifteen")

	if tree.Size() != 3 {
		t.Errorf("Expected size 3, got %d", tree.Size())
	}

	// For this simplified implementation, deletion is not fully implemented
	// so we expect it to return false
	if tree.Delete(5) {
		t.Errorf("Expected deletion to return false (not yet implemented)")
	}

	if tree.Size() != 3 { // Size should remain the same since deletion isn't implemented
		t.Errorf("Expected size to remain 3 (deletion not implemented), got %d", tree.Size())
	}

	if !tree.Contains(5) { // Key should still be present since deletion isn't implemented
		t.Errorf("Expected key 5 to still be present (deletion not implemented)")
	}

	if tree.Delete(99) { // Trying to delete non-existent key
		t.Errorf("Expected deletion of non-existent key to return false")
	}
}

func TestRBTree_ConcurrentAccess(t *testing.T) {
	tree := New[int, string]()

	// Test concurrent access doesn't cause race conditions
	done := make(chan bool, 10)

	// Multiple goroutines setting values
	for i := 0; i < 5; i++ {
		go func(val int) {
			tree.Set(val, string(rune('A'+val)))
			done <- true
		}(i)
	}

	// Multiple goroutines getting values
	for i := 0; i < 5; i++ {
		go func(val int) {
			tree.Get(val)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRBTree_InOrderTraversal(t *testing.T) {
	tree := New[int, string]()

	// Insert values in random order
	tree.Set(10, "ten")
	tree.Set(5, "five")
	tree.Set(15, "fifteen")
	tree.Set(3, "three")
	tree.Set(7, "seven")

	results := make([]int, 0, 5)
	tree.InOrderTraversal(func(key int, value string) {
		results = append(results, key)
	})

	// In-order traversal should give us keys in sorted order
	expected := []int{3, 5, 7, 10, 15}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("Expected %v at index %d, got %v", v, i, results[i])
		}
	}
}
