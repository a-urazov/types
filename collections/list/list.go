package list

import (
	"sync"
	"types/sort"
)

// List представляет собой универсальный список.
type List[T comparable] struct {
	items []T
	mu    sync.RWMutex
}

// New создает новый список.
func New[T comparable]() *List[T] {
	return &List[T]{
		items: make([]T, 0),
	}
}

// Add добавляет элемент в конец списка.
func (l *List[T]) Add(item T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, item)
}

// Insert добавляет элемент по указанному индексу.
func (l *List[T]) Insert(index int, item T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index > len(l.items) {
		return false
	}
	l.items = append(l.items[:index], append([]T{item}, l.items[index:]...)...)
	return true
}

// Remove удаляет первое вхождение элемента из списка.
func (l *List[T]) Remove(item T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			return true
		}
	}
	return false
}

// Get возвращает элемент по указанному индексу.
func (l *List[T]) Get(index int) (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, false
	}
	return l.items[index], true
}

// Size возвращает количество элементов в списке.
func (l *List[T]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items)
}

// IndexOf возвращает индекс первого вхождения элемента в списке.
// Возвращает -1, если элемент не найден.
func (l *List[T]) IndexOf(item T) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i, v := range l.items {
		if v == item {
			return i
		}
	}
	return -1
}

// Contains проверяет, существует ли элемент в списке.
func (l *List[T]) Contains(item T) bool {
	return l.IndexOf(item) != -1
}

// IsEmpty возвращает true, если список пуст.
func (l *List[T]) IsEmpty() bool {
	return l.Size() == 0
}

// Clear очищает список.
func (l *List[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = make([]T, 0)
}

// ToSlice возвращает срез, содержащий все элементы списка.
func (l *List[T]) ToSlice() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	arr := make([]T, len(l.items))
	copy(arr, l.items)
	return arr
}

// AddRange добавляет элементы из среза в список.
func (l *List[T]) AddRange(items []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, items...)
}

func (l *List[T]) RemoveAt(index int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index >= len(l.items) {
		return false
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	return true
}

func (l *List[T]) ForEach(action func(T)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, item := range l.items {
		action(item)
	}
}

func (l *List[T]) Filter(predicate func(T) bool) *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	filtered := New[T]()
	for _, item := range l.items {
		if predicate(item) {
			filtered.Add(item)
		}
	}
	return filtered
}

// Reverse изменяет порядок элементов в списке на обратный.
func (l *List[T]) Reverse() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, j := 0, len(l.items)-1; i < j; i, j = i+1, j-1 {
		l.items[i], l.items[j] = l.items[j], l.items[i]
	}
}

func (l *List[T]) Sort(cmp func (a, b T) bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	sort.Slice[T](l.items, func(a, b T) bool {
		return cmp(a, b)
	})
}

func (l *List[T]) ToArray() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	arr := make([]T, len(l.items))
	copy(arr, l.items)
	return arr
}
