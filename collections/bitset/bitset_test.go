package bitset

import (
	"testing"
)

func TestBitSetBasicOperations(t *testing.T) {
	bs := New()

	// Test empty set
	if !bs.IsEmpty() {
		t.Error("Expected empty BitSet")
	}
	if bs.Size() != 0 {
		t.Errorf("Expected size 0, got %d", bs.Size())
	}

	// Test setting values
	bs.Set(0)
	bs.Set(1)
	bs.Set(63)
	bs.Set(64)
	bs.Set(100)

	if bs.IsEmpty() {
		t.Error("Expected non-empty BitSet")
	}
	if bs.Size() != 5 {
		t.Errorf("Expected size 5, got %d", bs.Size())
	}

	// Test getting values
	if !bs.Get(0) || !bs.Get(1) || !bs.Get(63) || !bs.Get(64) || !bs.Get(100) {
		t.Error("Expected all set values to be present")
	}
	if bs.Get(2) || bs.Get(65) {
		t.Error("Expected unset values to be absent")
	}

	// Test clearing values
	bs.Clear(1)
	bs.Clear(64)

	if bs.Get(1) || bs.Get(64) {
		t.Error("Expected cleared values to be absent")
	}
	if bs.Size() != 3 {
		t.Errorf("Expected size 3 after clearing, got %d", bs.Size())
	}

	// Test clearing non-existent values
	bs.Clear(999)
	if bs.Size() != 3 {
		t.Errorf("Expected size to remain 3, got %d", bs.Size())
	}
}

func TestBitSetWithCapacity(t *testing.T) {
	bs := WithCapacity(200)

	// Should be able to set values up to 200 without reallocation issues
	for i := 0; i <= 200; i += 10 {
		bs.Set(i)
	}

	if bs.Size() != 21 { // 0, 10, 20, ..., 200 = 21 values
		t.Errorf("Expected size 21, got %d", bs.Size())
	}

	for i := 0; i <= 200; i += 10 {
		if !bs.Get(i) {
			t.Errorf("Expected value %d to be set", i)
		}
	}
}

func TestBitSetMinAndMax(t *testing.T) {
	bs := New()

	// Test empty set
	if min := bs.Min(); min != -1 {
		t.Errorf("Expected Min() to return -1 for empty set, got %d", min)
	}
	if max := bs.Max(); max != -1 {
		t.Errorf("Expected Max() to return -1 for empty set, got %d", max)
	}

	// Add some values
	values := []int{100, 50, 200, 10, 150}
	for _, v := range values {
		bs.Set(v)
	}

	if min := bs.Min(); min != 10 {
		t.Errorf("Expected Min() to return 10, got %d", min)
	}
	if max := bs.Max(); max != 200 {
		t.Errorf("Expected Max() to return 200, got %d", max)
	}
}

func TestBitSetToSliceAndForEach(t *testing.T) {
	bs := New()
	expected := []int{5, 10, 15, 20, 25, 100, 200}
	for _, v := range expected {
		bs.Set(v)
	}

	// Test ToSlice
	slice := bs.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}

	// Verify all expected values are in the slice
	valueMap := make(map[int]bool)
	for _, v := range slice {
		valueMap[v] = true
	}
	for _, v := range expected {
		if !valueMap[v] {
			t.Errorf("Expected value %d in slice", v)
		}
	}

	// Test ForEach
	var forEachResult []int
	bs.ForEach(func(value int) {
		forEachResult = append(forEachResult, value)
	})

	if len(forEachResult) != len(expected) {
		t.Errorf("Expected ForEach result length %d, got %d", len(expected), len(forEachResult))
	}

	// Verify all expected values are in ForEach result
	forEachMap := make(map[int]bool)
	for _, v := range forEachResult {
		forEachMap[v] = true
	}
	for _, v := range expected {
		if !forEachMap[v] {
			t.Errorf("Expected value %d in ForEach result", v)
		}
	}
}

func TestBitSetNextSetBit(t *testing.T) {
	bs := New()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		bs.Set(v)
	}

	// Test starting from 0
	if next := bs.NextSetBit(0); next != 10 {
		t.Errorf("Expected NextSetBit(0) to return 10, got %d", next)
	}

	// Test starting from just before a set bit
	if next := bs.NextSetBit(19); next != 20 {
		t.Errorf("Expected NextSetBit(19) to return 20, got %d", next)
	}

	// Test starting from a set bit
	if next := bs.NextSetBit(20); next != 20 {
		t.Errorf("Expected NextSetBit(20) to return 20, got %d", next)
	}

	// Test starting after the last set bit
	if next := bs.NextSetBit(60); next != -1 {
		t.Errorf("Expected NextSetBit(60) to return -1, got %d", next)
	}

	// Test empty set
	emptyBs := New()
	if next := emptyBs.NextSetBit(0); next != -1 {
		t.Errorf("Expected NextSetBit(0) on empty set to return -1, got %d", next)
	}
}

