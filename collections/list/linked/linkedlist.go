package linked

import "sync"

// node представляет узел в связанном списке.
type node[T any] struct {
	value T
	prev  *node[T]
	next  *node[T]
}

// List представляет собой универсальный двусвязный список.
type List[T any] struct {
	head *node[T]
	tail *node[T]
	size int
	mu   sync.RWMutex
}

// New создает новый LinkedList.
func New[T any]() *List[T] {
	return &List[T]{}
}

// AddFirst добавляет новое значение в начало списка.
func (l *List[T]) AddFirst(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	newNode := &node[T]{value: value, next: l.head}
	if l.head != nil {
		l.head.prev = newNode
	}
	l.head = newNode
	if l.tail == nil {
		l.tail = newNode
	}
	l.size++
}

// AddLast добавляет новое значение в конец списка.
func (l *List[T]) AddLast(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	newNode := &node[T]{value: value, prev: l.tail}
	if l.tail != nil {
		l.tail.next = newNode
	}
	l.tail = newNode
	if l.head == nil {
		l.head = newNode
	}
	l.size++
}

// RemoveFirst удаляет первый элемент из списка.
func (l *List[T]) RemoveFirst() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.head == nil {
		return false
	}
	l.head = l.head.next
	if l.head != nil {
		l.head.prev = nil
	} else {
		l.tail = nil // Список теперь пуст
	}
	l.size--
	return true
}

// RemoveLast удаляет последний элемент из списка.
func (l *List[T]) RemoveLast() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tail == nil {
		return false
	}
	l.tail = l.tail.prev
	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil // Список теперь пуст
	}
	l.size--
	return true
}

// First возвращает значение первого элемента.
func (l *List[T]) First() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// Last возвращает значение последнего элемента.
func (l *List[T]) Last() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.tail == nil {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// Size возвращает количество элементов в списке.
func (l *List[T]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.size
}

// IsEmpty возвращает true, если список пуст.
func (l *List[T]) IsEmpty() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.size == 0
}

// Clear очищает список.
func (l *List[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.head = nil
	l.tail = nil
	l.size = 0
}

// ToSlice возвращает срез, содержащий все элементы списка.
func (l *List[T]) ToSlice() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	slice := make([]T, 0, l.size)
	for n := l.head; n != nil; n = n.next {
		slice = append(slice, n.value)
	}
	return slice
}
