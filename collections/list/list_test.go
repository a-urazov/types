package list

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	l := New[int]()
	if l.Size() != 0 {
		t.Errorf("New list should be empty, but has size %d", l.Size())
	}
}

func TestAdd(t *testing.T) {
	l := New[string]()
	l.Add("a")
	l.Add("b")
	if l.Size() != 2 {
		t.Errorf("Expected size 2, got %d", l.Size())
	}
	val, ok := l.Get(0)
	if !ok || val != "a" {
		t.Errorf("Expected to get 'a' at index 0, got '%v'", val)
	}
	val, ok = l.Get(1)
	if !ok || val != "b" {
		t.Errorf("Expected to get 'b' at index 1, got '%v'", val)
	}
}

func TestRemove(t *testing.T) {
	l := New[int]()
	l.Add(1)
	l.Add(2)
	l.Add(3)

	if !l.Remove(2) {
		t.Error("Expected to remove 2, but it failed")
	}
	if l.Size() != 2 {
		t.Errorf("Expected size 2 after removing, got %d", l.Size())
	}
	if l.Contains(2) {
		t.Error("List should not contain 2 after removal")
	}

	if l.Remove(4) {
		t.Error("Should not be able to remove non-existent item")
	}
}

func TestGet(t *testing.T) {
	l := New[float64]()
	l.Add(1.1)
	val, ok := l.Get(0)
	if !ok || val != 1.1 {
		t.Errorf("Expected to get 1.1, got %v", val)
	}
	_, ok = l.Get(1)
	if ok {
		t.Error("Getting from out of bounds index should fail")
	}
	_, ok = l.Get(-1)
	if ok {
		t.Error("Getting from negative index should fail")
	}
}

func TestSize(t *testing.T) {
	l := New[int]()
	if l.Size() != 0 {
		t.Errorf("Empty list should have size 0, got %d", l.Size())
	}
	l.Add(10)
	if l.Size() != 1 {
		t.Errorf("List should have size 1, got %d", l.Size())
	}
}

func TestIndexOf(t *testing.T) {
	l := New[string]()
	l.Add("apple")
	l.Add("banana")
	l.Add("apple")

	if i := l.IndexOf("banana"); i != 1 {
		t.Errorf("Expected index of 'banana' to be 1, got %d", i)
	}
	if i := l.IndexOf("apple"); i != 0 {
		t.Errorf("Expected index of 'apple' to be 0, got %d", i)
	}
	if i := l.IndexOf("cherry"); i != -1 {
		t.Errorf("Expected index of 'cherry' to be -1, got %d", i)
	}
}

func TestContains(t *testing.T) {
	l := New[int]()
	l.Add(100)
	if !l.Contains(100) {
		t.Error("List should contain 100")
	}
	if l.Contains(200) {
		t.Error("List should not contain 200")
	}
}

func TestClear(t *testing.T) {
	l := New[int]()
	l.Add(1)
	l.Add(2)
	l.Clear()
	if l.Size() != 0 {
		t.Errorf("List should be empty after Clear, but size is %d", l.Size())
	}
	if l.Contains(1) {
		t.Error("List should not contain 1 after Clear")
	}
}

func TestInsert(t *testing.T) {
	l := New[string]()
	l.Add("a")
	l.Add("c")

	if !l.Insert(1, "b") {
		t.Error("Failed to insert 'b'")
	}
	if l.Size() != 3 {
		t.Errorf("Expected size 3 after insert, got %d", l.Size())
	}
	val, _ := l.Get(1)
	if val != "b" {
		t.Errorf("Expected 'b' at index 1, got '%s'", val)
	}

	if l.Insert(4, "d") {
		t.Error("Should not be able to insert at out of bounds index")
	}
	if l.Insert(-1, "d") {
		t.Error("Should not be able to insert at negative index")
	}
}

func TestRemoveAt(t *testing.T) {
	l := New[int]()
	l.Add(10)
	l.Add(20)
	l.Add(30)

	if !l.RemoveAt(1) {
		t.Error("Failed to remove at index 1")
	}
	if l.Size() != 2 {
		t.Errorf("Expected size 2 after RemoveAt, got %d", l.Size())
	}
	if l.Contains(20) {
		t.Error("List should not contain 20 after RemoveAt")
	}
	val, _ := l.Get(1)
	if val != 30 {
		t.Errorf("Expected 30 at index 1, got %d", val)
	}

	if l.RemoveAt(2) {
		t.Error("Should not be able to remove at out of bounds index")
	}
}

func TestForEach(t *testing.T) {
	l := New[int]()
	l.Add(1)
	l.Add(2)
	l.Add(3)

	sum := 0
	l.ForEach(func(item int) {
		sum += item
	})

	if sum != 6 {
		t.Errorf("Expected sum 6 from ForEach, got %d", sum)
	}
}

func TestToArray(t *testing.T) {
	l := New[string]()
	l.Add("x")
	l.Add("y")

	arr := l.ToArray()
	if len(arr) != 2 {
		t.Fatalf("Expected array of length 2, got %d", len(arr))
	}
	if arr[0] != "x" || arr[1] != "y" {
		t.Errorf("Array content mismatch: got %v", arr)
	}

	// Modify original list, array should not change
	l.Add("z")
	if len(arr) != 2 {
		t.Error("Array should not be modified after list changes")
	}
}

func TestConcurrency(t *testing.T) {
	l := New[string]()
	var wg sync.WaitGroup

	// Test concurrent additions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Add("item_" + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	if l.Size() != 100 {
		t.Errorf("Expected size 100 after concurrent adds, got %d", l.Size())
	}

	// Test concurrent reads and removals
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			l.Contains("item_" + strconv.Itoa(i))
		}(i)
		go func(i int) {
			defer wg.Done()
			l.Remove("item_" + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	if l.Size() != 50 {
		t.Errorf("Expected size 50 after concurrent removals, got %d", l.Size())
	}

	// Final check for any race conditions with the race detector (go test -race)
	fmt.Println("Concurrency test finished, run with -race flag to check for data races.")
}