func TestBitSetCloneAndEquals(t *testing.T) {
	bs1 := New()
	values := []int{1, 5, 10, 15, 20}
	for _, v := range values {
		bs1.Set(v)
	}

	// Test Clone
	bs2 := bs1.Clone()
	if !bs1.Equals(bs2) {
		t.Error("Cloned BitSet should be equal to original")
	}

	// Modify clone and check they're no longer equal
	bs2.Set(25)
	if bs1.Equals(bs2) {
		t.Error("Modified clone should not be equal to original")
	}

	// Test Equals with different sizes
	bs3 := New()
	bs3.Set(1)
	bs3.Set(5)
	if bs1.Equals(bs3) {
		t.Error("BitSets with different content should not be equal")
	}

	// Test Equals with empty sets
	empty1 := New()
	empty2 := New()
	if !empty1.Equals(empty2) {
		t.Error("Empty BitSets should be equal")
	}
}

func TestBitSetUnion(t *testing.T) {
	bs1 := New()
	bs1.Set(1)
	bs1.Set(3)
	bs1.Set(5)

	bs2 := New()
	bs2.Set(2)
	bs2.Set(4)
	bs2.Set(6)

	// Perform union
	bs1.Union(bs2)

	// Check that all values from both sets are present
	expectedValues := []int{1, 2, 3, 4, 5, 6}
	for _, v := range expectedValues {
		if !bs1.Get(v) {
			t.Errorf("Expected value %d to be present after union", v)
		}
	}

	// Check size
	if bs1.Size() != 6 {
		t.Errorf("Expected size 6 after union, got %d", bs1.Size())
	}

	// Original bs2 should be unchanged
	if bs2.Size() != 3 {
		t.Errorf("Original BitSet should be unchanged, expected size 3, got %d", bs2.Size())
	}
}

func TestBitSetIntersection(t *testing.T) {
	bs1 := New()
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)
	bs1.Set(4)

	bs2 := New()
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)
	bs2.Set(6)

	// Perform intersection
	bs1.Intersection(bs2)

	// Only common values should remain
	if bs1.Size() != 2 {
		t.Errorf("Expected size 2 after intersection, got %d", bs1.Size())
	}
	if !bs1.Get(3) || !bs1.Get(4) {
		t.Error("Expected common values 3 and 4 to remain")
	}
	if bs1.Get(1) || bs1.Get(2) {
		t.Error("Expected non-common values 1 and 2 to be removed")
	}
}

func TestBitSetDifference(t *testing.T) {
	bs1 := New()
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)
	bs1.Set(4)

	bs2 := New()
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)

	// Perform difference (bs1 - bs2)
	bs1.Difference(bs2)

	// Values in bs1 but not in bs2 should remain
	if bs1.Size() != 2 {
		t.Errorf("Expected size 2 after difference, got %d", bs1.Size())
	}
	if !bs1.Get(1) || !bs1.Get(2) {
		t.Error("Expected values 1 and 2 to remain")
	}
	if bs1.Get(3) || bs1.Get(4) {
		t.Error("Expected values 3 and 4 to be removed")
	}
}

func TestBitSetSymmetricDifference(t *testing.T) {
	bs1 := New()
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)

	bs2 := New()
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)

	// Perform symmetric difference (values in either set but not both)
	bs1.SymmetricDifference(bs2)

	// Expected: 1, 2, 4, 5 (3 is in both, so excluded)
	if bs1.Size() != 4 {
		t.Errorf("Expected size 4 after symmetric difference, got %d", bs1.Size())
	}

	expected := []int{1, 2, 4, 5}
	for _, v := range expected {
		if !bs1.Get(v) {
			t.Errorf("Expected value %d to be present", v)
		}
	}
	if bs1.Get(3) {
		t.Error("Expected value 3 to be absent (was in both sets)")
	}
}

func TestBitSetNegativeValues(t *testing.T) {
	bs := New()

	// Negative values should be ignored
	bs.Set(-1)
	bs.Set(-100)

	if !bs.IsEmpty() {
		t.Error("BitSet should ignore negative values and remain empty")
	}

	if bs.Get(-1) {
		t.Error("Get(-1) should return false")
	}

	bs.Clear(-5)
	// Should not panic or cause issues
}
