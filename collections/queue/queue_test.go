package queue

import (
	"sync"
	"testing"
)

func TestNewQueue(t *testing.T) {
	q := New[int]()
	if !q.IsEmpty() || q.Size() != 0 {
		t.Error("New queue should be empty")
	}
}

func TestEnqueue(t *testing.T) {
	q := New[string]()
	q.Enqueue("first")
	q.Enqueue("second")
	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %d", q.Size())
	}
	val, _ := q.Peek()
	if val != "first" {
		t.Errorf("Expected to peek 'first', got '%s'", val)
	}
}

func TestDequeue(t *testing.T) {
	q := New[int]()
	q.Enqueue(10)
	q.Enqueue(20)

	val, ok := q.Dequeue()
	if !ok || val != 10 {
		t.Errorf("Expected to dequeue 10, got %d", val)
	}
	if q.Size() != 1 {
		t.Errorf("Expected size 1 after dequeue, got %d", q.Size())
	}

	val, ok = q.Dequeue()
	if !ok || val != 20 {
		t.Errorf("Expected to dequeue 20, got %d", val)
	}
	if !q.IsEmpty() {
		t.Error("Queue should be empty after dequeuing all items")
	}

	_, ok = q.Dequeue()
	if ok {
		t.Error("Dequeue on empty queue should return false")
	}
}

func TestPeek(t *testing.T) {
	q := New[float64]()
	q.Enqueue(3.14)

	val, ok := q.Peek()
	if !ok || val != 3.14 {
		t.Errorf("Peek should return 3.14, got %v", val)
	}
	if q.Size() != 1 {
		t.Error("Peek should not remove the item from the queue")
	}

	q.Dequeue()
	_, ok = q.Peek()
	if ok {
		t.Error("Peek on empty queue should return false")
	}
}

func TestClearQueue(t *testing.T) {
	q := New[int]()
	q.Enqueue(1)
	q.Clear()
	if !q.IsEmpty() {
		t.Error("Queue should be empty after Clear")
	}
}

func TestQueueConcurrency(t *testing.T) {
	q := New[int]()
	var wg sync.WaitGroup

	// Concurrent enqueues
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(i)
		}(i)
	}
	wg.Wait()

	if q.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent enqueues, got %d", q.Size())
	}

	// Concurrent dequeues
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.Dequeue()
		}()
	}
	wg.Wait()

	if q.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent dequeues, got %d", q.Size())
	}
}
