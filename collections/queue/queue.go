package queue

import (
	"types/collections/list"
)

// Queue представляет собой универсальную очередь FIFO.
type Queue[T any] struct {
	items *list.List[T]
}

// New создает новую очередь.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items: list.New[T](),
	}
}

// Enqueue добавляет элемент в конец очереди.
func (q *Queue[T]) Enqueue(item T) {
	q.items.Add(item)
}

// Dequeue удаляет и возвращает элемент из начала очереди.
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.items.IsEmpty() {
		var zero T
		return zero, false
	}
	item, _ := q.items.Get(0)
	q.items.RemoveAt(0)
	return item, true
}

// Peek возвращает элемент в начале очереди, не удаляя его.
func (q *Queue[T]) Peek() (T, bool) {
	if q.items.IsEmpty() {
		var zero T
		return zero, false
	}
	return q.items.Get(0)
}

// Size возвращает количество элементов в очереди.
func (q *Queue[T]) Size() int {
	return q.items.Size()
}

// IsEmpty возвращает true, если очередь пуста.
func (q *Queue[T]) IsEmpty() bool {
	return q.items.IsEmpty()
}

// Clear очищает очередь.
func (q *Queue[T]) Clear() {
	q.items.Clear()
}
