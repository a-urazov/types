package heap

import (
	"testing"
)

func TestMinHeap(t *testing.T) {
	minHeap := New[int](func(a, b int) bool { return a < b })
	minHeap.Push(3)
	minHeap.Push(1)
	minHeap.Push(4)
	minHeap.Push(2)

	if size := minHeap.Size(); size != 4 {
		t.Errorf("Ожидался размер 4, получено %d", size)
	}

	expected := []int{1, 2, 3, 4}
	for _, val := range expected {
		popVal, ok := minHeap.Pop()
		if !ok {
			t.Fatalf("Ожидалось извлечение значения, но куча пуста")
		}
		if popVal != val {
			t.Errorf("Ожидалось извлечение %d, получено %d", val, popVal)
		}
	}

	if !minHeap.IsEmpty() {
		t.Errorf("Ожидалось, что куча будет пуста")
	}
}

func TestMaxHeap(t *testing.T) {
	maxHeap := New[int](func(a, b int) bool { return a > b })
	maxHeap.Push(3)
	maxHeap.Push(1)
	maxHeap.Push(4)
	maxHeap.Push(2)

	if size := maxHeap.Size(); size != 4 {
		t.Errorf("Ожидался размер 4, получено %d", size)
	}

	expected := []int{4, 3, 2, 1}
	for _, val := range expected {
		popVal, ok := maxHeap.Pop()
		if !ok {
			t.Fatalf("Ожидалось извлечение значения, но куча пуста")
		}
		if popVal != val {
			t.Errorf("Ожидалось извлечение %d, получено %d", val, popVal)
		}
	}

	if !maxHeap.IsEmpty() {
		t.Errorf("Ожидалось, что куча будет пуста")
	}
}

func TestPeek(t *testing.T) {
	h := New[string](func(a, b string) bool { return a < b })
	h.Push("b")
	h.Push("a")
	h.Push("c")

	val, ok := h.Peek()
	if !ok || val != "a" {
		t.Errorf("Ожидалось значение 'a', получено '%s'", val)
	}

	if size := h.Size(); size != 3 {
		t.Errorf("Ожидался размер 3 после Peek, получено %d", size)
	}
}

func TestPopEmpty(t *testing.T) {
	h := New[int](func(a, b int) bool { return a < b })
	_, ok := h.Pop()
	if ok {
		t.Errorf("Ожидалось, что Pop для пустой кучи вернет false")
	}
}
