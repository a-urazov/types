package priority

import (
	"sync"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	pq := New[string, int]()
	if !pq.IsEmpty() || pq.Size() != 0 {
		t.Error("New priority queue should be empty")
	}
}

func TestEnqueueAndDequeue(t *testing.T) {
	pq := New[string, int]()
	pq.Enqueue("low", 1)
	pq.Enqueue("high", 10)
	pq.Enqueue("medium", 5)

	if pq.Size() != 3 {
		t.Errorf("Expected size 3, got %d", pq.Size())
	}

	item, ok := pq.Dequeue()
	if !ok || item != "high" {
		t.Errorf("Expected to dequeue 'high', got '%s'", item)
	}

	item, ok = pq.Dequeue()
	if !ok || item != "medium" {
		t.Errorf("Expected to dequeue 'medium', got '%s'", item)
	}

	item, ok = pq.Dequeue()
	if !ok || item != "low" {
		t.Errorf("Expected to dequeue 'low', got '%s'", item)
	}

	if !pq.IsEmpty() {
		t.Error("Queue should be empty")
	}
}

func TestDequeueFromEmpty(t *testing.T) {
	pq := New[int, int]()
	_, ok := pq.Dequeue()
	if ok {
		t.Error("Dequeue on empty queue should return false")
	}
}

func TestPriorityQueueConcurrency(t *testing.T) {
	pq := New[int, int]()
	var wg sync.WaitGroup

	// Concurrent enqueues
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pq.Enqueue(i, i) // Value and priority are the same
		}(i)
	}
	wg.Wait()

	if pq.Size() != 100 {
		t.Errorf("Expected size 100, got %d", pq.Size())
	}

	// Concurrent dequeues
	results := make(chan int, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if val, ok := pq.Dequeue(); ok {
				results <- val
			}
		}()
	}
	wg.Wait()
	close(results)

	if pq.Size() != 50 {
		t.Errorf("Expected size 50 after dequeues, got %d", pq.Size())
	}

	// Check that we dequeued the highest priority items
	for val := range results {
		if val < 50 {
			t.Errorf("Dequeued a low priority item %d, expected items > 49", val)
		}
	}
}
