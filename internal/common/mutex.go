package common

import (
	"maps"
	"sync"
)

// RWLocker provides a common interface for read-write mutex operations
type RWLocker interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

// Vector provides thread-safe operations on slices
type Vector[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewVector creates a new thread-safe slice
func NewVector[T any]() *Vector[T] {
	return &Vector[T]{
		items: make([]T, 0),
	}
}

// WithWriteLock executes the given function with write lock
func (s *Vector[T]) WithWriteLock(fn func(items []T) []T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = fn(s.items)
}

// WithReadLock executes the given function with read lock
func (s *Vector[T]) WithReadLock(fn func(items []T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s.items)
}

// GetItems returns a copy of the items (thread-safe)
func (s *Vector[T]) GetItems() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]T, len(s.items))
	copy(items, s.items)
	return items
}

// SetItems replaces all items (thread-safe)
func (s *Vector[T]) SetItems(items []T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]T, len(items))
	copy(s.items, items)
}

// Len returns the length of the slice (thread-safe)
func (s *Vector[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Map provides thread-safe operations on maps
type Map[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewMap creates a new thread-safe map
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		items: make(map[K]V),
	}
}

// WithWriteLock executes the given function with write lock
func (m *Map[K, V]) WithWriteLock(fn func(items map[K]V) map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = fn(m.items)
}

// WithReadLock executes the given function with read lock
func (m *Map[K, V]) WithReadLock(fn func(items map[K]V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fn(m.items)
}

// GetItems returns a copy of the items (thread-safe)
func (m *Map[K, V]) GetItems() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make(map[K]V, len(m.items))
	maps.Copy(items, m.items)
	return items
}

// SetItems replaces all items (thread-safe)
func (m *Map[K, V]) SetItems(items map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[K]V, len(items))
	maps.Copy(m.items, items)
}

// Len returns the length of the map (thread-safe)
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}
