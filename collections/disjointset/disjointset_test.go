package disjointset

import (
	"reflect"
	"testing"
)

func TestDisjointSetBasicOperations(t *testing.T) {
	ds := New[int]()

	// Test empty set
	if !ds.IsEmpty() {
		t.Error("Expected empty disjoint set")
	}
	if ds.Size() != 0 {
		t.Errorf("Expected size 0, got %d", ds.Size())
	}

	// Make sets
	ds.MakeSet(1)
	ds.MakeSet(2)
	ds.MakeSet(3)

	if ds.IsEmpty() {
		t.Error("Expected non-empty disjoint set")
	}
	if ds.Size() != 3 {
		t.Errorf("Expected size 3, got %d", ds.Size())
	}

	// Find should return the element itself for single-element sets
	if root, ok := ds.Find(1); !ok || root != 1 {
		t.Errorf("Expected Find(1) to return 1, got %d", root)
	}
	if root, ok := ds.Find(2); !ok || root != 2 {
		t.Errorf("Expected Find(2) to return 2, got %d", root)
	}
	if root, ok := ds.Find(3); !ok || root != 3 {
		t.Errorf("Expected Find(3) to return 3, got %d", root)
	}

	// Connected should return false for different elements
	if ds.Connected(1, 2) {
		t.Error("Expected 1 and 2 to be in different sets")
	}
	if ds.Connected(1, 3) {
		t.Error("Expected 1 and 3 to be in different sets")
	}
	if ds.Connected(2, 3) {
		t.Error("Expected 2 and 3 to be in different sets")
	}

	// Union sets
	if !ds.Union(1, 2) {
		t.Error("Expected Union(1, 2) to return true")
	}
	if ds.Union(1, 2) {
		t.Error("Expected Union(1, 2) to return false (already connected)")
	}

	// After union, 1 and 2 should have same root
	root1, ok1 := ds.Find(1)
	root2, ok2 := ds.Find(2)
	if !ok1 || !ok2 || root1 != root2 {
		t.Errorf("Expected 1 and 2 to have same root after union, got %d and %d", root1, root2)
	}

	// Connected should return true for 1 and 2
	if !ds.Connected(1, 2) {
		t.Error("Expected 1 and 2 to be connected after union")
	}

	// 3 should still be separate
	if ds.Connected(1, 3) {
		t.Error("Expected 1 and 3 to be in different sets")
	}
	if ds.Connected(2, 3) {
		t.Error("Expected 2 and 3 to be in different sets")
	}

	// Union with 3
	if !ds.Union(2, 3) {
		t.Error("Expected Union(2, 3) to return true")
	}

	// Now all should be connected
	if !ds.Connected(1, 3) {
		t.Error("Expected 1 and 3 to be connected after union")
	}
	if !ds.Connected(2, 3) {
		t.Error("Expected 2 and 3 to be connected after union")
	}

	// All should have same root
	root3, ok3 := ds.Find(3)
	if !ok3 || root1 != root3 {
		t.Errorf("Expected all elements to have same root, got %d, %d, %d", root1, root2, root3)
	}

	// Set count should be 1
	if ds.SetCount() != 1 {
		t.Errorf("Expected set count 1, got %d", ds.SetCount())
	}
}

func TestDisjointSetUnionByRank(t *testing.T) {
	ds := New[string]()

	// Create multiple elements
	elements := []string{"a", "b", "c", "d", "e"}
	for _, elem := range elements {
		ds.MakeSet(elem)
	}

	// Union in a way that tests rank balancing
	ds.Union("a", "b") // rank 1
	ds.Union("c", "d") // rank 1
	ds.Union("a", "c") // should make one tree with rank 2

	// All should be connected
	for i := 0; i < len(elements)-1; i++ {
		for j := i + 1; j < len(elements); j++ {
			if !ds.Connected(elements[i], elements[j]) && elements[j] != "e" {
				t.Errorf("Expected %s and %s to be connected", elements[i], elements[j])
			}
		}
	}

	// "e" should be separate
	if ds.Connected("a", "e") {
		t.Error("Expected 'a' and 'e' to be in different sets")
	}

	// Set count should be 2
	if ds.SetCount() != 2 {
		t.Errorf("Expected set count 2, got %d", ds.SetCount())
	}
}

