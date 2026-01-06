package segmenttree

import (
	"testing"
)

func TestSegmentTreeSum(t *testing.T) {
	data := []int{1, 3, 5, 7, 9, 11}
	st := Sum(data)

	// Test basic queries
	if sum, ok := st.Query(0, 5); !ok || sum != 36 {
		t.Errorf("Expected sum 36 for range [0,5], got %d", sum)
	}
	if sum, ok := st.Query(0, 2); !ok || sum != 9 {
		t.Errorf("Expected sum 9 for range [0,2], got %d", sum)
	}
	if sum, ok := st.Query(3, 5); !ok || sum != 27 {
		t.Errorf("Expected sum 27 for range [3,5], got %d", sum)
	}
	if sum, ok := st.Query(2, 4); !ok || sum != 21 {
		t.Errorf("Expected sum 21 for range [2,4], got %d", sum)
	}

	// Test single element queries
	for i, expected := range data {
		if val, ok := st.Query(i, i); !ok || val != expected {
			t.Errorf("Expected %d for range [%d,%d], got %d", expected, i, i, val)
		}
	}

	// Test updates
	if !st.Update(2, 10) {
		t.Error("Expected Update(2, 10) to succeed")
	}
	if sum, ok := st.Query(0, 2); !ok || sum != 14 { // 1 + 3 + 10 = 14
		t.Errorf("Expected sum 14 after update, got %d", sum)
	}
	if sum, ok := st.Query(0, 5); !ok || sum != 41 { // 36 - 5 + 10 = 41
		t.Errorf("Expected total sum 41 after update, got %d", sum)
	}

	// Test invalid ranges
	if _, ok := st.Query(-1, 2); ok {
		t.Error("Expected Query(-1, 2) to fail")
	}
	if _, ok := st.Query(0, 10); ok {
		t.Error("Expected Query(0, 10) to fail")
	}
	if _, ok := st.Query(3, 1); ok {
		t.Error("Expected Query(3, 1) to fail")
	}

	// Test invalid updates
	if st.Update(10, 5) {
		t.Error("Expected Update(10, 5) to fail")
	}
	if st.Update(-1, 5) {
		t.Error("Expected Update(-1, 5) to fail")
	}
}

func TestSegmentTreeMin(t *testing.T) {
	data := []int{5, 2, 8, 1, 9, 3}
	st := Min(data)

	// Test basic queries
	if min, ok := st.Query(0, 5); !ok || min != 1 {
		t.Errorf("Expected min 1 for range [0,5], got %d", min)
	}
	if min, ok := st.Query(0, 2); !ok || min != 2 {
		t.Errorf("Expected min 2 for range [0,2], got %d", min)
	}
	if min, ok := st.Query(3, 5); !ok || min != 1 {
		t.Errorf("Expected min 1 for range [3,5], got %d", min)
	}
	if min, ok := st.Query(1, 4); !ok || min != 1 {
		t.Errorf("Expected min 1 for range [1,4], got %d", min)
	}

	// Test updates
	if !st.Update(3, 10) {
		t.Error("Expected Update(3, 10) to succeed")
	}
	if min, ok := st.Query(0, 5); !ok || min != 2 {
		t.Errorf("Expected min 2 after update, got %d", min)
	}
	if min, ok := st.Query(3, 5); !ok || min != 3 {
		t.Errorf("Expected min 3 for range [3,5] after update, got %d", min)
	}
}

func TestSegmentTreeMax(t *testing.T) {
	data := []int{5, 2, 8, 1, 9, 3}
	st := Max(data)

	// Test basic queries
	if max, ok := st.Query(0, 5); !ok || max != 9 {
		t.Errorf("Expected max 9 for range [0,5], got %d", max)
	}
	if max, ok := st.Query(0, 2); !ok || max != 8 {
		t.Errorf("Expected max 8 for range [0,2], got %d", max)
	}
	if max, ok := st.Query(3, 5); !ok || max != 9 {
		t.Errorf("Expected max 9 for range [3,5], got %d", max)
	}
	if max, ok := st.Query(1, 4); !ok || max != 9 {
		t.Errorf("Expected max 9 for range [1,4], got %d", max)
	}

	// Test updates
	if !st.Update(4, 4) {
		t.Error("Expected Update(4, 4) to succeed")
	}
	if max, ok := st.Query(0, 5); !ok || max != 8 {
		t.Errorf("Expected max 8 after update, got %d", max)
	}
	if max, ok := st.Query(3, 5); !ok || max != 4 {
		t.Errorf("Expected max 4 for range [3,5] after update, got %d", max)
	}
}

