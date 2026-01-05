package deque

// Deque представляет двустороннюю очередь (deque).
type Deque[T any] struct {
	elements []T
}

// New создает и возвращает новую пустую Deque.
func New[T any]() *Deque[T] {
	return &Deque[T]{
		elements: make([]T, 0),
	}
}

// PushFront добавляет элемент в начало очереди.
func (d *Deque[T]) PushFront(item T) {
	d.elements = append([]T{item}, d.elements...)
}

// PushBack добавляет элемент в конец очереди.
func (d *Deque[T]) PushBack(item T) {
	d.elements = append(d.elements, item)
}

// PopFront удаляет и возвращает элемент из начала очереди.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) PopFront() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	item := d.elements[0]
	d.elements = d.elements[1:]
	return item, true
}

// PopBack удаляет и возвращает элемент из конца очереди.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) PopBack() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	item := d.elements[len(d.elements)-1]
	d.elements = d.elements[:len(d.elements)-1]
	return item, true
}

// Front возвращает элемент из начала очереди без его удаления.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) Front() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	return d.elements[0], true
}

// Back возвращает элемент из конца очереди без его удаления.
// Возвращает элемент и булево значение, указывающее на успех операции.
func (d *Deque[T]) Back() (T, bool) {
	if d.IsEmpty() {
		var zero T
		return zero, false
	}
	return d.elements[len(d.elements)-1], true
}

// IsEmpty возвращает true, если очередь пуста, иначе false.
func (d *Deque[T]) IsEmpty() bool {
	return d.Size() == 0
}

// Size возвращает количество элементов в очереди.
func (d *Deque[T]) Size() int {
	return len(d.elements)
}
