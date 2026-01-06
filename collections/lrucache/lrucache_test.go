package lrucache

import (
	"testing"
)

func TestLRUCacheBasicOperations(t *testing.T) {
	cache := New[string, int](3)

	// Test empty cache
	if !cache.IsEmpty() {
		t.Error("Expected empty cache")
	}
	if cache.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cache.Size())
	}

	// Add items
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	// Get items
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected Get('a') to return 1, got %d", val)
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected Get('b') to return 2, got %d", val)
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("Expected Get('c') to return 3, got %d", val)
	}

	// Test Contains
	if !cache.Contains("a") || !cache.Contains("b") || !cache.Contains("c") {
		t.Error("Expected all keys to be present")
	}
	if cache.Contains("d") {
		t.Error("Expected key 'd' to be absent")
	}

	// Update existing item
	cache.Put("a", 10)
	if val, ok := cache.Get("a"); !ok || val != 10 {
		t.Errorf("Expected Get('a') to return 10 after update, got %d", val)
	}

	// Remove item
	if !cache.Remove("b") {
		t.Error("Expected Remove('b') to return true")
	}
	if cache.Contains("b") {
		t.Error("Expected key 'b' to be removed")
	}
	if cache.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", cache.Size())
	}

	// Try to remove non-existent item
	if cache.Remove("x") {
		t.Error("Expected Remove('x') to return false")
	}
}

func TestLRUCacheEviction(t *testing.T) {
	cache := New[string, int](2)

	// Fill cache to capacity
	cache.Put("a", 1)
	cache.Put("b", 2)

	// Access 'a' to make it most recently used
	cache.Get("a")

	// Add new item - should evict 'b' (least recently used)
	cache.Put("c", 3)

	// Check that 'b' is evicted and others remain
	if cache.Contains("b") {
		t.Error("Expected 'b' to be evicted")
	}
	if !cache.Contains("a") || !cache.Contains("c") {
		t.Error("Expected 'a' and 'c' to remain")
	}
	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}

	// Add another item - should evict 'a' (now least recently used)
	cache.Put("d", 4)
	if cache.Contains("a") {
		t.Error("Expected 'a' to be evicted")
	}
	if !cache.Contains("c") || !cache.Contains("d") {
		t.Error("Expected 'c' and 'd' to remain")
	}
}

func TestLRUCachePeek(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	// Peek should not affect LRU order
	// Current order: MRU=b, LRU=a
	if val, ok := cache.Peek("a"); !ok || val != 1 {
		t.Errorf("Expected Peek('a') to return 1, got %d", val)
	}

	// Add new item - should evict 'a' (least recently used, Peek doesn't change order)
	cache.Put("c", 3)

	if cache.Contains("a") {
		t.Error("Expected 'a' to be evicted")
	}
	if !cache.Contains("b") || !cache.Contains("c") {
		t.Error("Expected 'b' and 'c' to remain")
	}
}

func TestLRUCacheResize(t *testing.T) {
	cache := New[string, int](5)

	// Add more items than new capacity
	for i := 0; i < 5; i++ {
		cache.Put(string(rune('a'+i)), i+1)
	}

	// Resize to smaller capacity
	cache.Resize(2)

	if cache.Size() != 2 {
		t.Errorf("Expected size 2 after resize, got %d", cache.Size())
	}

	// The two most recently added items should remain
	if !cache.Contains("d") || !cache.Contains("e") {
		t.Error("Expected 'd' and 'e' to remain after resize")
	}

	// Resize to larger capacity
	cache.Resize(10)
	cache.Put("f", 6)
	cache.Put("g", 7)

	if cache.Size() != 4 {
		t.Errorf("Expected size 4 after adding to larger cache, got %d", cache.Size())
	}
}

