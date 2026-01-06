package priority

import (
	"cmp"
	"types/collections/heap"
)

// item - это элемент в очереди с приоритетом.
type item[TElement any, TPriority cmp.Ordered] struct {
	value    TElement
	priority TPriority
}

// Queue представляет собой очередь с приоритетом.
type Queue[TElement any, TPriority cmp.Ordered] struct {
	heap *heap.Heap[*item[TElement, TPriority]]
}

// New создает новую PriorityQueue.
func New[TElement any, TPriority cmp.Ordered]() *Queue[TElement, TPriority] {
	// Для максимальной приоритетной очереди мы хотим, чтобы элемент с большим значением приоритета был "меньше".
	maxHeap := heap.New(func(a, b *item[TElement, TPriority]) bool {
		return a.priority > b.priority
	})
	pq := &Queue[TElement, TPriority]{
		heap: maxHeap,
	}
	return pq
}

// Enqueue добавляет элемент в очередь с заданным приоритетом.
func (pq *Queue[TElement, TPriority]) Enqueue(value TElement, priority TPriority) {
	newItem := &item[TElement, TPriority]{
		value:    value,
		priority: priority,
	}
	pq.heap.Push(newItem)
}

// Dequeue удаляет и возвращает элемент с наивысшим приоритетом.
func (pq *Queue[TElement, TPriority]) Dequeue() (TElement, bool) {
	it, ok := pq.heap.Pop()
	if !ok {
		var zero TElement
		return zero, false
	}
	return it.value, true
}

// Size возвращает количество элементов в очереди.
func (pq *Queue[TElement, TPriority]) Size() int {
	return pq.heap.Size()
}

// IsEmpty возвращает true, если очередь пуста.
func (pq *Queue[TElement, TPriority]) IsEmpty() bool {
	return pq.heap.IsEmpty()
}
