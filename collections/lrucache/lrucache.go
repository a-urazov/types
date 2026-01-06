package lrucache

import (
	"sync"
)

// Node represents a node in the doubly linked list
type Node[K comparable, V any] struct {
	key   K
	value V
	prev  *Node[K, V]
	next  *Node[K, V]
}

// LRUCache represents a thread-safe LRU (Least Recently Used) cache
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*Node[K, V]
	head     *Node[K, V] // dummy head
	tail     *Node[K, V] // dummy tail
	mutex    sync.RWMutex
}

// New creates a new LRUCache with the given capacity
func New[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		capacity = 1
	}

	// Create dummy head and tail nodes
	head := &Node[K, V]{}
	tail := &Node[K, V]{}
	head.next = tail
	tail.prev = head

	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*Node[K, V]),
		head:     head,
		tail:     tail,
	}
}

// Get retrieves the value associated with the given key
// Returns the value and true if the key exists, false otherwise
// The accessed item is moved to the front (most recently used)
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	node, exists := lru.cache[key]
	if !exists {
		var zero V
		return zero, false
	}

	// Move to front (most recently used)
	lru.moveToFront(node)
	return node.value, true
}

// Put adds or updates a key-value pair in the cache
// If the cache is at capacity, the least recently used item is evicted
func (lru *LRUCache[K, V]) Put(key K, value V) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if node, exists := lru.cache[key]; exists {
		// Update existing key
		node.value = value
		lru.moveToFront(node)
		return
	}

	// Add new key
	newNode := &Node[K, V]{key: key, value: value}
	lru.cache[key] = newNode
	lru.addToFront(newNode)

	// Check capacity and evict if necessary
	if len(lru.cache) > lru.capacity {
		lru.evictLRU()
	}
}

// Remove removes a key-value pair from the cache
// Returns true if the key existed, false otherwise
func (lru *LRUCache[K, V]) Remove(key K) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	node, exists := lru.cache[key]
	if !exists {
		return false
	}

	lru.removeNode(node)
	delete(lru.cache, key)
	return true
}

// Peek retrieves the value without marking it as recently used
// Returns the value and true if the key exists, false otherwise
func (lru *LRUCache[K, V]) Peek(key K) (V, bool) {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	node, exists := lru.cache[key]
	if !exists {
		var zero V
		return zero, false
	}

	return node.value, true
}

// Contains checks if a key exists in the cache
func (lru *LRUCache[K, V]) Contains(key K) bool {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	_, exists := lru.cache[key]
	return exists
}

// Size returns the current number of items in the cache
func (lru *LRUCache[K, V]) Size() int {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	return len(lru.cache)
}

// Capacity returns the maximum capacity of the cache
func (lru *LRUCache[K, V]) Capacity() int {
	return lru.capacity
}

// IsEmpty returns true if the cache is empty
func (lru *LRUCache[K, V]) IsEmpty() bool {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	return len(lru.cache) == 0
}

// Clear removes all items from the cache
func (lru *LRUCache[K, V]) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.cache = make(map[K]*Node[K, V])
	lru.head.next = lru.tail
	lru.tail.prev = lru.head
}

// Keys returns a slice of all keys in the cache (in LRU order, least recent first)
func (lru *LRUCache[K, V]) Keys() []K {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	keys := make([]K, 0, len(lru.cache))
	current := lru.tail.prev
	for current != lru.head {
		keys = append(keys, current.key)
		current = current.prev
	}
	return keys
}

// Values returns a slice of all values in the cache (in LRU order, least recent first)
func (lru *LRUCache[K, V]) Values() []V {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	values := make([]V, 0, len(lru.cache))
	current := lru.tail.prev
	for current != lru.head {
		values = append(values, current.value)
		current = current.prev
	}
	return values
}

// ForEach iterates over all key-value pairs in LRU order (least recent first)
func (lru *LRUCache[K, V]) ForEach(fn func(key K, value V)) {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	current := lru.tail.prev
	for current != lru.head {
		fn(current.key, current.value)
		current = current.prev
	}
}

// ForEachMRU iterates over all key-value pairs in MRU order (most recent first)
func (lru *LRUCache[K, V]) ForEachMRU(fn func(key K, value V)) {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	current := lru.head.next
	for current != lru.tail {
		fn(current.key, current.value)
		current = current.next
	}
}

// moveToFront moves a node to the front of the list (most recently used)
func (lru *LRUCache[K, V]) moveToFront(node *Node[K, V]) {
	lru.removeNode(node)
	lru.addToFront(node)
}

// addToFront adds a node to the front of the list (after dummy head)
func (lru *LRUCache[K, V]) addToFront(node *Node[K, V]) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}

// removeNode removes a node from the doubly linked list
func (lru *LRUCache[K, V]) removeNode(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// evictLRU removes the least recently used item (before dummy tail)
func (lru *LRUCache[K, V]) evictLRU() {
	if lru.tail.prev == lru.head {
		// Cache is empty
		return
	}

	lruNode := lru.tail.prev
	lru.removeNode(lruNode)
	delete(lru.cache, lruNode.key)
}

// Resize changes the capacity of the cache
// If the new capacity is smaller than the current size, excess items are evicted
func (lru *LRUCache[K, V]) Resize(newCapacity int) {
	if newCapacity <= 0 {
		newCapacity = 1
	}

	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.capacity = newCapacity

	// Evict items if necessary
	for len(lru.cache) > lru.capacity {
		lru.evictLRU()
	}
}