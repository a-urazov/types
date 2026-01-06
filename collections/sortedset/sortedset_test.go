package sortedset

import (
	"testing"
)

func TestSortedSetBasicOperations(t *testing.T) {
	// Create a new sorted set for integers
	set := New(func(a, b int) bool { return a < b })

	// Test empty set
	if !set.IsEmpty() {
		t.Errorf("Expected empty set, got size %d", set.Size())
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0, got %d", set.Size())
	}

	// Test adding elements
	if !set.Add(3) {
		t.Error("Expected Add(3) to return true")
	}
	if !set.Add(1) {
		t.Error("Expected Add(1) to return true")
	}
	if !set.Add(4) {
		t.Error("Expected Add(4) to return true")
	}
	if !set.Add(2) {
		t.Error("Expected Add(2) to return true")
	}
	if set.Add(3) {
		t.Error("Expected Add(3) to return false (duplicate)")
	}

	// Test size and contains
	if set.Size() != 4 {
		t.Errorf("Expected size 4, got %d", set.Size())
	}
	if !set.Contains(1) || !set.Contains(2) || !set.Contains(3) || !set.Contains(4) {
		t.Error("Expected all added elements to be present")
	}
	if set.Contains(5) {
		t.Error("Expected 5 not to be present")
	}

	// Test ToSlice returns elements in sorted order
	slice := set.ToSlice()
	expected := []int{1, 2, 3, 4}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, slice[i])
		}
	}

	// Test First and Last
	if first, ok := set.First(); !ok || first != 1 {
		t.Errorf("Expected First() to return 1, got %d", first)
	}
	if last, ok := set.Last(); !ok || last != 4 {
		t.Errorf("Expected Last() to return 4, got %d", last)
	}
}

func TestSortedSetRemove(t *testing.T) {
	set := New(func(a, b int) bool { return a < b })

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)

	// Remove existing element
	if !set.Remove(2) {
		t.Error("Expected Remove(2) to return true")
	}
	if set.Contains(2) {
		t.Error("Expected 2 to be removed")
	}
	if set.Size() != 3 {
		t.Errorf("Expected size 3 after removal, got %d", set.Size())
	}

	// Remove non-existing element
	if set.Remove(5) {
		t.Error("Expected Remove(5) to return false")
	}
	if set.Size() != 3 {
		t.Errorf("Expected size to remain 3, got %d", set.Size())
	}

	// Remove all elements
	set.Remove(1)
	set.Remove(3)
	set.Remove(4)
	if !set.IsEmpty() {
		t.Error("Expected set to be empty after removing all elements")
	}
}

func TestSortedSetCeilingFloor(t *testing.T) {
	set := New(func(a, b int) bool { return a < b })

	// Add elements: 1, 3, 5, 7, 9
	elements := []int{1, 3, 5, 7, 9}
	for _, e := range elements {
		set.Add(e)
	}

	// Test Ceiling
	if val, ok := set.Ceiling(0); !ok || val != 1 {
		t.Errorf("Ceiling(0) should return 1, got %d", val)
	}
	if val, ok := set.Ceiling(2); !ok || val != 3 {
		t.Errorf("Ceiling(2) should return 3, got %d", val)
	}
	if val, ok := set.Ceiling(5); !ok || val != 5 {
		t.Errorf("Ceiling(5) should return 5, got %d", val)
	}
	if val, ok := set.Ceiling(10); ok {
		t.Errorf("Ceiling(10) should not exist, but got %d", val)
	}

	// Test Floor
	if val, ok := set.Floor(0); ok {
		t.Errorf("Floor(0) should not exist, but got %d", val)
	}
	if val, ok := set.Floor(2); !ok || val != 1 {
		t.Errorf("Floor(2) should return 1, got %d", val)
	}
	if val, ok := set.Floor(5); !ok || val != 5 {
		t.Errorf("Floor(5) should return 5, got %d", val)
	}
	if val, ok := set.Floor(10); !ok || val != 9 {
		t.Errorf("Floor(10) should return 9, got %d", val)
	}
}

func TestSortedSetClear(t *testing.T) {
	set := New(func(a, b int) bool { return a < b })

	set.Add(1)
	set.Add(2)
	set.Add(3)

	set.Clear()

	if !set.IsEmpty() {
		t.Error("Expected set to be empty after Clear()")
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", set.Size())
	}
}

func TestSortedSetForEach(t *testing.T) {
	set := New(func(a, b int) bool { return a < b })

	elements := []int{5, 2, 8, 1, 9, 3}
	for _, e := range elements {
		set.Add(e)
	}

	var result []int
	set.ForEach(func(item int) {
		result = append(result, item)
	})

	// Result should be sorted
	expected := []int{1, 2, 3, 5, 8, 9}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %d at index %d in ForEach result, got %d", v, i, result[i])
		}
	}
}

func TestSortedSetStringComparison(t *testing.T) {
	// Test with strings
	set := New(func(a, b string) bool { return a < b })

	strings := []string{"zebra", "apple", "banana", "cherry"}
	for _, s := range strings {
		set.Add(s)
	}

	slice := set.ToSlice()
	expected := []string{"apple", "banana", "cherry", "zebra"}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %s at index %d, got %s", v, i, slice[i])
		}
	}
}