func TestSegmentTreeCustomMergeFunction(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	// Custom merge function: product
	productFn := func(a, b float64) float64 {
		return a * b
	}

	st := New(data, productFn)

	// Test product queries
	if prod, ok := st.Query(0, 2); !ok || prod != 6.0 { // 1 * 2 * 3 = 6
		t.Errorf("Expected product 6.0 for range [0,2], got %f", prod)
	}
	if prod, ok := st.Query(1, 3); !ok || prod != 24.0 { // 2 * 3 * 4 = 24
		t.Errorf("Expected product 24.0 for range [1,3], got %f", prod)
	}

	// Test update
	if !st.Update(1, 10.0) {
		t.Error("Expected Update(1, 10.0) to succeed")
	}
	if prod, ok := st.Query(0, 2); !ok || prod != 30.0 { // 1 * 10 * 3 = 30
		t.Errorf("Expected product 30.0 after update, got %f", prod)
	}
}

func TestSegmentTreeEdgeCases(t *testing.T) {
	// Empty tree
	emptySt := Sum([]int{})
	if !emptySt.IsEmpty() {
		t.Error("Expected empty tree to be empty")
	}
	if emptySt.Size() != 0 {
		t.Errorf("Expected size 0 for empty tree, got %d", emptySt.Size())
	}
	if _, ok := emptySt.Query(0, 0); ok {
		t.Error("Expected Query on empty tree to fail")
	}
	if emptySt.Update(0, 1) {
		t.Error("Expected Update on empty tree to fail")
	}

	// Single element tree
	singleSt := Sum([]int{42})
	if singleSt.IsEmpty() {
		t.Error("Expected single element tree to not be empty")
	}
	if singleSt.Size() != 1 {
		t.Errorf("Expected size 1 for single element tree, got %d", singleSt.Size())
	}
	if val, ok := singleSt.Query(0, 0); !ok || val != 42 {
		t.Errorf("Expected value 42 for single element tree, got %d", val)
	}
	if !singleSt.Update(0, 100) {
		t.Error("Expected Update on single element tree to succeed")
	}
	if val, ok := singleSt.Query(0, 0); !ok || val != 100 {
		t.Errorf("Expected updated value 100, got %d", val)
	}

	// Two element tree
	twoSt := Min([]int{10, 5})
	if min, ok := twoSt.Query(0, 1); !ok || min != 5 {
		t.Errorf("Expected min 5 for two element tree, got %d", min)
	}
}

func TestSegmentTreeGetAndToSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	st := Sum(data)

	// Test Get
	for i, expected := range data {
		if val, ok := st.Get(i); !ok || val != expected {
			t.Errorf("Expected Get(%d) to return %d, got %d", i, expected, val)
		}
	}

	// Test invalid Get
	if _, ok := st.Get(10); ok {
		t.Error("Expected Get(10) to fail")
	}
	if _, ok := st.Get(-1); ok {
		t.Error("Expected Get(-1) to fail")
	}

	// Test ToSlice
	slice := st.ToSlice()
	if len(slice) != len(data) {
		t.Errorf("Expected slice length %d, got %d", len(data), len(slice))
	}
	for i, expected := range data {
		if slice[i] != expected {
			t.Errorf("Expected slice[%d] to be %d, got %d", i, expected, slice[i])
		}
	}

	// Modify original and check slice is unchanged
	st.Update(0, 100)
	if slice[0] != 1 {
		t.Error("Expected original slice to be unchanged after tree update")
	}
}

func TestSegmentTreeResize(t *testing.T) {
	data := []int{1, 2, 3}
	st := Sum(data)

	// Resize to larger size
	st.Resize(5)
	if st.Size() != 5 {
		t.Errorf("Expected size 5 after resize, got %d", st.Size())
	}
	if sum, ok := st.Query(0, 2); !ok || sum != 6 {
		t.Errorf("Expected sum 6 for original range, got %d", sum)
	}
	if sum, ok := st.Query(3, 4); !ok || sum != 0 {
		t.Errorf("Expected sum 0 for new range, got %d", sum)
	}

	// Resize to smaller size
	st.Resize(2)
	if st.Size() != 2 {
		t.Errorf("Expected size 2 after resize, got %d", st.Size())
	}
	if sum, ok := st.Query(0, 1); !ok || sum != 3 {
		t.Errorf("Expected sum 3 for resized range, got %d", sum)
	}

	// Resize to zero
	st.Resize(0)
	if !st.IsEmpty() {
		t.Error("Expected empty tree after resize to 0")
	}
}

func TestSegmentTreeClear(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	st := Sum(data)

	st.Clear()

	if !st.IsEmpty() {
		t.Error("Expected empty tree after Clear()")
	}
	if st.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", st.Size())
	}
	if _, ok := st.Query(0, 0); ok {
		t.Error("Expected Query to fail on cleared tree")
	}
}

func TestSegmentTreeConcurrency(t *testing.T) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i + 1
	}
	st := Sum(data)
	done := make(chan bool)

	// Updaters
	go func() {
		for i := 0; i < 50; i++ {
			st.Update(i, i*2)
		}
		done <- true
	}()

	// Query workers
	go func() {
		for i := 0; i < 25; i++ {
			st.Query(i, i+10)
			st.Get(i)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final state is consistent
	if st.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent operations, got %d", st.Size())
	}
}