func TestDisjointSetPathCompression(t *testing.T) {
	ds := New[int]()

	// Create a chain: 1-2-3-4-5
	for i := 1; i <= 5; i++ {
		ds.MakeSet(i)
	}
	ds.Union(1, 2)
	ds.Union(2, 3)
	ds.Union(3, 4)
	ds.Union(4, 5)

	// After Find(5), path compression should make 5 point directly to root
	root, ok := ds.Find(5)
	if !ok {
		t.Error("Expected Find(5) to succeed")
	}

	// Check that intermediate nodes now point to root (path compression)
	// We can't directly check internal structure, but we can verify performance
	// by ensuring subsequent operations are fast
	root2, ok2 := ds.Find(3)
	if !ok2 || root2 != root {
		t.Errorf("Expected Find(3) to return same root %d, got %d", root, root2)
	}
}

func TestDisjointSetSetsAndElements(t *testing.T) {
	ds := New[string]()

	// Add elements and create sets
	ds.MakeSet("apple")
	ds.MakeSet("banana")
	ds.MakeSet("cherry")
	ds.MakeSet("date")

	ds.Union("apple", "banana")
	ds.Union("cherry", "date")

	// Get all elements
	elements := ds.Elements()
	if len(elements) != 4 {
		t.Errorf("Expected 4 elements, got %d", len(elements))
	}

	// Get all sets
	sets := ds.Sets()
	if len(sets) != 2 {
		t.Errorf("Expected 2 sets, got %d", len(sets))
	}

	// Each set should have 2 elements
	setSizes := []int{len(sets[0]), len(sets[1])}
	if setSizes[0] != 2 || setSizes[1] != 2 {
		t.Errorf("Expected both sets to have size 2, got %v", setSizes)
	}

	// Verify set contents
	allElements := make(map[string]bool)
	for _, set := range sets {
		for _, elem := range set {
			allElements[elem] = true
		}
	}
	expected := map[string]bool{"apple": true, "banana": true, "cherry": true, "date": true}
	if !reflect.DeepEqual(allElements, expected) {
		t.Errorf("Expected elements %v, got %v", expected, allElements)
	}
}

func TestDisjointSetForEachSet(t *testing.T) {
	ds := New[int]()

	// Create sets
	ds.MakeSet(1)
	ds.MakeSet(2)
	ds.MakeSet(3)
	ds.MakeSet(4)

	ds.Union(1, 2)
	ds.Union(3, 4)

	var setCount int
	var totalElements int
	ds.ForEachSet(func(elements []int) {
		setCount++
		totalElements += len(elements)
	})

	if setCount != 2 {
		t.Errorf("Expected 2 sets in ForEachSet, got %d", setCount)
	}
	if totalElements != 4 {
		t.Errorf("Expected 4 total elements, got %d", totalElements)
	}
}

func TestDisjointSetNonExistentElements(t *testing.T) {
	ds := New[string]()

	// Find on non-existent element
	if _, ok := ds.Find("nonexistent"); ok {
		t.Error("Expected Find on non-existent element to return false")
	}

	// Connected with non-existent elements
	if ds.Connected("a", "b") {
		t.Error("Expected Connected with non-existent elements to return false")
	}

	// Union with non-existent elements should create them
	if !ds.Union("x", "y") {
		t.Error("Expected Union with non-existent elements to return true")
	}
	if ds.Size() != 2 {
		t.Errorf("Expected size 2 after Union with non-existent elements, got %d", ds.Size())
	}
	if !ds.Connected("x", "y") {
		t.Error("Expected x and y to be connected after Union")
	}
}

func TestDisjointSetClear(t *testing.T) {
	ds := New[int]()

	ds.MakeSet(1)
	ds.MakeSet(2)
	ds.Union(1, 2)

	ds.Clear()

	if !ds.IsEmpty() {
		t.Error("Expected empty disjoint set after Clear()")
	}
	if ds.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", ds.Size())
	}

	// Operations on cleared set should work normally
	ds.MakeSet(10)
	if ds.Size() != 1 {
		t.Errorf("Expected size 1 after adding to cleared set, got %d", ds.Size())
	}
}

func TestDisjointSetConcurrency(t *testing.T) {
	ds := New[int]()
	done := make(chan bool)

	// Create initial sets
	for i := 0; i < 10; i++ {
		ds.MakeSet(i)
	}

	// Perform unions concurrently
	go func() {
		for i := 0; i < 5; i++ {
			ds.Union(i, i+5)
		}
		done <- true
	}()

	// Check connections concurrently
	go func() {
		for i := 0; i < 5; i++ {
			ds.Connected(i, i+5)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state
	if ds.SetCount() != 5 {
		t.Errorf("Expected 5 sets after concurrent unions, got %d", ds.SetCount())
	}

	// Each pair should be connected
	for i := 0; i < 5; i++ {
		if !ds.Connected(i, i+5) {
			t.Errorf("Expected %d and %d to be connected", i, i+5)
		}
	}
}