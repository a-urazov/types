package sort

import (
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

// TestSlice тестирует функцию Slice с различными сценариями данных
func TestSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Single element",
			input:    []int{5},
			expected: []int{5},
		},
		{
			name:     "Already sorted",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Reverse sorted",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Random order",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
			expected: []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9},
		},
		{
			name:     "Duplicate elements",
			input:    []int{4, 2, 4, 2, 4, 2},
			expected: []int{2, 2, 2, 4, 4, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создать копию входных данных, чтобы не изменять оригинал
			data := make([]int, len(tc.input))
			copy(data, tc.input)

			// Отсортировать, используя нашу функцию
			Slice(data, func(a, b int) bool { return a < b })

			// Проверить, совпадает ли результат с ожидаемым
			if !reflect.DeepEqual(data, tc.expected) {
				t.Errorf("Slice() = %v, want %v", data, tc.expected)
			}
		})
	}
}

// TestSliceString тестирует функцию Slice с данными типа string
func TestSliceString(t *testing.T) {
	input := []string{"banana", "apple", "cherry", "date"}
	expected := []string{"apple", "banana", "cherry", "date"}

	data := make([]string, len(input))
	copy(data, input)

	Slice(data, func(a, b string) bool { return a < b })

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Slice() = %v, want %v", data, expected)
	}
}

// BenchmarkSlice тестирует производительность функции Slice с различными размерами данных
func BenchmarkSlice(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			data := make([]int, size)

			for i := 0; i < b.N; i++ {
				// Переинициализировать данные случайными значениями для каждой итерации
				rand.Seed(time.Now().UnixNano())
				for j := 0; j < size; j++ {
					data[j] = rand.Intn(size * 10)
				}

				b.StartTimer()
				Slice(data, func(a, b int) bool { return a < b })
				b.StopTimer()
			}
		})
	}
}

// BenchmarkSliceSorted тестирует производительность функции Slice с уже отсортированными данными
func BenchmarkSliceSorted(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run("Sorted_"+strconv.Itoa(size), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				Slice(data, func(a, b int) bool { return a < b })
				b.StopTimer()
			}
		})
	}
}

// BenchmarkSliceReverse тестирует производительность функции Slice с обратно отсортированными данными
func BenchmarkSliceReverse(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run("Reverse_"+strconv.Itoa(size), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = size - i
			}

			for i := 0; i < b.N; i++ {
				b.StartTimer()
				Slice(data, func(a, b int) bool { return a < b })
				b.StopTimer()
			}
		})
	}
}

// BenchmarkSliceBuiltIn сравнивает нашу реализацию со встроенной сортировкой Go
func BenchmarkSliceBuiltIn(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run("BuiltIn_"+strconv.Itoa(size), func(b *testing.B) {
			data := make([]int, size)

			for i := 0; i < b.N; i++ {
				// Переинициализировать данные случайными значениями для каждой итерации
				rand.Seed(time.Now().UnixNano())
				for j := 0; j < size; j++ {
					data[j] = rand.Intn(size * 10)
				}

				b.StartTimer()
				sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
				b.StopTimer()
			}
		})
	}
}