func TestLRUCacheKeysAndValues(t *testing.T) {
	cache := New[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Access 'a' to make it most recent
	cache.Get("a")

	// Keys and Values should be in LRU order (least recent first)
	keys := cache.Keys()
	expectedKeys := []string{"b", "c", "a"}
	for i, key := range expectedKeys {
		if keys[i] != key {
			t.Errorf("Expected key %s at index %d, got %s", key, i, keys[i])
		}
	}

	values := cache.Values()
	expectedValues := []int{2, 3, 1}
	for i, val := range expectedValues {
		if values[i] != val {
			t.Errorf("Expected value %d at index %d, got %d", val, i, values[i])
		}
	}
}

func TestLRUCacheForEach(t *testing.T) {
	cache := New[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Access 'a' to make it most recent
	cache.Get("a")

	// Test ForEach (LRU order)
	var forEachKeys []string
	var forEachValues []int
	cache.ForEach(func(key string, value int) {
		forEachKeys = append(forEachKeys, key)
		forEachValues = append(forEachValues, value)
	})

	expectedKeys := []string{"b", "c", "a"}
	expectedValues := []int{2, 3, 1}
	for i, key := range expectedKeys {
		if forEachKeys[i] != key {
			t.Errorf("ForEach: Expected key %s at index %d, got %s", key, i, forEachKeys[i])
		}
		if forEachValues[i] != expectedValues[i] {
			t.Errorf("ForEach: Expected value %d at index %d, got %d", expectedValues[i], i, forEachValues[i])
		}
	}

	// Test ForEachMRU (MRU order)
	var forEachMRUKeys []string
	var forEachMRUValues []int
	cache.ForEachMRU(func(key string, value int) {
		forEachMRUKeys = append(forEachMRUKeys, key)
		forEachMRUValues = append(forEachMRUValues, value)
	})

	expectedMRUKeys := []string{"a", "c", "b"}
	expectedMRUValues := []int{1, 3, 2}
	for i, key := range expectedMRUKeys {
		if forEachMRUKeys[i] != key {
			t.Errorf("ForEachMRU: Expected key %s at index %d, got %s", key, i, forEachMRUKeys[i])
		}
		if forEachMRUValues[i] != expectedMRUValues[i] {
			t.Errorf("ForEachMRU: Expected value %d at index %d, got %d", expectedMRUValues[i], i, forEachMRUValues[i])
		}
	}
}

func TestLRUCacheClear(t *testing.T) {
	cache := New[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)

	cache.Clear()

	if !cache.IsEmpty() {
		t.Error("Expected empty cache after Clear()")
	}
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", cache.Size())
	}
	if cache.Contains("a") || cache.Contains("b") {
		t.Error("Expected all keys to be removed after Clear()")
	}
}

func TestLRUCacheEdgeCases(t *testing.T) {
	// Test with zero/negative capacity
	cache := New[string, int](0)
	cache.Put("a", 1)
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 with zero capacity (should default to 1), got %d", cache.Size())
	}

	// Test Get on empty cache
	emptyCache := New[string, int](2)
	if val, ok := emptyCache.Get("nonexistent"); ok {
		t.Errorf("Expected Get on empty cache to return false, got value %v", val)
	}

	// Test Remove on empty cache
	if emptyCache.Remove("nonexistent") {
		t.Error("Expected Remove on empty cache to return false")
	}

	// Test Peek on empty cache
	if val, ok := emptyCache.Peek("nonexistent"); ok {
		t.Errorf("Expected Peek on empty cache to return false, got value %v", val)
	}
}

func TestLRUCacheConcurrency(t *testing.T) {
	cache := New[int, string](100)
	done := make(chan bool)

	// Writers
	go func() {
		for i := 0; i < 50; i++ {
			cache.Put(i, string(rune('a'+i%26)))
		}
		done <- true
	}()

	// Readers
	go func() {
		for i := 0; i < 25; i++ {
			cache.Get(i)
			cache.Contains(i)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state
	if cache.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent operations, got %d", cache.Size())
	}

	// All keys 0-49 should be present
	for i := 0; i < 50; i++ {
		if !cache.Contains(i) {
			t.Errorf("Expected key %d to be present after concurrent operations", i)
		}
	}
}
