package circularlist

import (
	"testing"
)

func TestCircularListBasicOperations(t *testing.T) {
	cl := New[int]()

	// Test empty list
	if !cl.IsEmpty() {
		t.Error("Expected empty list")
	}
	if cl.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cl.Size())
	}

	// Add elements
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)

	if cl.IsEmpty() {
		t.Error("Expected non-empty list")
	}
	if cl.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cl.Size())
	}

	// Test Get
	if val, ok := cl.Get(0); !ok || val != 1 {
		t.Errorf("Expected Get(0) to return 1, got %d", val)
	}
	if val, ok := cl.Get(1); !ok || val != 2 {
		t.Errorf("Expected Get(1) to return 2, got %d", val)
	}
	if val, ok := cl.Get(2); !ok || val != 3 {
		t.Errorf("Expected Get(2) to return 3, got %d", val)
	}

	// Test invalid Get
	if _, ok := cl.Get(3); ok {
		t.Error("Expected Get(3) to return false")
	}
	if _, ok := cl.Get(-1); ok {
		t.Error("Expected Get(-1) to return false")
	}

	// Test Contains
	if !cl.Contains(2) {
		t.Error("Expected list to contain 2")
	}
	if cl.Contains(4) {
		t.Error("Expected list to not contain 4")
	}

	// Test IndexOf
	if idx := cl.IndexOf(2); idx != 1 {
		t.Errorf("Expected IndexOf(2) to return 1, got %d", idx)
	}
	if idx := cl.IndexOf(4); idx != -1 {
		t.Errorf("Expected IndexOf(4) to return -1, got %d", idx)
	}

	// Test Set
	if !cl.Set(1, 20) {
		t.Error("Expected Set(1, 20) to return true")
	}
	if val, ok := cl.Get(1); !ok || val != 20 {
		t.Errorf("Expected Get(1) to return 20 after Set, got %d", val)
	}

	// Test invalid Set
	if cl.Set(5, 100) {
		t.Error("Expected Set(5, 100) to return false")
	}
}

func TestCircularListAddAt(t *testing.T) {
	cl := New[string]()

	// Add at beginning
	if !cl.AddAt(0, "first") {
		t.Error("Expected AddAt(0, 'first') to return true")
	}

	// Add at end
	if !cl.AddAt(1, "last") {
		t.Error("Expected AddAt(1, 'last') to return true")
	}

	// Add in middle
	if !cl.AddAt(1, "middle") {
		t.Error("Expected AddAt(1, 'middle') to return true")
	}

	if cl.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cl.Size())
	}

	// Check order
	expected := []string{"first", "middle", "last"}
	for i, exp := range expected {
		if val, ok := cl.Get(i); !ok || val != exp {
			t.Errorf("Expected Get(%d) to return %s, got %s", i, exp, val)
		}
	}

	// Test invalid AddAt
	if cl.AddAt(-1, "invalid") {
		t.Error("Expected AddAt(-1, 'invalid') to return false")
	}
	if cl.AddAt(5, "invalid") {
		t.Error("Expected AddAt(5, 'invalid') to return false")
	}
}

func TestCircularListRemove(t *testing.T) {
	cl := New[int]()

	// Add elements
	cl.Add(1)
	cl.Add(2)
	cl.Add(3)
	cl.Add(2) // duplicate

	// Remove first occurrence of 2
	if !cl.Remove(2) {
		t.Error("Expected Remove(2) to return true")
	}
	if cl.Size() != 3 {
		t.Errorf("Expected size 3 after removal, got %d", cl.Size())
	}

	// Check remaining elements
	expected := []int{1, 3, 2} // second 2 should remain
	for i, exp := range expected {
		if val, ok := cl.Get(i); !ok || val != exp {
			t.Errorf("Expected Get(%d) to return %d, got %d", i, exp, val)
		}
	}

	// Remove non-existent element
	if cl.Remove(5) {
		t.Error("Expected Remove(5) to return false")
	}

	// Remove all elements
	if !cl.Remove(1) {
		t.Error("Expected Remove(1) to return true")
	}
	if !cl.Remove(3) {
		t.Error("Expected Remove(3) to return true")
	}
	if !cl.Remove(2) {
		t.Error("Expected Remove(2) to return true")
	}

	if !cl.IsEmpty() {
		t.Error("Expected empty list after removing all elements")
	}
}

func TestCircularListRemoveAt(t *testing.T) {
	cl := New[int]()

	// Add elements
	cl.Add(10)
	cl.Add(20)
	cl.Add(30)

	// Remove middle element
	if !cl.RemoveAt(1) {
		t.Error("Expected RemoveAt(1) to return true")
	}
	if cl.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", cl.Size())
	}

	// Check remaining elements
	if val, ok := cl.Get(0); !ok || val != 10 {
		t.Errorf("Expected Get(0) to return 10, got %d", val)
	}
	if val, ok := cl.Get(1); !ok || val != 30 {
		t.Errorf("Expected Get(1) to return 30, got %d", val)
	}

	// Test invalid RemoveAt
	if cl.RemoveAt(5) {
		t.Error("Expected RemoveAt(5) to return false")
	}
	if cl.RemoveAt(-1) {
		t.Error("Expected RemoveAt(-1) to return false")
	}
}

