package concurrentmap

import (
	"fmt"
	"sync"
)

const (
	defaultShardCount = 32
)

// shard represents a single shard of the concurrent map
type shard[K comparable, V any] struct {
	items map[K]V
	mutex sync.RWMutex
}

// ConcurrentMap represents a thread-safe map with sharding for better performance
type ConcurrentMap[K comparable, V any] struct {
	shards    []*shard[K, V]
	shardMask uint32
}

// New creates a new ConcurrentMap with the default number of shards
func New[K comparable, V any]() *ConcurrentMap[K, V] {
	return NewWithShardCount[K, V](defaultShardCount)
}

// NewWithShardCount creates a new ConcurrentMap with the specified number of shards
func NewWithShardCount[K comparable, V any](shardCount int) *ConcurrentMap[K, V] {
	if shardCount <= 0 {
		shardCount = defaultShardCount
	}

	// Ensure shardCount is a power of 2 for efficient hashing
	if (shardCount & (shardCount - 1)) != 0 {
		// Find next power of 2
		shardCount = 1
		for shardCount < defaultShardCount {
			shardCount <<= 1
		}
	}

	shards := make([]*shard[K, V], shardCount)
	for i := range shards {
		shards[i] = &shard[K, V]{
			items: make(map[K]V),
		}
	}

	return &ConcurrentMap[K, V]{
		shards:    shards,
		shardMask: uint32(shardCount - 1),
	}
}

// getShard returns the shard for the given key
func (cm *ConcurrentMap[K, V]) getShard(key K) *shard[K, V] {
	// Use a simple hash function to determine shard
	// In a real implementation, we'd use a proper hash function
	// For now, we'll use a basic approach with type conversion
	hash := cm.hashKey(key)
	return cm.shards[hash&cm.shardMask]
}

// hashKey generates a hash for the given key
// This is a simplified version - in a production system, you'd use a proper hash function
func (cm *ConcurrentMap[K, V]) hashKey(key K) uint32 {
	// Use Go's runtime hash function by converting to string and back
	// This is a simplified implementation - in a real system we'd use a proper hash function
	// such as FNV or similar
	var hash uint32
	switch v := any(key).(type) {
	case string:
		for _, b := range []byte(v) {
			hash = hash*31 + uint32(b)
		}
	case int:
		hash = uint32(v)
	case int32:
		hash = uint32(v)
	case int64:
		hash = uint32(v)
	case uint:
		hash = uint32(v)
	case uint32:
		hash = v
	case uint64:
		hash = uint32(v)
	default:
		// For other types, convert to string representation and hash
		s := any(key)
		str := fmt.Sprintf("%v", s)
		for _, b := range []byte(str) {
			hash = hash*31 + uint32(b)
		}
	}
	return hash
}

// Set adds or updates a key-value pair in the map
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.Lock()
	defer shardPtr.mutex.Unlock()

	shardPtr.items[key] = value
}

// Get retrieves the value for the given key
// Returns the value and true if the key exists, false otherwise
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.RLock()
	defer shardPtr.mutex.RUnlock()

	value, exists := shardPtr.items[key]
	return value, exists
}

// Delete removes a key-value pair from the map
// Returns true if the key existed and was removed, false otherwise
func (cm *ConcurrentMap[K, V]) Delete(key K) bool {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.Lock()
	defer shardPtr.mutex.Unlock()

	_, exists := shardPtr.items[key]
	if exists {
		delete(shardPtr.items, key)
	}
	return exists
}

// Contains checks if a key exists in the map
func (cm *ConcurrentMap[K, V]) Contains(key K) bool {
	_, exists := cm.Get(key)
	return exists
}

// Len returns the total number of elements in the map
func (cm *ConcurrentMap[K, V]) Len() int {
	total := 0
	for _, shard := range cm.shards {
		shard.mutex.RLock()
		total += len(shard.items)
		shard.mutex.RUnlock()
	}
	return total
}

// IsEmpty returns true if the map is empty
func (cm *ConcurrentMap[K, V]) IsEmpty() bool {
	return cm.Len() == 0
}

