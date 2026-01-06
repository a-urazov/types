package multimap

import (
	"reflect"
	"testing"
)

func TestMultiMapBasicOperations(t *testing.T) {
	mm := New[string, int]()

	// Test empty multimap
	if !mm.IsEmpty() {
		t.Error("Expected empty multimap")
	}
	if mm.Size() != 0 {
		t.Errorf("Expected size 0, got %d", mm.Size())
	}
	if mm.KeySize() != 0 {
		t.Errorf("Expected key size 0, got %d", mm.KeySize())
	}

	// Test Put
	mm.Put("fruits", 1)
	mm.Put("fruits", 2)
	mm.Put("vegetables", 3)
	mm.Put("fruits", 4)

	if mm.IsEmpty() {
		t.Error("Expected non-empty multimap")
	}
	if mm.Size() != 4 {
		t.Errorf("Expected size 4, got %d", mm.Size())
	}
	if mm.KeySize() != 2 {
		t.Errorf("Expected key size 2, got %d", mm.KeySize())
	}

	// Test Get
	fruitValues := mm.Get("fruits")
	expectedFruits := []int{1, 2, 4}
	if !reflect.DeepEqual(fruitValues, expectedFruits) {
		t.Errorf("Expected fruits %v, got %v", expectedFruits, fruitValues)
	}

	vegValues := mm.Get("vegetables")
	expectedVeg := []int{3}
	if !reflect.DeepEqual(vegValues, expectedVeg) {
		t.Errorf("Expected vegetables %v, got %v", expectedVeg, vegValues)
	}

	// Test GetFirst
	if val, ok := mm.GetFirst("fruits"); !ok || val != 1 {
		t.Errorf("Expected first fruit to be 1, got %d", val)
	}

	// Test GetLast
	if val, ok := mm.GetLast("fruits"); !ok || val != 4 {
		t.Errorf("Expected last fruit to be 4, got %d", val)
	}

	// Test non-existent key
	nonExistent := mm.Get("nonexistent")
	if len(nonExistent) != 0 {
		t.Errorf("Expected empty slice for non-existent key, got %v", nonExistent)
	}
}

func TestMultiMapContains(t *testing.T) {
	mm := New[int, string]()

	mm.Put(1, "apple")
	mm.Put(1, "banana")
	mm.Put(2, "carrot")

	// Test ContainsKey
	if !mm.ContainsKey(1) || !mm.ContainsKey(2) {
		t.Error("Expected keys 1 and 2 to exist")
	}
	if mm.ContainsKey(3) {
		t.Error("Expected key 3 to not exist")
	}

	// Test ContainsValue
	if !mm.ContainsValue("apple") || !mm.ContainsValue("carrot") {
		t.Error("Expected values 'apple' and 'carrot' to exist")
	}
	if mm.ContainsValue("orange") {
		t.Error("Expected value 'orange' to not exist")
	}

	// Test ContainsKeyValue
	if !mm.ContainsKeyValue(1, "apple") || !mm.ContainsKeyValue(2, "carrot") {
		t.Error("Expected key-value pairs to exist")
	}
	if mm.ContainsKeyValue(1, "carrot") {
		t.Error("Expected key-value pair (1, 'carrot') to not exist")
	}
}

func TestMultiMapRemove(t *testing.T) {
	mm := New[string, int]()

	mm.Put("letters", 1)
	mm.Put("letters", 2)
	mm.Put("letters", 3)
	mm.Put("numbers", 10)

	// Test Remove specific value
	if !mm.Remove("letters", 2) {
		t.Error("Expected Remove('letters', 2) to return true")
	}
	remainingLetters := mm.Get("letters")
	expectedRemaining := []int{1, 3}
	if !reflect.DeepEqual(remainingLetters, expectedRemaining) {
		t.Errorf("Expected remaining letters %v, got %v", expectedRemaining, remainingLetters)
	}

	// Test Remove non-existent value
	if mm.Remove("letters", 5) {
		t.Error("Expected Remove('letters', 5) to return false")
	}

	// Test Remove non-existent key
	if mm.Remove("nonexistent", 1) {
		t.Error("Expected Remove('nonexistent', 1) to return false")
	}

	// Test RemoveAll
	if !mm.RemoveAll("letters") {
		t.Error("Expected RemoveAll('letters') to return true")
	}
	if len(mm.Get("letters")) != 0 {
		t.Error("Expected 'letters' to have no values after RemoveAll")
	}
	if mm.KeySize() != 1 {
		t.Errorf("Expected key size 1 after RemoveAll, got %d", mm.KeySize())
	}
}

func TestMultiMapValuesSize(t *testing.T) {
	mm := New[string, int]()

	mm.Put("a", 1)
	mm.Put("a", 2)
	mm.Put("a", 3)
	mm.Put("b", 10)

	if size := mm.ValuesSize("a"); size != 3 {
		t.Errorf("Expected values size of 'a' to be 3, got %d", size)
	}
	if size := mm.ValuesSize("b"); size != 1 {
		t.Errorf("Expected values size of 'b' to be 1, got %d", size)
	}
	if size := mm.ValuesSize("nonexistent"); size != 0 {
		t.Errorf("Expected values size of non-existent key to be 0, got %d", size)
	}
}

