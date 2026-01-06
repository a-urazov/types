package multimap

import (
	"sync"
)

// MultiMap represents a thread-safe map where one key can have multiple values
type MultiMap[K comparable, V any] struct {
	items map[K][]V
	mutex sync.RWMutex
}

// New creates a new empty MultiMap
func New[K comparable, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{
		items: make(map[K][]V),
	}
}

// Put adds a value to the list of values for the given key
func (mm *MultiMap[K, V]) Put(key K, value V) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if values, exists := mm.items[key]; exists {
		mm.items[key] = append(values, value)
	} else {
		mm.items[key] = []V{value}
	}
}

// Get returns all values associated with the given key
// Returns an empty slice if the key does not exist
func (mm *MultiMap[K, V]) Get(key K) []V {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists {
		// Return a copy of the slice to prevent external modification
		result := make([]V, len(values))
		copy(result, values)
		return result
	}
	return []V{}
}

// GetFirst returns the first value associated with the given key
// Returns the value and true if the key exists, false otherwise
func (mm *MultiMap[K, V]) GetFirst(key K) (V, bool) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists && len(values) > 0 {
		return values[0], true
	}
	var zero V
	return zero, false
}

// GetLast returns the last value associated with the given key
// Returns the value and true if the key exists, false otherwise
func (mm *MultiMap[K, V]) GetLast(key K) (V, bool) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists && len(values) > 0 {
		return values[len(values)-1], true
	}
	var zero V
	return zero, false
}

// Remove removes a specific value from the list of values for the given key
// Returns true if the value was found and removed, false otherwise
func (mm *MultiMap[K, V]) Remove(key K, value V) bool {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if values, exists := mm.items[key]; exists {
		for i, v := range values {
			if any(v) == any(value) {
				// Remove the value at index i
				newValues := append(values[:i], values[i+1:]...)
				if len(newValues) == 0 {
					delete(mm.items, key)
				} else {
					mm.items[key] = newValues
				}
				return true
			}
		}
	}
	return false
}

// RemoveAll removes all values associated with the given key
// Returns true if the key existed, false otherwise
func (mm *MultiMap[K, V]) RemoveAll(key K) bool {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	_, exists := mm.items[key]
	if exists {
		delete(mm.items, key)
	}
	return exists
}

// RemoveValue removes the first occurrence of a value from the list of values for the given key
// Returns true if the value was found and removed, false otherwise
func (mm *MultiMap[K, V]) RemoveValue(key K, value V) bool {
	return mm.Remove(key, value)
}

// ContainsKey checks if the given key exists in the multimap
func (mm *MultiMap[K, V]) ContainsKey(key K) bool {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	_, exists := mm.items[key]
	return exists
}

// ContainsValue checks if the given value exists for any key in the multimap
func (mm *MultiMap[K, V]) ContainsValue(value V) bool {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	for _, values := range mm.items {
		for _, v := range values {
			if any(v) == any(value) {
				return true
			}
		}
	}
	return false
}

// ContainsKeyValue checks if the given key-value pair exists in the multimap
func (mm *MultiMap[K, V]) ContainsKeyValue(key K, value V) bool {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists {
		for _, v := range values {
			if any(v) == any(value) {
				return true
			}
		}
	}
	return false
}

// Size returns the total number of key-value pairs in the multimap
func (mm *MultiMap[K, V]) Size() int {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	total := 0
	for _, values := range mm.items {
		total += len(values)
	}
	return total
}

// KeySize returns the number of keys in the multimap
func (mm *MultiMap[K, V]) KeySize() int {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	return len(mm.items)
}

// ValuesSize returns the number of values associated with the given key
func (mm *MultiMap[K, V]) ValuesSize(key K) int {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists {
		return len(values)
	}
	return 0
}

// IsEmpty returns true if the multimap is empty
func (mm *MultiMap[K, V]) IsEmpty() bool {
	return mm.Size() == 0
}

// Clear removes all key-value pairs from the multimap
func (mm *MultiMap[K, V]) Clear() {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mm.items = make(map[K][]V)
}

// Keys returns a slice of all keys in the multimap
func (mm *MultiMap[K, V]) Keys() []K {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	keys := make([]K, 0, len(mm.items))
	for k := range mm.items {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the multimap
func (mm *MultiMap[K, V]) Values() []V {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	var values []V
	for _, valueList := range mm.items {
		values = append(values, valueList...)
	}
	return values
}

// Entries returns a slice of all key-value pairs in the multimap
func (mm *MultiMap[K, V]) Entries() []struct {
	Key   K
	Value V
} {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	var entries []struct {
		Key   K
		Value V
	}
	for k, valueList := range mm.items {
		for _, v := range valueList {
			entries = append(entries, struct {
				Key   K
				Value V
			}{k, v})
		}
	}
	return entries
}

// ForEach iterates over all key-value pairs in the multimap
func (mm *MultiMap[K, V]) ForEach(fn func(key K, value V)) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	for k, valueList := range mm.items {
		for _, v := range valueList {
			fn(k, v)
		}
	}
}

// ForEachKey iterates over all values for a specific key
func (mm *MultiMap[K, V]) ForEachKey(key K, fn func(value V)) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	if values, exists := mm.items[key]; exists {
		for _, v := range values {
			fn(v)
		}
	}
}

// ReplaceValues replaces all values for the given key with a new list of values
// Returns the old values associated with the key
func (mm *MultiMap[K, V]) ReplaceValues(key K, values []V) []V {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	oldValues := mm.items[key]
	if len(values) == 0 {
		delete(mm.items, key)
	} else {
		mm.items[key] = values
	}
	return oldValues
}

// PutAll adds all values from the provided slice to the list of values for the given key
func (mm *MultiMap[K, V]) PutAll(key K, values []V) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if existingValues, exists := mm.items[key]; exists {
		mm.items[key] = append(existingValues, values...)
	} else {
		mm.items[key] = values
	}
}

// Set replaces all values for the given key with a single value
func (mm *MultiMap[K, V]) Set(key K, value V) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mm.items[key] = []V{value}
}

// ToMap returns a map where each key maps to a slice of its values
func (mm *MultiMap[K, V]) ToMap() map[K][]V {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	result := make(map[K][]V, len(mm.items))
	for k, v := range mm.items {
		// Create a copy of the slice to prevent external modification
		valuesCopy := make([]V, len(v))
		copy(valuesCopy, v)
		result[k] = valuesCopy
	}
	return result
}