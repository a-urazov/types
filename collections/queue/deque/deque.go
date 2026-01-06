package deque

import "types/collections/list"

// Deque представляет двустороннюю очередь (deque).
type Deque[T any] struct {
	elements *list.List[T]
}

// New создает и возвращает новую пустую Deque.
func New[T any]() *Deque[T] {
	return &Deque[T]{
		elements: list.New[T](),
	}
}

// PushFront добавляет элемент в начало очереди.
func (d *Deque[T]) PushFront(item T) {
	d.elements.Insert(0, item)
}

// PushBack добавляет элемент в конец очереди.
func (d *Deque[T]) PushBack(item T) {
	d.elements.Add(item)
}

// PopFront удаляет и возвращает элемент из начала очереди.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) PopFront() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	item, _ := d.elements.Get(0)
	d.elements.RemoveAt(0)
	return item, true
}

// PopBack удаляет и возвращает элемент из конца очереди.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) PopBack() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	size := d.elements.Size()
	item, _ := d.elements.Get(size - 1)
	d.elements.RemoveAt(size - 1)
	return item, true
}

// Front возвращает элемент из начала очереди без его удаления.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) Front() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	return d.elements.Get(0)
}

// Back возвращает элемент из конца очереди без его удаления.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) Back() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	size := d.elements.Size()
	return d.elements.Get(size - 1)
}

// IsEmpty возвращает true, если очередь пуста, иначе false.
func (d *Deque[T]) IsEmpty() bool {
	return d.elements.IsEmpty()
}

// Size возвращает количество элементов в очереди.
func (d *Deque[T]) Size() int {
	return d.elements.Size()
}
