package fenwicktree

import (
	"testing"
)

func TestFenwickTreeBasicOperations(t *testing.T) {
	ft := New(10)

	// Test empty tree
	if !ft.IsEmpty() {
		t.Error("Expected empty tree")
	}
	if ft.Size() != 10 {
		t.Errorf("Expected size 10, got %d", ft.Size())
	}
	if ft.Query(9) != 0 {
		t.Errorf("Expected sum 0 for empty tree, got %d", ft.Query(9))
	}

	// Add values
	ft.Add(0, 1)
	ft.Add(1, 2)
	ft.Add(3, 3)
	ft.Add(5, 4)
	ft.Add(9, 5)

	// Test Query (prefix sums)
	if sum := ft.Query(0); sum != 1 {
		t.Errorf("Expected Query(0) to be 1, got %d", sum)
	}
	if sum := ft.Query(1); sum != 3 { // 1+2
		t.Errorf("Expected Query(1) to be 3, got %d", sum)
	}
	if sum := ft.Query(2); sum != 3 { // 1+2+0
		t.Errorf("Expected Query(2) to be 3, got %d", sum)
	}
	if sum := ft.Query(3); sum != 6 { // 1+2+0+3
		t.Errorf("Expected Query(3) to be 6, got %d", sum)
	}
	if sum := ft.Query(9); sum != 15 { // 1+2+0+3+0+4+0+0+0+5
		t.Errorf("Expected Query(9) to be 15, got %d", sum)
	}

	// Test QueryRange
	if sum := ft.QueryRange(0, 1); sum != 3 {
		t.Errorf("Expected QueryRange(0,1) to be 3, got %d", sum)
	}
	if sum := ft.QueryRange(2, 4); sum != 3 { // 0+3+0
		t.Errorf("Expected QueryRange(2,4) to be 3, got %d", sum)
	}
	if sum := ft.QueryRange(5, 9); sum != 9 { // 4+0+0+0+5
		t.Errorf("Expected QueryRange(5,9) to be 9, got %d", sum)
	}

	// Test invalid ranges
	if sum := ft.QueryRange(5, 3); sum != 0 {
		t.Errorf("Expected QueryRange(5,3) to be 0, got %d", sum)
	}
	if sum := ft.QueryRange(-1, 2); sum != 3 {
		t.Errorf("Expected QueryRange(-1,2) to be 3, got %d", sum)
	}
	if sum := ft.QueryRange(8, 12); sum != 5 {
		t.Errorf("Expected QueryRange(8,12) to be 5, got %d", sum)
	}
}

func TestFenwickTreeFromSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	ft := FromSlice(data)

	if ft.Size() != 5 {
		t.Errorf("Expected size 5, got %d", ft.Size())
	}

	// Test prefix sums
	if sum := ft.Query(0); sum != 1 {
		t.Errorf("Expected Query(0) to be 1, got %d", sum)
	}
	if sum := ft.Query(2); sum != 6 { // 1+2+3
		t.Errorf("Expected Query(2) to be 6, got %d", sum)
	}
	if sum := ft.Query(4); sum != 15 { // 1+2+3+4+5
		t.Errorf("Expected Query(4) to be 15, got %d", sum)
	}

	// Test range sums
	if sum := ft.QueryRange(1, 3); sum != 9 { // 2+3+4
		t.Errorf("Expected QueryRange(1,3) to be 9, got %d", sum)
	}
}

func TestFenwickTreeGetAndSet(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	ft := FromSlice(data)

	// Test Get
	for i, expected := range data {
		if val := ft.Get(i); val != expected {
			t.Errorf("Expected Get(%d) to be %d, got %d", i, expected, val)
		}
	}

	// Test Set
	ft.Set(2, 35) // Change 30 to 35
	if val := ft.Get(2); val != 35 {
		t.Errorf("Expected Get(2) to be 35 after Set, got %d", val)
	}

	// Check that sums are updated correctly
	if sum := ft.Query(2); sum != 65 { // 10+20+35
		t.Errorf("Expected Query(2) to be 65 after Set, got %d", sum)
	}
	if sum := ft.Query(4); sum != 155 { // 10+20+35+40+50
		t.Errorf("Expected Query(4) to be 155 after Set, got %d", sum)
	}

	// Test invalid Get/Set
	if val := ft.Get(10); val != 0 {
		t.Errorf("Expected Get(10) to be 0, got %d", val)
	}
	ft.Set(-1, 100) // Should do nothing
	if ft.Query(4) != 155 {
		t.Error("Set with invalid index should not change sums")
	}
}

func TestFenwickTreeToSliceAndForEach(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	ft := FromSlice(data)

	// Test ToSlice
	slice := ft.ToSlice()
	if len(slice) != len(data) {
		t.Errorf("Expected slice length %d, got %d", len(data), len(slice))
	}
	for i, expected := range data {
		if slice[i] != expected {
			t.Errorf("Expected slice[%d] to be %d, got %d", i, expected, slice[i])
		}
	}

	// Test ForEach
	var forEachResult []int
	ft.ForEach(func(index, value int) {
		forEachResult = append(forEachResult, value)
	})
	if len(forEachResult) != len(data) {
		t.Errorf("Expected ForEach result length %d, got %d", len(data), len(forEachResult))
	}
	for i, expected := range data {
		if forEachResult[i] != expected {
			t.Errorf("Expected ForEach result[%d] to be %d, got %d", i, expected, forEachResult[i])
		}
	}
}

func TestFenwickTreeClear(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	ft := FromSlice(data)

	ft.Clear()

	if !ft.IsEmpty() {
		t.Error("Expected empty tree after Clear()")
	}
	if ft.Query(4) != 0 {
		t.Errorf("Expected sum 0 after Clear(), got %d", ft.Query(4))
	}
}

func TestFenwickTreeResize(t *testing.T) {
	data := []int{1, 2, 3}
	ft := FromSlice(data)

	// Resize to larger size
	ft.Resize(5)
	if ft.Size() != 5 {
		t.Errorf("Expected size 5 after resize, got %d", ft.Size())
	}
	if sum := ft.Query(2); sum != 6 { // Original data preserved
		t.Errorf("Expected sum 6 for original range, got %d", sum)
	}
	if sum := ft.Query(4); sum != 6 { // New elements are zero
		t.Errorf("Expected sum 6 for new range, got %d", sum)
	}

	// Resize to smaller size
	ft.Resize(2)
	if ft.Size() != 2 {
		t.Errorf("Expected size 2 after resize, got %d", ft.Size())
	}
	if sum := ft.Query(1); sum != 3 { // First 2 elements preserved
		t.Errorf("Expected sum 3 for resized range, got %d", sum)
	}

	// Resize to zero
	ft.Resize(0)
	if !ft.IsEmpty() {
		t.Error("Expected empty tree after resize to 0")
	}
}

func TestFenwickTreeConcurrency(t *testing.T) {
	ft := New(100)
	done := make(chan bool)

	// Adders
	go func() {
		for i := 0; i < 100; i++ {
			ft.Add(i, 1)
		}
		done <- true
	}()

	// Query workers
	go func() {
		for i := 0; i < 50; i++ {
			ft.Query(i)
			ft.QueryRange(i, i+10)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state
	if ft.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent operations, got %d", ft.Size())
	}
	if sum := ft.Query(99); sum != 100 {
		t.Errorf("Expected total sum 100, got %d", sum)
	}
}
