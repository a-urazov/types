package sorted

import (
	"sync"
	"testing"
)

func TestNewSortedDictionary(t *testing.T) {
	d := New[int, string]()
	if d.Size() != 0 {
		t.Error("New sorted dictionary should be empty")
	}
}

func TestSetAndGetSorted(t *testing.T) {
	d := New[int, string]()
	d.Set(3, "three")
	d.Set(1, "one")
	d.Set(2, "two")

	if d.Size() != 3 {
		t.Errorf("Expected size 3, got %d", d.Size())
	}

	// Check order
	keys := d.Keys()
	if keys[0] != 1 || keys[1] != 2 || keys[2] != 3 {
		t.Errorf("Keys are not sorted: %v", keys)
	}

	values := d.Values()
	if values[0] != "one" || values[1] != "two" || values[2] != "three" {
		t.Errorf("Values are not in correct order: %v", values)
	}

	// Test Get
	val, ok := d.Get(2)
	if !ok || val != "two" {
		t.Errorf("Failed to get value for key 2. Got: %s", val)
	}

	// Test Update
	d.Set(2, "deux")
	val, _ = d.Get(2)
	if val != "deux" {
		t.Errorf("Failed to update value. Expected 'deux', got '%s'", val)
	}
}

func TestRemoveSorted(t *testing.T) {
	d := New[string, int]()
	d.Set("b", 2)
	d.Set("a", 1)
	d.Set("c", 3)

	if !d.Remove("b") {
		t.Error("Failed to remove key 'b'")
	}
	if d.Size() != 2 {
		t.Errorf("Expected size 2, got %d", d.Size())
	}
	if d.ContainsKey("b") {
		t.Error("Dictionary should not contain key 'b' after removal")
	}
	keys := d.Keys()
	if keys[0] != "a" || keys[1] != "c" {
		t.Errorf("Keys are not correct after removal: %v", keys)
	}
}

func TestSortedDictionaryConcurrency(t *testing.T) {
	d := New[int, int]()
	var wg sync.WaitGroup

	// Concurrent sets
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.Set(i, i*10)
		}(i)
	}
	wg.Wait()

	if d.Size() != 100 {
		t.Errorf("Expected size 100, got %d", d.Size())
	}
	keys := d.Keys()
	if keys[50] != 50 {
		t.Error("Data corruption in keys")
	}
	val, _ := d.Get(50)
	if val != 500 {
		t.Error("Data corruption in values")
	}

	// Concurrent removes and gets
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			d.Remove(i)
		}(i)
		go func(i int) {
			defer wg.Done()
			d.Get(i + 50)
		}(i)
	}
	wg.Wait()

	if d.Size() != 50 {
		t.Errorf("Expected size 50 after removals, got %d", d.Size())
	}
	if d.ContainsKey(25) {
		t.Error("Key 25 should have been removed")
	}
}
