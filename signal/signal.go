// Package signal provides generic signaling mechanisms for inter-process communication
// using existing data structures from the types library.
package signal

import (
	"sync"
	"types/collections/queue"
	"types/collections/set"
)

// Signal represents a communication channel that can broadcast messages to multiple listeners
type Signal[T any] struct {
	mutex     sync.RWMutex
	listeners map[string]chan T
	nextID    int
}

// BufferedSignal uses a queue from the collections package to buffer signals
type BufferedSignal[T any] struct {
	mutex      sync.RWMutex
	listeners  map[string]chan T
	signalBuf  *queue.Queue[T]  // Using queue from collections
	listenerID *set.Set[string] // Using set from collections
	nextID     int
	bufSize    int
}

// New creates a new Signal instance
func New[T any]() *Signal[T] {
	return &Signal[T]{
		listeners: make(map[string]chan T),
		nextID:    0,
	}
}

// NewBuffered creates a new BufferedSignal instance with a specified buffer size
func NewBuffered[T any](bufSize int) *BufferedSignal[T] {
	return &BufferedSignal[T]{
		listeners:  make(map[string]chan T),
		signalBuf:  queue.New[T](),
		listenerID: set.New[string](),
		nextID:     0,
		bufSize:    bufSize,
	}
}

// Subscribe adds a new listener to the signal and returns a channel to receive messages
// and an unsubscribe function
func (s *Signal[T]) Subscribe(bufferSize int) (<-chan T, func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := s.generateID()
	ch := make(chan T, bufferSize)
	s.listeners[id] = ch

	unsubscribe := func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if listenerCh, exists := s.listeners[id]; exists {
			close(listenerCh)
			delete(s.listeners, id)
		}
	}

	return ch, unsubscribe
}

// Subscribe adds a new listener to the buffered signal and returns a channel to receive messages
// and an unsubscribe function
func (bs *BufferedSignal[T]) Subscribe(bufferSize int) (<-chan T, func()) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	id := bs.generateID()
	ch := make(chan T, bufferSize)
	bs.listeners[id] = ch
	bs.listenerID.Add(id)

	unsubscribe := func() {
		bs.mutex.Lock()
		defer bs.mutex.Unlock()
		if listenerCh, exists := bs.listeners[id]; exists {
			close(listenerCh)
			delete(bs.listeners, id)
			bs.listenerID.Remove(id)
		}
	}

	return ch, unsubscribe
}

// Emit sends a value to all subscribed listeners
func (s *Signal[T]) Emit(value T) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, ch := range s.listeners {
		select {
		case ch <- value:
		default:
			// Skip if channel is full (non-blocking)
		}
	}
}

// Emit sends a value to all subscribed listeners, using the buffered queue
func (bs *BufferedSignal[T]) Emit(value T) {
	bs.mutex.Lock()

	// Add to internal buffer
	bs.signalBuf.Enqueue(value)

	// Keep buffer size within limit
	for bs.signalBuf.Size() > bs.bufSize {
		_, _ = bs.signalBuf.Dequeue() // Discard oldest item
	}

	listeners := make([]chan T, 0, len(bs.listeners))
	for _, ch := range bs.listeners {
		listeners = append(listeners, ch)
	}
	bs.mutex.Unlock()

	// Send to all listeners
	for _, ch := range listeners {
		select {
		case ch <- value:
		default:
			// Skip if channel is full (non-blocking)
		}
	}
}

// EmitSync sends a value to all subscribed listeners synchronously (blocking until all receive)
func (s *Signal[T]) EmitSync(value T) {
	s.mutex.RLock()
	listeners := make([]chan T, 0, len(s.listeners))
	for _, ch := range s.listeners {
		listeners = append(listeners, ch)
	}
	s.mutex.RUnlock()

	for _, ch := range listeners {
		ch <- value
	}
}

// EmitSync sends a value to all subscribed listeners synchronously (blocking until all receive)
func (bs *BufferedSignal[T]) EmitSync(value T) {
	bs.mutex.Lock()

	// Add to internal buffer
	bs.signalBuf.Enqueue(value)

	// Keep buffer size within limit
	for bs.signalBuf.Size() > bs.bufSize {
		_, _ = bs.signalBuf.Dequeue() // Discard oldest item
	}

	listeners := make([]chan T, 0, len(bs.listeners))
	for _, ch := range bs.listeners {
		listeners = append(listeners, ch)
	}
	bs.mutex.Unlock()

	for _, ch := range listeners {
		ch <- value
	}
}

// Broadcast sends a value to all subscribed listeners concurrently
func (s *Signal[T]) Broadcast(value T) {
	s.mutex.RLock()
	listeners := make([]chan T, 0, len(s.listeners))
	for _, ch := range s.listeners {
		listeners = append(listeners, ch)
	}
	s.mutex.RUnlock()

	var wg sync.WaitGroup
	for _, ch := range listeners {
		wg.Add(1)
		go func(c chan T) {
			defer wg.Done()
			c <- value
		}(ch)
	}
	wg.Wait()
}

// Broadcast sends a value to all subscribed listeners concurrently
func (bs *BufferedSignal[T]) Broadcast(value T) {
	bs.mutex.Lock()

	// Add to internal buffer
	bs.signalBuf.Enqueue(value)

	// Keep buffer size within limit
	for bs.signalBuf.Size() > bs.bufSize {
		bs.signalBuf.Dequeue()
	}

	listeners := make([]chan T, 0, len(bs.listeners))
	for _, ch := range bs.listeners {
		listeners = append(listeners, ch)
	}
	bs.mutex.Unlock()

	var wg sync.WaitGroup
	for _, ch := range listeners {
		wg.Add(1)
		go func(c chan T) {
			defer wg.Done()
			c <- value
		}(ch)
	}
	wg.Wait()
}

// generateID creates a unique ID for a listener
func (s *Signal[T]) generateID() string {
	id := s.nextID
	s.nextID++
	return "listener_" + string(rune(id+'0')) // Simple ID generation
}

// generateID creates a unique ID for a buffered signal listener
func (bs *BufferedSignal[T]) generateID() string {
	id := bs.nextID
	bs.nextID++
	return "buffered_listener_" + string(rune(id+'0')) // Simple ID generation
}

// GetBuffer returns a copy of the signal buffer
func (bs *BufferedSignal[T]) GetBuffer() []T {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	size := bs.signalBuf.Size()
	if size == 0 {
		return []T{}
	}

	items := make([]T, 0, size)
	tempQueue := queue.New[T]()

	// Copy items maintaining order
	for i := 0; i < size; i++ {
		item, ok := bs.signalBuf.Dequeue()
		if ok {
			items = append(items, item)
			tempQueue.Enqueue(item)
		}
	}

	// Restore original queue
	for !tempQueue.IsEmpty() {
		item, ok := tempQueue.Dequeue()
		if ok {
			bs.signalBuf.Enqueue(item)
		}
	}

	return items
}

// ClearBuffer clears the internal buffer
func (bs *BufferedSignal[T]) ClearBuffer() {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Dequeue all items to clear the buffer
	for !bs.signalBuf.IsEmpty() {
		_, _ = bs.signalBuf.Dequeue() // Discard the return values
	}
}
