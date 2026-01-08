package sort

import (
	"runtime"
	"sync"
)

// MaxGoroutines определяет максимальное количество горутин для параллельной сортировки
var MaxGoroutines = runtime.GOMAXPROCS(0)

func Slice[T any](s []T, less func(a, b T) bool) {
	hybridSort(s, less)
}

func hybridSort[T any](s []T, less func(a, b T) bool) {
	if len(s) <= 1 {
		return
	}

	const threshold = 16 // Оптимизированный порог для cache-friendly производительности
	if len(s) <= threshold {
		insertionSort(s, less)
		return
	}

	parallelSort(s, less)
}

func insertionSort[T any](s []T, less func(a, b T) bool) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && less(s[j], s[j-1]); j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func parallelSort[T any](s []T, less func(a, b T) bool) {
	if len(s) <= 1 {
		return
	}

	// Адаптивный расчет глубины параллелизма на основе размера среза и доступных ядер
	depthLimit := calculateDepthLimit(len(s))

	// Использовать параллельную интроспективную сортировку (Quicksort + Heapsort защита)
	parallelIntrosortDepth(s, less, depthLimit, 2*depthLimit)
}

func calculateDepthLimit(n int) int {
	// Базовая глубина на основе log2(n), используя порог 16 для Insertion Sort
	depth := 0
	for n > 16 {
		n >>= 1
		depth++
	}

	// Ограничиваем глубину на основе количества доступных ядер
	// maxDepth ≈ log2(GOMAXPROCS), чтобы не превышать горутины > GOMAXPROCS*2
	maxDepth := 0
	for cores := MaxGoroutines; cores > 1; cores >>= 1 {
		maxDepth++
	}

	if depth > maxDepth {
		depth = maxDepth
	}

	return depth
}

func parallelIntrosortDepth[T any](s []T, less func(a, b T) bool, depth, maxDepth int) {
	for len(s) > 16 {
		if depth == 0 {
			// Если достигли предела глубины, использовать Heapsort для избежания worst-case
			heapsort(s, less)
			return
		}

		// Выбрать опорный элемент (median-of-three для лучшей производительности)
		pivot := medianOfThree(s, less)

		// Разделить срез на три части: < pivot, == pivot, > pivot
		left, right := partition(s, less, pivot)

		// Рекурсивно сортировать меньшую половину параллельно, большую - последовательно
		if len(s[:left]) < len(s[right:]) {
			// Меньшая половина слева
			sortPartitions(s[:left], s[right:], less, depth, maxDepth)
		} else {
			// Меньшая половина справа
			sortPartitions(s[right:], s[:left], less, depth, maxDepth)
		}
	}

	// Для малых срезов (<=16) использовать Insertion Sort
	insertionSort(s, less)
}

func partition[T any](s []T, less func(a, b T) bool, pivot T) (left, right int) {
	// Трёхсторонее разделение (3-way partition) для обработки дубликатов
	l, m, r := 0, 0, len(s)

	for m < r {
		cmp := compare(s[m], pivot, less)
		switch cmp {
		case -1: // s[m] < pivot
			s[l], s[m] = s[m], s[l]
			l++
			m++
		case 0: // s[m] == pivot
			m++
		case 1: // s[m] > pivot
			r--
			s[m], s[r] = s[r], s[m]
		}
	}

	return l, r
}

func compare[T any](a, b T, less func(a, b T) bool) int {
	// Возвращает: -1 если a < b, 0 если a == b, 1 если a > b
	if less(a, b) {
		return -1
	}
	if less(b, a) {
		return 1
	}
	return 0
}

func medianOfThree[T any](s []T, less func(a, b T) bool) T {
	// Выбрать медиану из первого, среднего и последнего элементов
	a, b, c := s[0], s[len(s)/2], s[len(s)-1]

	if less(a, b) {
		if less(b, c) {
			return b // a < b < c
		} else if less(a, c) {
			return c // a < c < b
		} else {
			return a // c < a < b
		}
	} else {
		if less(a, c) {
			return a // b < a < c
		} else if less(b, c) {
			return c // b < c < a
		} else {
			return b // c < b < a
		}
	}
}

func heapsort[T any](s []T, less func(a, b T) bool) {
	// Построить max-heap (или min-heap в зависимости от less)
	n := len(s)
	for i := n/2 - 1; i >= 0; i-- {
		heapify(s, i, n, less)
	}

	// Извлечь элементы из heap один за другим
	for i := n - 1; i > 0; i-- {
		s[0], s[i] = s[i], s[0]
		heapify(s, 0, i, less)
	}
}

func heapify[T any](s []T, i, n int, less func(a, b T) bool) {
	// Восстановить свойство max-heap для поддерева с корнем в i
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && less(s[largest], s[left]) {
		largest = left
	}

	if right < n && less(s[largest], s[right]) {
		largest = right
	}

	if largest != i {
		s[i], s[largest] = s[largest], s[i]
		heapify(s, largest, n, less)
	}
}

func sortPartitions[T any](smaller, larger []T, less func(a, b T) bool, depth, maxDepth int) {
	if depth > 1 && len(smaller) > 16 {
		var wg sync.WaitGroup
		wg.Go(func() {
			parallelIntrosortDepth(smaller, less, depth-1, maxDepth)
		})
		parallelIntrosortDepth(larger, less, depth-1, maxDepth)
		wg.Wait()
	} else {
		parallelIntrosortDepth(smaller, less, depth-1, maxDepth)
		parallelIntrosortDepth(larger, less, depth-1, maxDepth)
	}
}
