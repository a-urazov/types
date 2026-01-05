package priority

import (
	"cmp"
	"container/heap"
	"sync"
)

// item - это элемент в очереди с приоритетом.
type item[TElement any, TPriority cmp.Ordered] struct {
	value    TElement
	priority TPriority
	index    int // Индекс элемента в куче.
}

// internalHeap - это базовая реализация для нашей очереди с приоритетом.
type internalHeap[TElement any, TPriority cmp.Ordered] []*item[TElement, TPriority]

func (h internalHeap[TElement, TPriority]) Len() int { return len(h) }
func (h internalHeap[TElement, TPriority]) Less(i, j int) bool {
	// Мы хотим, чтобы Pop давал нам самый высокий, а не самый низкий приоритет, поэтому мы используем >
	return h[i].priority > h[j].priority
}
func (h internalHeap[TElement, TPriority]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *internalHeap[TElement, TPriority]) Push(x any) {
	n := len(*h)
	it := x.(*item[TElement, TPriority])
	it.index = n
	*h = append(*h, it)
}
func (h *internalHeap[TElement, TPriority]) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	old[n-1] = nil // избежать утечки памяти
	it.index = -1  // для безопасности
	*h = old[0 : n-1]
	return it
}

// Queue представляет собой очередь с приоритетом.
type Queue[TElement any, TPriority cmp.Ordered] struct {
	heap internalHeap[TElement, TPriority]
	mu   sync.Mutex
}

// New создает новую PriorityQueue.
func New[TElement any, TPriority cmp.Ordered]() *Queue[TElement, TPriority] {
	pq := &Queue[TElement, TPriority]{
		heap: make(internalHeap[TElement, TPriority], 0),
	}
	heap.Init(&pq.heap)
	return pq
}

// Enqueue добавляет элемент в очередь с заданным приоритетом.
func (pq *Queue[TElement, TPriority]) Enqueue(value TElement, priority TPriority) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	newItem := &item[TElement, TPriority]{
		value:    value,
		priority: priority,
	}
	heap.Push(&pq.heap, newItem)
}

// Dequeue удаляет и возвращает элемент с наивысшим приоритетом.
func (pq *Queue[TElement, TPriority]) Dequeue() (TElement, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.heap.Len() == 0 {
		var zero TElement
		return zero, false
	}
	it := heap.Pop(&pq.heap).(*item[TElement, TPriority])
	return it.value, true
}

// Size возвращает количество элементов в очереди.
func (pq *Queue[TElement, TPriority]) Size() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.heap.Len()
}

// IsEmpty возвращает true, если очередь пуста.
func (pq *Queue[TElement, TPriority]) IsEmpty() bool {
	return pq.Size() == 0
}
