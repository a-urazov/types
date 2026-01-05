package dictionary

import (
	"sort"
	"strconv"
	"sync"
	"testing"
)

func TestNewDictionary(t *testing.T) {
	d := New[string, int]()
	if d.Size() != 0 {
		t.Errorf("New dictionary should be empty, but has size %d", d.Size())
	}
}

func TestSetAndGet(t *testing.T) {
	d := New[string, int]()
	d.Set("one", 1)
	d.Set("two", 2)

	if val, ok := d.Get("one"); !ok || val != 1 {
		t.Errorf("Expected to get 1 for key 'one', got %d", val)
	}
	if d.Size() != 2 {
		t.Errorf("Expected size 2, got %d", d.Size())
	}

	d.Set("one", 11) // Test update
	if val, ok := d.Get("one"); !ok || val != 11 {
		t.Errorf("Expected to get updated value 11, got %d", val)
	}
}

func TestRemove(t *testing.T) {
	d := New[string, float64]()
	d.Set("pi", 3.14)

	if !d.Remove("pi") {
		t.Error("Expected to remove 'pi', but it failed")
	}
	if d.Size() != 0 {
		t.Errorf("Expected size 0 after removing, got %d", d.Size())
	}
	if _, ok := d.Get("pi"); ok {
		t.Error("Dictionary should not contain 'pi' after removal")
	}

	if d.Remove("e") {
		t.Error("Should not be able to remove non-existent key")
	}
}

func TestContainsKey(t *testing.T) {
	d := New[int, bool]()
	d.Set(42, true)
	if !d.ContainsKey(42) {
		t.Error("Dictionary should contain key 42")
	}
	if d.ContainsKey(100) {
		t.Error("Dictionary should not contain key 100")
	}
}

func TestKeys(t *testing.T) {
	d := New[string, int]()
	d.Set("c", 3)
	d.Set("a", 1)
	d.Set("b", 2)

	keys := d.Keys()
	sort.Strings(keys) // Sort for predictable order

	expected := []string{"a", "b", "c"}
	if len(keys) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("Expected key '%s' at index %d, got '%s'", expected[i], i, k)
		}
	}
}

func TestValues(t *testing.T) {
	d := New[string, int]()
	d.Set("c", 3)
	d.Set("a", 1)
	d.Set("b", 2)

	values := d.Values()
	sort.Ints(values) // Sort for predictable order

	expected := []int{1, 2, 3}
	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected value %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestClear(t *testing.T) {
	d := New[int, string]()
	d.Set(1, "one")
	d.Clear()
	if d.Size() != 0 {
		t.Errorf("Dictionary should be empty after Clear, but size is %d", d.Size())
	}
	if d.ContainsKey(1) {
		t.Error("Dictionary should not contain key 1 after Clear")
	}
}

func TestConcurrencyDictionary(t *testing.T) {
	d := New[int, string]()
	var wg sync.WaitGroup

	// Concurrent sets
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.Set(i, "value_"+strconv.Itoa(i))
		}(i)
	}
	wg.Wait()

	if d.Size() != 1000 {
		t.Errorf("Expected size 1000 after concurrent sets, got %d", d.Size())
	}

	// Concurrent gets and removes
	for i := 0; i < 500; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			d.Get(i)
		}(i)
		go func(i int) {
			defer wg.Done()
			d.Remove(i + 500) // Remove the other half
		}(i)
	}
	wg.Wait()

	if d.Size() != 500 {
		t.Errorf("Expected size 500 after concurrent operations, got %d", d.Size())
	}
	if d.ContainsKey(750) {
		t.Error("Key 750 should have been removed")
	}
	if !d.ContainsKey(250) {
		t.Error("Key 250 should still exist")
	}
}
