package concurrentmap

import (
	"sync"
	"testing"
)

func TestConcurrentMapBasicOperations(t *testing.T) {
	cm := New[int, string]()

	// Test empty map
	if !cm.IsEmpty() {
		t.Error("Expected empty map")
	}
	if cm.Len() != 0 {
		t.Errorf("Expected length 0, got %d", cm.Len())
	}

	// Test Set and Get
	cm.Set(1, "one")
	cm.Set(2, "two")
	cm.Set(3, "three")

	if cm.IsEmpty() {
		t.Error("Expected non-empty map")
	}
	if cm.Len() != 3 {
		t.Errorf("Expected length 3, got %d", cm.Len())
	}

	// Test Get
	if val, ok := cm.Get(1); !ok || val != "one" {
		t.Errorf("Expected Get(1) to return 'one', got %s", val)
	}
	if val, ok := cm.Get(2); !ok || val != "two" {
		t.Errorf("Expected Get(2) to return 'two', got %s", val)
	}
	if val, ok := cm.Get(3); !ok || val != "three" {
		t.Errorf("Expected Get(3) to return 'three', got %s", val)
	}

	// Test non-existent key
	if _, ok := cm.Get(4); ok {
		t.Error("Expected Get(4) to return false")
	}

	// Test Contains
	if !cm.Contains(1) || !cm.Contains(2) || !cm.Contains(3) {
		t.Error("Expected all keys to exist")
	}
	if cm.Contains(4) {
		t.Error("Expected key 4 to not exist")
	}

	// Test Update
	if !cm.Update(1, func(v string) string { return v + "-updated" }) {
		t.Error("Expected Update to return true for existing key")
	}
	if val, ok := cm.Get(1); !ok || val != "one-updated" {
		t.Errorf("Expected updated value 'one-updated', got %s", val)
	}

	// Test Update on non-existent key
	if cm.Update(4, func(v string) string { return v + "-updated" }) {
		t.Error("Expected Update to return false for non-existent key")
	}

	// Test Delete
	if !cm.Delete(2) {
		t.Error("Expected Delete(2) to return true")
	}
	if cm.Contains(2) {
		t.Error("Expected key 2 to be deleted")
	}
	if cm.Len() != 2 {
		t.Errorf("Expected length 2 after deletion, got %d", cm.Len())
	}

	// Test Delete on non-existent key
	if cm.Delete(2) {
		t.Error("Expected Delete(2) to return false for non-existent key")
	}
}

func TestConcurrentMapTrySetAndGetOrSet(t *testing.T) {
	cm := New[string, int]()

	// Test TrySet
	if !cm.TrySet("key1", 10) {
		t.Error("Expected TrySet to return true for new key")
	}
	if !cm.Contains("key1") {
		t.Error("Expected key1 to exist after TrySet")
	}
	if val, _ := cm.Get("key1"); val != 10 {
		t.Errorf("Expected value 10, got %d", val)
	}

	// Try to set same key again
	if cm.TrySet("key1", 20) {
		t.Error("Expected TrySet to return false for existing key")
	}
	if val, _ := cm.Get("key1"); val != 10 {
		t.Errorf("Expected value to remain 10, got %d", val)
	}

	// Test GetOrSet
	if val, existed := cm.GetOrSet("key2", 20); existed {
		t.Error("Expected GetOrSet to return false for new key")
	} else if val != 20 {
		t.Errorf("Expected GetOrSet to return value 20, got %d", val)
	}

	// Try GetOrSet on existing key
	if val, existed := cm.GetOrSet("key1", 30); !existed {
		t.Error("Expected GetOrSet to return true for existing key")
	} else if val != 10 {
		t.Errorf("Expected GetOrSet to return original value 10, got %d", val)
	}
}

func TestConcurrentMapKeysValuesItems(t *testing.T) {
	cm := New[int, string]()

	cm.Set(1, "one")
	cm.Set(2, "two")
	cm.Set(3, "three")

	// Test Keys
	keys := cm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Test Values
	values := cm.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test Items
	items := cm.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify all keys and values are present
	expected := map[int]string{1: "one", 2: "two", 3: "three"}
	for k, v := range items {
		if expected[k] != v {
			t.Errorf("Expected items[%d] = '%s', got '%s'", k, expected[k], v)
		}
	}
}

