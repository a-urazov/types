package set

import (
	"types/collections/dictionary"
)

// Set представляет собой набор уникальных элементов.
type Set[T comparable] struct {
	items *dictionary.Dictionary[T, struct{}]
}

// New создает новый HashSet.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		items: dictionary.New[T, struct{}](),
	}
}

// Add добавляет элемент в набор. Возвращает false, если элемент уже существует.
func (s *Set[T]) Add(item T) bool {
	if _, ok := s.items.Get(item); ok {
		return false
	}
	s.items.Set(item, struct{}{})
	return true
}

// Remove удаляет элемент из набора.
func (s *Set[T]) Remove(item T) bool {
	if _, ok := s.items.Get(item); !ok {
		return false
	}
	s.items.Remove(item)
	return true
}

// Contains проверяет, существует ли элемент в наборе.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items.Get(item)
	return ok
}

// Size возвращает количество элементов в наборе.
func (s *Set[T]) Size() int {
	return s.items.Size()
}

// Clear удаляет все элементы из набора.
func (s *Set[T]) Clear() {
	s.items.Clear()
}

// ToArray возвращает срез всех элементов в наборе.
func (s *Set[T]) ToArray() []T {
	arr := make([]T, 0, s.items.Size())
	for _, item := range s.items.Keys() {
		arr = append(arr, item)
	}
	return arr
}