// Clear removes all key-value pairs from the map
func (cm *ConcurrentMap[K, V]) Clear() {
	for _, shard := range cm.shards {
		shard.mutex.Lock()
		shard.items = make(map[K]V)
		shard.mutex.Unlock()
	}
}

// Keys returns a slice of all keys in the map
func (cm *ConcurrentMap[K, V]) Keys() []K {
	var keys []K
	for _, shard := range cm.shards {
		shard.mutex.RLock()
		for k := range shard.items {
			keys = append(keys, k)
		}
		shard.mutex.RUnlock()
	}
	return keys
}

// Values returns a slice of all values in the map
func (cm *ConcurrentMap[K, V]) Values() []V {
	var values []V
	for _, shard := range cm.shards {
		shard.mutex.RLock()
		for _, v := range shard.items {
			values = append(values, v)
		}
		shard.mutex.RUnlock()
	}
	return values
}

// Items returns a map of all key-value pairs
func (cm *ConcurrentMap[K, V]) Items() map[K]V {
	result := make(map[K]V)
	for _, shard := range cm.shards {
		shard.mutex.RLock()
		for k, v := range shard.items {
			result[k] = v
		}
		shard.mutex.RUnlock()
	}
	return result
}

// ForEach iterates over all key-value pairs in the map
// Note: This is not atomic across the entire map, as different shards may be locked at different times
func (cm *ConcurrentMap[K, V]) ForEach(fn func(key K, value V)) {
	for _, shard := range cm.shards {
		shard.mutex.RLock()
		for k, v := range shard.items {
			fn(k, v)
		}
		shard.mutex.RUnlock()
	}
}

// Size returns the number of shards in the map
func (cm *ConcurrentMap[K, V]) Size() int {
	return len(cm.shards)
}

// ResizeShards changes the number of shards in the map
// Note: This operation is expensive as it requires rehashing all items
func (cm *ConcurrentMap[K, V]) ResizeShards(newShardCount int) {
	if newShardCount <= 0 {
		return
	}

	// Ensure shardCount is a power of 2
	if (newShardCount & (newShardCount - 1)) != 0 {
		// Find next power of 2
		originalCount := newShardCount
		newShardCount = 1
		for newShardCount < originalCount {
			newShardCount <<= 1
		}
	}

	// Get all current items
	allItems := cm.Items()

	// Create new shards
	newShards := make([]*shard[K, V], newShardCount)
	for i := range newShards {
		newShards[i] = &shard[K, V]{
			items: make(map[K]V),
		}
	}

	// Rehash all items to new shards
	for k, v := range allItems {
		hash := cm.hashKey(k)
		shardIdx := hash & uint32(newShardCount-1)
		newShards[shardIdx].items[k] = v
	}

	// Update the map
	cm.shards = newShards
	cm.shardMask = uint32(newShardCount - 1)
}

// TrySet attempts to add a key-value pair if the key doesn't exist
// Returns true if the key was set, false if it already existed
func (cm *ConcurrentMap[K, V]) TrySet(key K, value V) bool {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.Lock()
	defer shardPtr.mutex.Unlock()

	if _, exists := shardPtr.items[key]; exists {
		return false
	}

	shardPtr.items[key] = value
	return true
}

// GetOrSet returns the value for the key if it exists, otherwise sets it to the given value
// Returns the value and a boolean indicating if it was already present
func (cm *ConcurrentMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.Lock()
	defer shardPtr.mutex.Unlock()

	if existingValue, exists := shardPtr.items[key]; exists {
		return existingValue, true
	}

	shardPtr.items[key] = value
	return value, false
}

// Update modifies the value for the given key using the provided function
// Returns true if the key existed and was updated, false otherwise
func (cm *ConcurrentMap[K, V]) Update(key K, updateFn func(V) V) bool {
	shardPtr := cm.getShard(key)
	shardPtr.mutex.Lock()
	defer shardPtr.mutex.Unlock()

	if value, exists := shardPtr.items[key]; exists {
		shardPtr.items[key] = updateFn(value)
		return true
	}
	return false
}