func TestConcurrentMapForEach(t *testing.T) {
	cm := New[int, int]()

	for i := 1; i <= 5; i++ {
		cm.Set(i, i*i) // Store square of i
	}

	sum := 0
	count := 0
	cm.ForEach(func(key int, value int) {
		if key*key != value {
			t.Errorf("Expected value %d for key %d, got %d", key*key, key, value)
		}
		sum += value
		count++
	})

	if count != 5 {
		t.Errorf("Expected to iterate over 5 items, got %d", count)
	}

	expectedSum := 1 + 4 + 9 + 16 + 25 // 1^2 + 2^2 + 3^2 + 4^2 + 5^2
	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}
}

func TestConcurrentMapClear(t *testing.T) {
	cm := New[string, int]()

	cm.Set("a", 1)
	cm.Set("b", 2)
	cm.Set("c", 3)

	if cm.IsEmpty() || cm.Len() != 3 {
		t.Error("Expected non-empty map with length 3")
	}

	cm.Clear()

	if !cm.IsEmpty() || cm.Len() != 0 {
		t.Error("Expected empty map after Clear()")
	}
}

func TestConcurrentMapSharding(t *testing.T) {
	cm := NewWithShardCount[int, int](4) // 4 shards

	// Add several items
	for i := 0; i < 100; i++ {
		cm.Set(i, i*2)
	}

	if cm.Len() != 100 {
		t.Errorf("Expected length 100, got %d", cm.Len())
	}

	// Verify all items are accessible
	for i := 0; i < 100; i++ {
		if val, ok := cm.Get(i); !ok || val != i*2 {
			t.Errorf("Expected Get(%d) to return %d, got %d", i, i*2, val)
		}
	}

	// Verify we have the expected number of shards
	if cm.Size() != 4 {
		t.Errorf("Expected 4 shards, got %d", cm.Size())
	}
}

func TestConcurrentMapConcurrency(t *testing.T) {
	cm := New[int, int]()
	const numGoroutines = 10
	const opsPerGoroutine = 100
	var wg sync.WaitGroup

	// Launch goroutines to perform concurrent operations
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			start := goroutineID * opsPerGoroutine
			end := start + opsPerGoroutine
			for i := start; i < end; i++ {
				cm.Set(i, i*2)
			}
		}(g)
	}

	wg.Wait()

	// Verify all items were set correctly
	if cm.Len() != numGoroutines*opsPerGoroutine {
		t.Errorf("Expected length %d, got %d", numGoroutines*opsPerGoroutine, cm.Len())
	}

	for i := 0; i < numGoroutines*opsPerGoroutine; i++ {
		if val, ok := cm.Get(i); !ok || val != i*2 {
			t.Errorf("Expected Get(%d) to return %d, got %d", i, i*2, val)
		}
	}

	// Launch goroutines to perform concurrent reads and writes
	var readWg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		readWg.Add(1)
		go func(goroutineID int) {
			defer readWg.Done()
			start := goroutineID * opsPerGoroutine
			end := start + opsPerGoroutine
			for i := start; i < end; i++ {
				// Read
				_, _ = cm.Get(i)
				// Update
				cm.Update(i, func(v int) int { return v + 1 })
			}
		}(g)
	}

	readWg.Wait()

	// Verify all items were updated correctly
	for i := 0; i < numGoroutines*opsPerGoroutine; i++ {
		expected := i*2 + 1 // original value plus 1 from update
		if val, ok := cm.Get(i); !ok || val != expected {
			t.Errorf("Expected Get(%d) to return %d, got %d", i, expected, val)
		}
	}
}

func TestConcurrentMapResizeShards(t *testing.T) {
	cm := NewWithShardCount[string, int](2) // Start with 2 shards

	// Add some items
	for i := 0; i < 50; i++ {
		cm.Set(string(rune('a'+i)), i)
	}

	// Resize to more shards
	cm.ResizeShards(8)

	// Verify all items are still accessible
	if cm.Len() != 50 {
		t.Errorf("Expected length 50 after resize, got %d", cm.Len())
	}

	for i := 0; i < 50; i++ {
		key := string(rune('a' + i))
		if val, ok := cm.Get(key); !ok || val != i {
			t.Errorf("Expected Get('%s') to return %d, got %d", key, i, val)
		}
	}

	// Verify new shard count
	if cm.Size() != 8 {
		t.Errorf("Expected 8 shards after resize, got %d", cm.Size())
	}
}