func TestMultiMapReplaceAndPutAll(t *testing.T) {
	mm := New[string, int]()

	mm.Put("test", 1)
	mm.Put("test", 2)

	// Test ReplaceValues
	oldValues := mm.ReplaceValues("test", []int{10, 20, 30})
	expectedOld := []int{1, 2}
	if !reflect.DeepEqual(oldValues, expectedOld) {
		t.Errorf("Expected old values %v, got %v", expectedOld, oldValues)
	}

	newValues := mm.Get("test")
	expectedNew := []int{10, 20, 30}
	if !reflect.DeepEqual(newValues, expectedNew) {
		t.Errorf("Expected new values %v, got %v", expectedNew, newValues)
	}

	// Test PutAll
	mm.PutAll("test", []int{40, 50})
	finalValues := mm.Get("test")
	expectedFinal := []int{10, 20, 30, 40, 50}
	if !reflect.DeepEqual(finalValues, expectedFinal) {
		t.Errorf("Expected final values %v, got %v", expectedFinal, finalValues)
	}
}

func TestMultiMapSet(t *testing.T) {
	mm := New[string, int]()

	mm.Put("test", 1)
	mm.Put("test", 2)

	// Test Set replaces all values
	mm.Set("test", 100)
	values := mm.Get("test")
	expected := []int{100}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("Expected values %v after Set, got %v", expected, values)
	}
}

func TestMultiMapKeysValuesEntries(t *testing.T) {
	mm := New[int, string]()

	mm.Put(1, "a")
	mm.Put(1, "b")
	mm.Put(2, "c")
	mm.Put(3, "d")

	// Test Keys
	keys := mm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Test Values
	values := mm.Values()
	if len(values) != 4 {
		t.Errorf("Expected 4 values, got %d", len(values))
	}

	// Test Entries
	entries := mm.Entries()
	if len(entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(entries))
	}

	// Check that all expected key-value pairs exist
	entryMap := make(map[int][]string)
	for _, entry := range entries {
		entryMap[entry.Key] = append(entryMap[entry.Key], entry.Value)
	}

	if len(entryMap[1]) != 2 || !contains(entryMap[1], "a") || !contains(entryMap[1], "b") {
		t.Error("Missing entries for key 1")
	}
	if len(entryMap[2]) != 1 || entryMap[2][0] != "c" {
		t.Error("Missing entry for key 2")
	}
	if len(entryMap[3]) != 1 || entryMap[3][0] != "d" {
		t.Error("Missing entry for key 3")
	}
}

func TestMultiMapForEach(t *testing.T) {
	mm := New[string, int]()

	mm.Put("numbers", 1)
	mm.Put("numbers", 2)
	mm.Put("numbers", 3)
	mm.Put("letters", 10)

	var forEachResult []struct {
		key   string
		value int
	}
	mm.ForEach(func(key string, value int) {
		forEachResult = append(forEachResult, struct {
			key   string
			value int
		}{key, value})
	})

	if len(forEachResult) != 4 {
		t.Errorf("Expected 4 results from ForEach, got %d", len(forEachResult))
	}

	// Check that all expected pairs exist
	pairMap := make(map[string][]int)
	for _, result := range forEachResult {
		pairMap[result.key] = append(pairMap[result.key], result.value)
	}

	if len(pairMap["numbers"]) != 3 {
		t.Error("Expected 3 values for 'numbers' in ForEach")
	}
	if len(pairMap["letters"]) != 1 {
		t.Error("Expected 1 value for 'letters' in ForEach")
	}
}

func TestMultiMapForEachKey(t *testing.T) {
	mm := New[string, int]()

	mm.Put("test", 1)
	mm.Put("test", 2)
	mm.Put("test", 3)

	var forEachKeyResult []int
	mm.ForEachKey("test", func(value int) {
		forEachKeyResult = append(forEachKeyResult, value)
	})

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(forEachKeyResult, expected) {
		t.Errorf("Expected ForEachKey result %v, got %v", expected, forEachKeyResult)
	}

	// Test ForEachKey on non-existent key
	var emptyResult []int
	mm.ForEachKey("nonexistent", func(value int) {
		emptyResult = append(emptyResult, value)
	})
	if len(emptyResult) != 0 {
		t.Error("Expected empty result for ForEachKey on non-existent key")
	}
}

func TestMultiMapToMap(t *testing.T) {
	mm := New[string, int]()

	mm.Put("a", 1)
	mm.Put("a", 2)
	mm.Put("b", 3)

	resultMap := mm.ToMap()
	if len(resultMap) != 2 {
		t.Errorf("Expected 2 keys in result map, got %d", len(resultMap))
	}

	if len(resultMap["a"]) != 2 || resultMap["a"][0] != 1 || resultMap["a"][1] != 2 {
		t.Errorf("Expected a=[1,2], got a=%v", resultMap["a"])
	}

	if len(resultMap["b"]) != 1 || resultMap["b"][0] != 3 {
		t.Errorf("Expected b=[3], got b=%v", resultMap["b"])
	}
}

func TestMultiMapClear(t *testing.T) {
	mm := New[string, int]()

	mm.Put("a", 1)
	mm.Put("a", 2)
	mm.Put("b", 3)

	mm.Clear()

	if !mm.IsEmpty() {
		t.Error("Expected empty multimap after Clear()")
	}
	if mm.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", mm.Size())
	}
	if mm.KeySize() != 0 {
		t.Errorf("Expected key size 0 after Clear(), got %d", mm.KeySize())
	}
	if len(mm.Get("a")) != 0 {
		t.Error("Expected empty values after Clear()")
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}