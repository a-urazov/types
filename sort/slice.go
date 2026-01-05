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

	const threshold = 20 // Увеличенный порог для лучшей производительности с большими срезами
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

	// Определить оптимальное количество горутин на основе размера среза и доступных процессоров
	depthLimit := 0
	for n := len(s); n > 20; n >>= 1 {
		depthLimit++
	}

	// Использовать подход пула рабочих для ограничения создания горутин
	parallelSortDepth(s, less, depthLimit)
}

func parallelSortDepth[T any](s []T, less func(a, b T) bool, depth int) {
	if len(s) <= 1 {
		return
	}

	const threshold = 20
	if len(s) <= threshold || depth == 0 {
		insertionSort(s, less)
		return
	}

	mid := len(s) / 2
	left := s[:mid]
	right := s[mid:]

	// Ограничить создание горутин, используя контроль рекурсии на основе глубины
	if depth > 1 {
		// Сортировать левую и правую половины параллельно
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			parallelSortDepth(left, less, depth-1)
		}()
		go func() {
			defer wg.Done()
			parallelSortDepth(right, less, depth-1)
		}()

		// Дождаться, пока обе половины будут отсортированы
		wg.Wait()
	} else {
		// Для более глубоких уровней сортировать последовательно, чтобы избежать чрезмерного создания горутин
		parallelSortDepth(left, less, depth-1)
		parallelSortDepth(right, less, depth-1)
	}

	// Объединить отсортированные половины
	merge(s, left, right, less)
}

func merge[T any](s, left, right []T, less func(a, b T) bool) {
	i, j, k := 0, 0, 0

	// Оптимизировать основной цикл слияния
	for i < len(left) && j < len(right) {
		if !less(right[j], left[i]) { // Использовать !less(right, left) вместо less(left, right) для сохранения стабильности
			s[k] = left[i]
			i++
		} else {
			s[k] = right[j]
			j++
		}
		k++
	}

	// Скопировать оставшиеся элементы - использовать copy для лучшей производительности
	copy(s[k:], left[i:])
	copy(s[k+len(left)-i:], right[j:])
}
