package streams

import (
	"sort"
)

func (s *stream[T]) Map(f func(T) T) Stream[T] {
	result := make([]T, len(s.elements))
	for i, v := range s.elements {
		result[i] = f(v)
	}
	return &stream[T]{elements: result}
}

func (s *stream[T]) Filter(f func(T) bool) Stream[T] {
	var result []T
	for _, v := range s.elements {
		if f(v) {
			result = append(result, v)
		}
	}
	return &stream[T]{elements: result}
}

func (s *stream[T]) Sorted(comparator func(T, T) int) Stream[T] {
	result := make([]T, len(s.elements))
	copy(result, s.elements)
	sort.Slice(result, func(i, j int) bool {
		return comparator(result[i], result[j]) < 0
	})
	return &stream[T]{elements: result}
}

func (s *stream[T]) Distinct() Stream[T] {
	if len(s.elements) == 0 {
		return &stream[T]{elements: []T{}}
	}

	// Use a map to store unique elements. This requires T to be comparable.
	// The zero-sized struct is used as the value to minimize memory usage.
	uniqueSet := make(map[any]struct{})
	var result []T

	for _, v := range s.elements {
		// We need to cast to `any` to use it as a map key since T is not guaranteed to be comparable at compile time.
		// This is a common workaround in Go for generic types.
		if _, ok := uniqueSet[v]; !ok {
			uniqueSet[v] = struct{}{}
			result = append(result, v)
		}
	}

	return &stream[T]{elements: result}
}

func (s *stream[T]) ForEach(f func(T)) {
	for _, v := range s.elements {
		f(v)
	}
}

func (s *stream[T]) Collect() []T {
	return s.elements
}

func (s *stream[T]) Reduce(f func(T, T) T) (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	acc := s.elements[0]
	for i := 1; i < len(s.elements); i++ {
		acc = f(acc, s.elements[i])
	}
	return acc, true
}

func (s *stream[T]) Count() int {
	return len(s.elements)
}

func (s *stream[T]) AnyMatch(f func(T) bool) bool {
	for _, v := range s.elements {
		if f(v) {
			return true
		}
	}
	return false
}

func (s *stream[T]) AllMatch(f func(T) bool) bool {
	for _, v := range s.elements {
		if !f(v) {
			return false
		}
	}
	return true
}

func (s *stream[T]) NoneMatch(f func(T) bool) bool {
	return !s.AnyMatch(f)
}

func (s *stream[T]) FindFirst() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[0], true
}