func TestCircularListToSliceAndForEach(t *testing.T) {
	cl := New[int]()

	// Add elements
	for i := 1; i <= 5; i++ {
		cl.Add(i * 10)
	}

	// Test ToSlice
	slice := cl.ToSlice()
	expected := []int{10, 20, 30, 40, 50}
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}
	for i, exp := range expected {
		if slice[i] != exp {
			t.Errorf("Expected slice[%d] to be %d, got %d", i, exp, slice[i])
		}
	}

	// Test ForEach
	var forEachResult []int
	cl.ForEach(func(index int, value int) {
		forEachResult = append(forEachResult, value)
	})
	if len(forEachResult) != len(expected) {
		t.Errorf("Expected ForEach result length %d, got %d", len(expected), len(forEachResult))
	}
	for i, exp := range expected {
		if forEachResult[i] != exp {
			t.Errorf("Expected ForEach result[%d] to be %d, got %d", i, exp, forEachResult[i])
		}
	}

	// Test ReverseForEach
	var reverseResult []int
	cl.ReverseForEach(func(index int, value int) {
		reverseResult = append(reverseResult, value)
	})
	expectedReverse := []int{50, 40, 30, 20, 10}
	for i, exp := range expectedReverse {
		if reverseResult[i] != exp {
			t.Errorf("Expected ReverseForEach result[%d] to be %d, got %d", i, exp, reverseResult[i])
		}
	}
}

func TestCircularListCircularOperations(t *testing.T) {
	cl := New[int]()

	// Add elements
	for i := 1; i <= 5; i++ {
		cl.Add(i)
	}

	// Test GetNext
	if next, ok := cl.GetNext(1); !ok || next != 2 {
		t.Errorf("Expected GetNext(1) to return 2, got %d", next)
	}
	if next, ok := cl.GetNext(5); !ok || next != 1 {
		t.Errorf("Expected GetNext(5) to return 1 (circular), got %d", next)
	}

	// Test GetPrev
	if prev, ok := cl.GetPrev(2); !ok || prev != 1 {
		t.Errorf("Expected GetPrev(2) to return 1, got %d", prev)
	}
	if prev, ok := cl.GetPrev(1); !ok || prev != 5 {
		t.Errorf("Expected GetPrev(1) to return 5 (circular), got %d", prev)
	}

	// Test non-existent element
	if _, ok := cl.GetNext(10); ok {
		t.Error("Expected GetNext(10) to return false")
	}
	if _, ok := cl.GetPrev(10); ok {
		t.Error("Expected GetPrev(10) to return false")
	}
}

func TestCircularListRotation(t *testing.T) {
	cl := New[int]()

	// Add elements
	for i := 1; i <= 5; i++ {
		cl.Add(i)
	}

	// Rotate left by 2
	cl.RotateLeft(2)

	// After rotating left by 2: [3, 4, 5, 1, 2]
	expected := []int{3, 4, 5, 1, 2}
	for i, exp := range expected {
		if val, ok := cl.Get(i); !ok || val != exp {
			t.Errorf("Expected Get(%d) to return %d after left rotation, got %d", i, exp, val)
		}
	}

	// Rotate right by 3 (should be equivalent to left rotation by -3)
	cl.RotateRight(3)

	// After rotating right by 3: [1, 2, 3, 4, 5]
	expected = []int{1, 2, 3, 4, 5}
	for i, exp := range expected {
		if val, ok := cl.Get(i); !ok || val != exp {
			t.Errorf("Expected Get(%d) to return %d after right rotation, got %d", i, exp, val)
		}
	}

	// Test rotation by size (should be no change)
	cl.RotateLeft(5)
	for i, exp := range expected {
		if val, ok := cl.Get(i); !ok || val != exp {
			t.Errorf("Expected Get(%d) to return %d after rotation by size, got %d", i, exp, val)
		}
	}
}

func TestCircularListClear(t *testing.T) {
	cl := New[string]()

	cl.Add("a")
	cl.Add("b")
	cl.Add("c")

	cl.Clear()

	if !cl.IsEmpty() {
		t.Error("Expected empty list after Clear()")
	}
	if cl.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", cl.Size())
	}
	if _, ok := cl.Get(0); ok {
		t.Error("Expected Get(0) to return false after Clear()")
	}
}

func TestCircularListSingleElement(t *testing.T) {
	cl := New[int]()

	// Add single element
	cl.Add(42)

	if cl.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cl.Size())
	}

	// Test Get
	if val, ok := cl.Get(0); !ok || val != 42 {
		t.Errorf("Expected Get(0) to return 42, got %d", val)
	}

	// Test circular properties
	if next, ok := cl.GetNext(42); !ok || next != 42 {
		t.Errorf("Expected GetNext(42) to return 42 in single element list, got %d", next)
	}
	if prev, ok := cl.GetPrev(42); !ok || prev != 42 {
		t.Errorf("Expected GetPrev(42) to return 42 in single element list, got %d", prev)
	}

	// Remove the single element
	if !cl.Remove(42) {
		t.Error("Expected Remove(42) to return true")
	}
	if !cl.IsEmpty() {
		t.Error("Expected empty list after removing single element")
	}
}

func TestCircularListConcurrency(t *testing.T) {
	cl := New[int]()
	done := make(chan bool)

	// Writers
	go func() {
		for i := 0; i < 50; i++ {
			cl.Add(i)
		}
		done <- true
	}()

	// Readers
	go func() {
		for i := 0; i < 25; i++ {
			cl.Size()
			if cl.Size() > 0 {
				cl.Get(0)
				cl.Contains(i)
			}
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state
	if cl.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent operations, got %d", cl.Size())
	}
}