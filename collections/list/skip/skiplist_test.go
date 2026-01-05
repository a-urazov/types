package skip

import (
	"reflect"
	"testing"
)

func TestNewSkipList(t *testing.T) {
	sl := New[int]()
	if sl == nil {
		t.Error("New() should not return nil")
	}
	if !sl.IsEmpty() {
		t.Error("New SkipList should be empty")
	}
	if sl.Size() != 0 {
		t.Errorf("New SkipList size should be 0, got %d", sl.Size())
	}
}

func TestInsert(t *testing.T) {
	sl := New[string]()
	values := []string{"apple", "banana", "cherry", "date", "elderberry"}

	for _, v := range values {
		sl.Insert(v)
		if !sl.Search(v) {
			t.Errorf("SkipList should contain inserted value %s", v)
		}
		if sl.Size() != countInsertedValues(values[:indexOf(values, v)+1]) {
			t.Errorf("SkipList size should be %d after inserting %d values, got %d", countInsertedValues(values[:indexOf(values, v)+1]), indexOf(values, v)+1, sl.Size())
		}
	}
}

func TestSearch(t *testing.T) {
	sl := New[int]()
	values := []int{10, 5, 15, 3, 7, 12, 18}

	for _, v := range values {
		sl.Insert(v)
	}

	for _, v := range values {
		if !sl.Search(v) {
			t.Errorf("SkipList should contain inserted value %d", v)
		}
	}

	nonExistentValues := []int{1, 6, 11, 20}
	for _, v := range nonExistentValues {
		if sl.Search(v) {
			t.Errorf("SkipList should not contain non-existent value %d", v)
		}
	}
}

func TestDelete(t *testing.T) {
	sl := New[int]()
	values := []int{10, 5, 15, 3, 7, 12, 18}

	for _, v := range values {
		sl.Insert(v)
	}

	// Delete a value
	if !sl.Delete(5) {
		t.Error("Delete should return true for existing value")
	}
	if sl.Search(5) {
		t.Error("SkipList should not contain deleted value 5")
	}
	if sl.Size() != 6 {
		t.Errorf("SkipList size should be 6 after deleting one value, got %d", sl.Size())
	}

	// Delete a non-existent value
	if sl.Delete(99) {
		t.Error("Delete should return false for non-existent value")
	}

	// Delete all values
	for _, v := range values {
		sl.Delete(v)
		if sl.Search(v) {
			t.Errorf("SkipList should not contain deleted value %d", v)
		}
	}

	if !sl.IsEmpty() {
		t.Error("SkipList should be empty after deleting all values")
	}
}

func TestSize(t *testing.T) {
	sl := New[string]()
	values := []string{"a", "bb", "ccc", "dddd"}

	if sl.Size() != 0 {
		t.Errorf("Empty SkipList size should be 0, got %d", sl.Size())
	}

	for i, v := range values {
		sl.Insert(v)
		if sl.Size() != i+1 {
			t.Errorf("SkipList size should be %d after inserting %d values, got %d", i+1, i+1, sl.Size())
		}
	}

	// Delete some values
	sl.Delete("bb")
	if sl.Size() != 3 {
		t.Errorf("SkipList size should be 3 after deleting 'bb', got %d", sl.Size())
	}

	sl.Delete("cccc") // Non-existent
	if sl.Size() != 3 {
		t.Errorf("SkipList size should still be 3 after trying to delete non-existent 'cccc', got %d", sl.Size())
	}

	// Delete all
	for _, v := range values {
		sl.Delete(v)
	}
	if sl.Size() != 0 {
		t.Errorf("SkipList size should be 0 after deleting all values, got %d", sl.Size())
	}
}

func TestIsEmpty(t *testing.T) {
	sl := New[float64]()
	if !sl.IsEmpty() {
		t.Error("New SkipList should be empty")
	}

	sl.Insert(1.5)
	if sl.IsEmpty() {
		t.Error("SkipList should not be empty after insertion")
	}

	sl.Delete(1.5)
	if !sl.IsEmpty() {
		t.Error("SkipList should be empty after deleting the last value")
	}
}

// Вспомогательная функция для поиска индекса значения в срезе
func indexOf(slice any, item any) int {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("Slice argument is not a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(s.Index(i).Interface(), item) {
			return i
		}
	}
	return -1
}

// Вспомогательная функция для подсчета количества уникальных вставленных значений
func countInsertedValues(values []string) int {
	count := 0
	seen := make(map[string]bool)
	for _, word := range values {
		if !seen[word] {
			seen[word] = true
			count++
		}
	}
	return count
}
