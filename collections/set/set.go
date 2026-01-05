package set

import "sync"

// Set представляет собой набор уникальных элементов.
type Set[T comparable] struct {
	items map[T]struct{}
	mu    sync.RWMutex
}

// New создает новый HashSet.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

// Add добавляет элемент в набор. Возвращает false, если элемент уже существует.
func (s *Set[T]) Add(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[item]; ok {
		return false
	}
	s.items[item] = struct{}{}
	return true
}

// Remove удаляет элемент из набора.
func (s *Set[T]) Remove(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[item]; !ok {
		return false
	}
	delete(s.items, item)
	return true
}

// Contains проверяет, существует ли элемент в наборе.
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[item]
	return ok
}

// Size возвращает количество элементов в наборе.
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Clear удаляет все элементы из набора.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

// ToArray возвращает срез всех элементов в наборе.
func (s *Set[T]) ToArray() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	arr := make([]T, 0, len(s.items))
	for item := range s.items {
		arr = append(arr, item)
	}
	return arr
}
