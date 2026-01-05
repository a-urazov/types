package queue

import "sync"

// Queue представляет собой универсальную очередь FIFO.
type Queue[T any] struct {
	items []T
	mu    sync.RWMutex
}

// New создает новую очередь.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// Enqueue добавляет элемент в конец очереди.
func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue удаляет и возвращает элемент из начала очереди.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek возвращает элемент в начале очереди, не удаляя его.
func (q *Queue[T]) Peek() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Size возвращает количество элементов в очереди.
func (q *Queue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}

// IsEmpty возвращает true, если очередь пуста.
func (q *Queue[T]) IsEmpty() bool {
	return q.Size() == 0
}

// Clear очищает очередь.
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make([]T, 0)
}
