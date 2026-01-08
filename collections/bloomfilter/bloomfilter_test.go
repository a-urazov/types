package bloomfilter

import (
	"fmt"
	"testing"
)

func TestBloomFilterBasicOperations(t *testing.T) {
	// Create a BloomFilter with expected 1000 elements and 1% false positive rate
	bf := New(1000, 0.01)

	// Test empty filter
	if !bf.IsEmpty() {
		t.Error("Expected empty BloomFilter")
	}
	if bf.Size() != 0 {
		t.Errorf("Expected size 0, got %d", bf.Size())
	}

	// Add some elements
	elements := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, elem := range elements {
		bf.Put([]byte(elem))
	}

	// Check that added elements are present
	for _, elem := range elements {
		if !bf.MightContain([]byte(elem)) {
			t.Errorf("Expected element %s to be present", elem)
		}
	}

	// Size should be approximately the number of elements added
	if bf.Size() < 3 || bf.Size() > 10 {
		t.Errorf("Expected size between 3-10, got %d", bf.Size())
	}

	// Clear the filter
	bf.Clear()
	if !bf.IsEmpty() {
		t.Error("Expected empty BloomFilter after Clear()")
	}
}

func TestBloomFilterFalsePositives(t *testing.T) {
	// Create a small BloomFilter to increase chance of false positives
	bf := New(10, 0.1) // 10% false positive rate

	// Add some elements
	addedElements := []string{"test1", "test2", "test3", "test4", "test5"}
	for _, elem := range addedElements {
		bf.Put([]byte(elem))
	}

	// Check that added elements are definitely present
	for _, elem := range addedElements {
		if !bf.MightContain([]byte(elem)) {
			t.Errorf("Added element %s should be present", elem)
		}
	}

	// Test some non-added elements - some might return true (false positives)
	nonAddedElements := []string{"not1", "not2", "not3", "not4", "not5"}
	falsePositives := 0
	for _, elem := range nonAddedElements {
		if bf.MightContain([]byte(elem)) {
			falsePositives++
		}
	}

	// We expect some false positives but not all
	if falsePositives == len(nonAddedElements) {
		t.Error("All non-added elements returned true - filter might be saturated")
	}

	// The false positive rate should be roughly around 10%
	expectedFalsePositiveRate := bf.FalsePositiveRate()
	if expectedFalsePositiveRate < 0.0 || expectedFalsePositiveRate > 1.0 {
		t.Errorf("False positive rate should be between 0 and 1, got %f", expectedFalsePositiveRate)
	}
}

func TestBloomFilterWithSize(t *testing.T) {
	// Create BloomFilter with explicit size
	bf := WithSize(1024, 3) // 1024 bits, 3 hash functions

	if bf.Capacity() != 1024 {
		t.Errorf("Expected capacity 1024, got %d", bf.Capacity())
	}
	if bf.NumHashes() != 3 {
		t.Errorf("Expected 3 hash functions, got %d", bf.NumHashes())
	}

	// Add and test elements
	bf.Put([]byte("hello"))
	if !bf.MightContain([]byte("hello")) {
		t.Error("Expected 'hello' to be present")
	}
	if bf.MightContain([]byte("world")) {
		// This could be a false positive, but unlikely with this configuration
		// We'll accept it as possible
		t.Log("Got false positive for 'world' - this is acceptable")
	}
}

func TestBloomFilterMerge(t *testing.T) {
	// Create two compatible BloomFilters
	bf1 := New(100, 0.01)
	bf2 := New(100, 0.01)

	// Add different elements to each
	bf1.Put([]byte("apple"))
	bf1.Put([]byte("banana"))

	bf2.Put([]byte("cherry"))
	bf2.Put([]byte("date"))

	// Merge them
	err := bf1.Merge(bf2)
	if err != nil {
		t.Errorf("Expected no error when merging compatible filters, got %v", err)
	}

	// Check that all elements are present
	if !bf1.MightContain([]byte("apple")) {
		t.Error("Expected 'apple' to be present after merge")
	}
	if !bf1.MightContain([]byte("banana")) {
		t.Error("Expected 'banana' to be present after merge")
	}
	if !bf1.MightContain([]byte("cherry")) {
		t.Error("Expected 'cherry' to be present after merge")
	}
	if !bf1.MightContain([]byte("date")) {
		t.Error("Expected 'date' to be present after merge")
	}

	// Test merging incompatible filters
	bf3 := New(50, 0.01) // Different size
	err = bf1.Merge(bf3)
	if err == nil {
		t.Error("Expected error when merging incompatible filters")
	}
	if err != ErrIncompatibleFilters {
		t.Errorf("Expected ErrIncompatibleFilters, got %v", err)
	}
}

func TestBloomFilterClone(t *testing.T) {
	bf := New(100, 0.01)
	bf.Put([]byte("test"))

	// Clone the filter
	clone := bf.Clone()

	// Both should contain the same elements
	if !bf.MightContain([]byte("test")) {
		t.Error("Original filter should contain 'test'")
	}
	if !clone.MightContain([]byte("test")) {
		t.Error("Cloned filter should contain 'test'")
	}

	// Modify original and check clone is unchanged
	bf.Put([]byte("another"))
	if clone.MightContain([]byte("another")) {
		t.Error("Clone should not be affected by changes to original")
	}
}

func TestBloomFilterEdgeCases(t *testing.T) {
	// Test with invalid parameters
	bf1 := New(0, 0.0) // Should use defaults
	if bf1.Capacity() == 0 {
		t.Error("BloomFilter with invalid params should still work")
	}

	bf2 := New(-1, 2.0) // Should use defaults
	if bf2.Capacity() == 0 {
		t.Error("BloomFilter with invalid params should still work")
	}

	// Test with empty data
	bf3 := New(10, 0.1)
	bf3.Put([]byte{})
	if !bf3.MightContain([]byte{}) {
		t.Error("Empty byte slice should be handled correctly")
	}

	// Test very large number of hash functions
	bf4 := WithSize(1024, 100) // Should cap at 32
	if bf4.NumHashes() != 32 {
		t.Errorf("Expected max 32 hash functions, got %d", bf4.NumHashes())
	}
}

func TestBloomFilterConcurrency(t *testing.T) {
	const itemFormat = "item-%d"
	bf := New(1000, 0.01)
	done := make(chan bool)

	// Add elements concurrently
	go func() {
		for i := range 100 {
			bf.Put([]byte(fmt.Sprintf(itemFormat, i)))
		}
		done <- true
	}()

	// Check elements concurrently
	go func() {
		for i := range 50 {
			bf.MightContain([]byte(fmt.Sprintf(itemFormat, i)))
		}
		done <- true
	}()

	<-done
	<-done

	// All added elements should be present
	for i := 0; i < 100; i++ {
		if !bf.MightContain([]byte(fmt.Sprintf(itemFormat, i))) {
			t.Errorf("Element item-%d should be present", i)
		}
	}
}