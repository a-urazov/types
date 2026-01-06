package heap

import (
	"container/heap"
	"sync"
)

// Comparator - это функция сравнения, которая определяет порядок элементов в куче.
// Для минимальной кучи она должна возвращать true, если a < b.
// Для максимальной кучи она должна возвращать true, если a > b.
type Comparator[T any] func(a, b T) bool

// internalHeap - это внутренняя реализация для нашей кучи.
type internalHeap[T any] struct {
	elements   []T
	comparator Comparator[T]
}

// Len возвращает количество элементов.
func (h internalHeap[T]) Len() int { return len(h.elements) }

// Less сравнивает два элемента.
func (h internalHeap[T]) Less(i, j int) bool {
	return h.comparator(h.elements[i], h.elements[j])
}

// Swap меняет местами два элемента.
func (h internalHeap[T]) Swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

// Push добавляет элемент в кучу.
func (h *internalHeap[T]) Push(x any) {
	h.elements = append(h.elements, x.(T))
}

// Pop удаляет элемент из кучи.
func (h *internalHeap[T]) Pop() any {
	old := h.elements
	n := len(old)
	x := old[n-1]
	var zero T
	old[n-1] = zero // Избегаем утечки памяти
	h.elements = old[0 : n-1]
	return x
}

// Heap представляет собой обобщенную, поточно-безопасную структуру данных кучи.
type Heap[T any] struct {
	internal *internalHeap[T]
	mu       sync.Mutex
}

// New создает новую кучу с заданным компаратором.
func New[T any](comparator Comparator[T]) *Heap[T] {
	h := &Heap[T]{
		internal: &internalHeap[T]{
			elements:   make([]T, 0),
			comparator: comparator,
		},
	}
	heap.Init(h.internal)
	return h
}

// Push добавляет элемент в кучу.
func (h *Heap[T]) Push(value T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	heap.Push(h.internal, value)
}

// Pop удаляет и возвращает верхний элемент из кучи.
func (h *Heap[T]) Pop() (T, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.internal.Len() == 0 {
		var zero T
		return zero, false
	}
	return heap.Pop(h.internal).(T), true
}

// Peek возвращает верхний элемент, не удаляя его.
func (h *Heap[T]) Peek() (T, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.internal.Len() == 0 {
		var zero T
		return zero, false
	}
	return h.internal.elements[0], true
}

// Size возвращает количество элементов в куче.
func (h *Heap[T]) Size() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.internal.Len()
}

// IsEmpty возвращает true, если куча пуста.
func (h *Heap[T]) IsEmpty() bool {
	return h.Size() == 0
}
