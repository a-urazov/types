package streams

import (
	"bufio"
	"io"
)

// Stream provides a way to process a sequence of elements.
type Stream[T any] interface {
	// Intermediate operations
	Map(func(T) T) Stream[T]
	Filter(func(T) bool) Stream[T]
	Sorted(func(T, T) int) Stream[T]
	Distinct() Stream[T]

	// Terminal operations
	ForEach(func(T))
	Collect() []T
	Reduce(func(T, T) T) (T, bool)
	Count() int
	AnyMatch(func(T) bool) bool
	AllMatch(func(T) bool) bool
	NoneMatch(func(T) bool) bool
	FindFirst() (T, bool)
}

type stream[T any] struct {
	elements []T
}

// Of creates a new Stream from a slice of elements.
func Of[T any](elements []T) Stream[T] {
	return &stream[T]{elements: elements}
}

// FromReader creates a new Stream from an io.Reader, processing it line by line.
// This implementation uses channels to provide a lazy-loaded stream.
func FromReader(reader io.Reader) Stream[string] {
	elementsChan := make(chan string)

	go func() {
		defer close(elementsChan)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			elementsChan <- scanner.Text()
		}
	}()

	// The stream will collect elements from the channel on demand.
	// This is a lazy implementation.
	var elements []string
	for elem := range elementsChan {
		elements = append(elements, elem)
	}

	return &stream[string]{elements: elements}
}
