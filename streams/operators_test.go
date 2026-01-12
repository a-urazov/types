package streams

import (
	"reflect"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	data := []int{1, 2, 3}
	s := Of(data)
	result := s.Map(func(i int) int { return i * 2 }).Collect()
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map() got = %v, want %v", result, expected)
	}
}

func TestFilter(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	s := Of(data)
	result := s.Filter(func(i int) bool { return i%2 == 0 }).Collect()
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter() got = %v, want %v", result, expected)
	}
}

func TestSorted(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	s := Of(data)
	result := s.Sorted(func(a, b int) int { return a - b }).Collect()
	expected := []int{1, 1, 3, 4, 5, 9}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Sorted() got = %v, want %v", result, expected)
	}
}

func TestReduce(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	s := Of(data)
	result, ok := s.Reduce(func(a, b int) int { return a + b })
	if !ok {
		t.Error("Reduce() returned not ok")
	}
	expected := 15
	if result != expected {
		t.Errorf("Reduce() got = %v, want %v", result, expected)
	}
}

func TestCount(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	s := Of(data)
	result := s.Count()
	expected := 5
	if result != expected {
		t.Errorf("Count() got = %v, want %v", result, expected)
	}
}

func TestAnyMatch(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	s := Of(data)
	if !s.AnyMatch(func(i int) bool { return i == 3 }) {
		t.Error("AnyMatch() should be true")
	}
	if s.AnyMatch(func(i int) bool { return i == 10 }) {
		t.Error("AnyMatch() should be false")
	}
}

func TestAllMatch(t *testing.T) {
	data := []int{2, 4, 6, 8}
	s := Of(data)
	if !s.AllMatch(func(i int) bool { return i%2 == 0 }) {
		t.Error("AllMatch() should be true")
	}
	data2 := []int{2, 4, 5, 8}
	s2 := Of(data2)
	if s2.AllMatch(func(i int) bool { return i%2 == 0 }) {
		t.Error("AllMatch() should be false")
	}
}

func TestNoneMatch(t *testing.T) {
	data := []int{1, 3, 5, 7}
	s := Of(data)
	if !s.NoneMatch(func(i int) bool { return i%2 == 0 }) {
		t.Error("NoneMatch() should be true")
	}
	data2 := []int{1, 3, 4, 7}
	s2 := Of(data2)
	if s2.NoneMatch(func(i int) bool { return i%2 == 0 }) {
		t.Error("NoneMatch() should be false")
	}
}

func TestFindFirst(t *testing.T) {
	data := []int{1, 2, 3}
	s := Of(data)
	result, ok := s.FindFirst()
	if !ok {
		t.Error("FindFirst() returned not ok")
	}
	if result != 1 {
		t.Errorf("FindFirst() got = %v, want %v", result, 1)
	}

	var emptyData []int
	sEmpty := Of(emptyData)
	_, ok = sEmpty.FindFirst()
	if ok {
		t.Error("FindFirst() on empty stream should return not ok")
	}
}

func TestStream_Distinct(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		expected []int
	}{
		{
			name:     "no duplicates",
			elements: []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "with duplicates",
			elements: []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all duplicates",
			elements: []int{1, 1, 1, 1, 1},
			expected: []int{1},
		},
		{
			name:     "empty stream",
			elements: []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Of(tt.elements)
			result := s.Distinct().Collect()

			// Since the order is not guaranteed, we should sort both slices before comparing
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
