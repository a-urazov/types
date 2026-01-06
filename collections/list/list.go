package list

import (
	"reflect"
	"types/internal/common"
	"types/sort"
)

// List представляет собой универсальный список.
type List[T any] struct {
	items *common.Vector[T]
}

// New создает новый список.
func New[T any]() *List[T] {
	return &List[T]{
		items: common.NewVector[T](),
	}
}

// Add добавляет элемент в конец списка.
func (l *List[T]) Add(item T) {
	l.items.Lock(func(items []T) []T {
		return append(items, item)
	})
}

// Insert добавляет элемент по указанному индексу.
func (l *List[T]) Insert(index int, item T) bool {
	var result bool
	l.items.Lock(func(items []T) []T {
		if index < 0 || index > len(items) {
			result = false
			return items
		}
		result = true
		return append(items[:index], append([]T{item}, items[index:]...)...)
	})
	return result
}

// Remove удаляет первое вхождение элемента из списка.
func (l *List[T]) Remove(item T) bool {
	var result bool
	l.items.Lock(func(items []T) []T {
		for i, v := range items {
			if reflect.DeepEqual(v, item) {
				result = true
				return append(items[:i], items[i+1:]...)
			}
		}
		result = false
		return items
	})
	return result
}

// Get возвращает элемент по указанному индексу.
func (l *List[T]) Get(index int) (T, bool) {
	var result T
	var ok bool
	l.items.RLock(func(items []T) {
		if index < 0 || index >= len(items) {
			var zero T
			result = zero
			ok = false
			return
		}
		result = items[index]
		ok = true
	})
	return result, ok
}

// Set устанавливает элемент по указанному индексу.
func (l *List[T]) Set(index int, item T) bool {
	var result bool
	l.items.Lock(func(items []T) []T {
		if index < 0 || index >= len(items) {
			result = false
			return items
		}
		items[index] = item
		result = true
		return items
	})
	return result
}

// Size возвращает количество элементов в списке.
func (l *List[T]) Size() int {
	return l.items.Len()
}

// IndexOf возвращает индекс первого вхождения элемента в списке.
// Возвращает -1, если элемент не найден.
func (l *List[T]) IndexOf(item T) int {
	var result int = -1
	l.items.RLock(func(items []T) {
		for i, v := range items {
			if reflect.DeepEqual(v, item) {
				result = i
				return
			}
		}
	})
	return result
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
	l.items.Lock(func(items []T) []T {
		return make([]T, 0)
	})
}

// ToSlice возвращает срез, содержащий все элементы списка.
func (l *List[T]) ToSlice() []T {
	return l.items.GetItems()
}

// AddRange добавляет элементы из среза в список.
func (l *List[T]) AddRange(items []T) {
	l.items.Lock(func(existing []T) []T {
		return append(existing, items...)
	})
}

func (l *List[T]) RemoveAt(index int) bool {
	var result bool
	l.items.Lock(func(items []T) []T {
		if index < 0 || index >= len(items) {
			result = false
			return items
		}
		result = true
		return append(items[:index], items[index+1:]...)
	})
	return result
}

func (l *List[T]) ForEach(action func(T)) {
	l.items.RLock(func(items []T) {
		for _, item := range items {
			action(item)
		}
	})
}

func (l *List[T]) Filter(predicate func(T) bool) *List[T] {
	filtered := New[T]()
	l.items.RLock(func(items []T) {
		for _, item := range items {
			if predicate(item) {
				filtered.Add(item)
			}
		}
	})
	return filtered
}

func (l *List[T]) Reduce(initial T, reducer func(T, T) T) T {
	var result T = initial
	l.items.RLock(func(items []T) {
		for _, item := range items {
			result = reducer(result, item)
		}
	})
	return result
}

// Reverse изменяет порядок элементов в списке на обратный.
func (l *List[T]) Reverse() {
	l.items.Lock(func(items []T) []T {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
		return items
	})
}

func (l *List[T]) Sort(cmp func(a, b T) bool) {
	l.items.Lock(func(items []T) []T {
		sort.Slice(items, func(a, b T) bool {
			return cmp(a, b)
		})
		return items
	})
}

func (l *List[T]) ToArray() []T {
	return l.items.GetItems()
}
