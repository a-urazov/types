package stack

import (
	"types/collections/list"
)

// Stack представляет собой универсальный стек LIFO.
type Stack[T any] struct {
	items *list.List[T]
}

// New создает новый стек.
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: list.New[T](),
	}
}

// Push добавляет элемент в верхнюю часть стека.
func (s *Stack[T]) Push(item T) {
	s.items.Add(item)
}

// Pop удаляет и возвращает элемент из верхней части стека.
func (s *Stack[T]) Pop() (T, bool) {
	if s.items.IsEmpty() {
		var zero T
		return zero, false
	}
	size := s.items.Size()
	item, _ := s.items.Get(size - 1)
	s.items.RemoveAt(size - 1)
	return item, true
}

// Peek возвращает элемент в верхней части стека, не удаляя его.
func (s *Stack[T]) Peek() (T, bool) {
	if s.items.IsEmpty() {
		var zero T
		return zero, false
	}
	size := s.items.Size()
	return s.items.Get(size - 1)
}

// Size возвращает количество элементов в стеке.
func (s *Stack[T]) Size() int {
	return s.items.Size()
}

// IsEmpty возвращает true, если стек пуст.
func (s *Stack[T]) IsEmpty() bool {
	return s.items.IsEmpty()
}

// Clear очищает стек.
func (s *Stack[T]) Clear() {
	s.items.Clear()
}
