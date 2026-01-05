package btree

import (
	"testing"
)

func TestNewBTree(t *testing.T) {
	bt := New[int]()
	if bt == nil {
		t.Error("New() should not return nil")
	}
	if !bt.IsEmpty() {
		t.Error("New BTree should be empty")
	}
	if bt.Size() != 0 {
		t.Errorf("New BTree size should be 0, got %d", bt.Size())
	}
}

func TestInsert(t *testing.T) {
	bt := New[int]()
	values := []int{10, 5, 20, 3, 7, 15, 25, 1, 4, 6, 8, 12, 17, 22, 27}

	for _, v := range values {
		bt.Insert(v)
		if !bt.Search(v) {
			t.Errorf("BTree should contain inserted value %d", v)
		}
		// Size check is O(n), so we'll do it occasionally for performance
		if len(values) <= 10 || (len(values) > 10 && indexOfInt(values, v)%5 == 0) {
			expectedSize := countUniqueInts(values[:indexOfInt(values, v)+1])
			if bt.Size() != expectedSize {
				t.Errorf("BTree size should be %d after inserting %d values, got %d", expectedSize, indexOfInt(values, v)+1, bt.Size())
			}
		}
	}
}

func TestSearch(t *testing.T) {
	bt := New[string]()
	values := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}

	for _, v := range values {
		bt.Insert(v)
	}

	for _, v := range values {
		if !bt.Search(v) {
			t.Errorf("BTree should contain inserted value %s", v)
		}
	}

	nonExistentValues := []string{"grapefruit", "kiwi", "lemon", "mango"}
	for _, v := range nonExistentValues {
		if bt.Search(v) {
			t.Errorf("BTree should not contain non-existent value %s", v)
		}
	}
}

func TestDelete(t *testing.T) {
	bt := New[int]()
	values := []int{10, 5, 20, 3, 7, 15, 25}

	for _, v := range values {
		bt.Insert(v)
	}

	// Удалить значение
	if !bt.Delete(5) {
		t.Error("Delete should return true for existing value")
	}
	if bt.Search(5) {
		t.Error("BTree should not contain deleted value 5")
	}
	if bt.Size() != 6 {
		t.Errorf("BTree size should be 6 after deleting one value, got %d", bt.Size())
	}

	// Удалить несуществующее значение
	if bt.Delete(99) {
		t.Error("Delete should return false for non-existent value")
	}

	// Удалить все значения
	for _, v := range values {
		bt.Delete(v)
		if bt.Search(v) {
			t.Errorf("BTree should not contain deleted value %d", v)
		}
	}

	if !bt.IsEmpty() {
		t.Error("BTree should be empty after deleting all values")
	}
}

func TestSize(t *testing.T) {
	bt := New[string]()
	values := []string{"a", "bb", "ccc", "dddd"}

	if bt.Size() != 0 {
		t.Errorf("Empty BTree size should be 0, got %d", bt.Size())
	}

	for i, v := range values {
		bt.Insert(v)
		if bt.Size() != i+1 {
			t.Errorf("BTree size should be %d after inserting %d values, got %d", i+1, i+1, bt.Size())
		}
	}

	// Удалить некоторые значения
	bt.Delete("bb")
	if bt.Size() != 3 {
		t.Errorf("BTree size should be 3 after deleting 'bb', got %d", bt.Size())
	}

	bt.Delete("zzz") // Несуществующее
	if bt.Size() != 3 {
		t.Errorf("BTree size should still be 3 after trying to delete non-existent 'zzz', got %d", bt.Size())
	}

	// Удалить все
	for _, v := range values {
		bt.Delete(v)
	}
	if bt.Size() != 0 {
		t.Errorf("BTree size should be 0 after deleting all values, got %d", bt.Size())
	}
}

func TestIsEmpty(t *testing.T) {
	bt := New[float64]()
	if !bt.IsEmpty() {
		t.Error("New BTree should be empty")
	}

	bt.Insert(1.5)
	if bt.IsEmpty() {
		t.Error("BTree should not be empty after insertion")
	}

	bt.Delete(1.5)
	if !bt.IsEmpty() {
		t.Error("BTree should be empty after deleting the last value")
	}
}

// Вспомогательная функция для поиска индекса int в срезе
func indexOfInt(slice []int, item int) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// Вспомогательная функция для подсчета количества уникальных целых чисел в срезе
func countUniqueInts(values []int) int {
	seen := make(map[int]bool)
	for _, v := range values {
		seen[v] = true
	}
	return len(seen)
}
