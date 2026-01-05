package stack

import "sync"

// Stack представляет собой универсальный стек LIFO.
type Stack[T any] struct {
	items []T
	mu    sync.RWMutex
}

// New создает новый стек.
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// Push добавляет элемент в верхнюю часть стека.
func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop удаляет и возвращает элемент из верхней части стека.
func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

// Peek возвращает элемент в верхней части стека, не удаляя его.
func (s *Stack[T]) Peek() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Size возвращает количество элементов в стеке.
func (s *Stack[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// IsEmpty возвращает true, если стек пуст.
func (s *Stack[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear очищает стек.
func (s *Stack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]T, 0)
}
