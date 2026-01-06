package heap

import (
	"container/heap"
	"sync"
	"types/collections/list"
)

// Comparator - это функция сравнения, которая определяет порядок элементов в куче.
// Для минимальной кучи она должна возвращать true, если a < b.
// Для максимальной кучи она должна возвращать true, если a > b.
type Comparator[T any] func(a, b T) bool

// internalHeap - это внутренняя реализация для нашей кучи.
type internalHeap[T any] struct {
	elements   *list.List[T]
	comparator Comparator[T]
}

// Len возвращает количество элементов.
func (h internalHeap[T]) Len() int { return h.elements.Size() }

// Less сравнивает два элемента.
func (h internalHeap[T]) Less(i, j int) bool {
	itemI, _ := h.elements.Get(i)
	itemJ, _ := h.elements.Get(j)
	return h.comparator(itemI, itemJ)
}

// Swap меняет местами два элемента.
func (h internalHeap[T]) Swap(i, j int) {
	itemI, _ := h.elements.Get(i)
	itemJ, _ := h.elements.Get(j)
	h.elements.Set(i, itemJ)
	h.elements.Set(j, itemI)
}

// Push добавляет элемент в кучу.
func (h *internalHeap[T]) Push(x any) {
	h.elements.Add(x.(T))
}

// Pop удаляет элемент из кучи.
func (h *internalHeap[T]) Pop() any {
	n := h.elements.Size()
	x, _ := h.elements.Get(n - 1)
	h.elements.RemoveAt(n - 1)
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
			elements:   list.New[T](),
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
	return h.internal.elements.Get(0)
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
