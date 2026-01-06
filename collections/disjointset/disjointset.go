package disjointset

import (
	"sync"
)

// Element represents an element in the disjoint set with its parent and rank
type Element[T comparable] struct {
	parent T
	rank   int
}

// DisjointSet represents a thread-safe disjoint set (Union-Find) data structure
type DisjointSet[T comparable] struct {
	elements map[T]*Element[T]
	mutex    sync.Mutex
}

// New creates a new empty DisjointSet
func New[T comparable]() *DisjointSet[T] {
	return &DisjointSet[T]{
		elements: make(map[T]*Element[T]),
	}
}

// MakeSet creates a new set containing only the given element
// If the element already exists, it does nothing
func (ds *DisjointSet[T]) MakeSet(element T) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if _, exists := ds.elements[element]; !exists {
		ds.elements[element] = &Element[T]{
			parent: element,
			rank:   0,
		}
	}
}

// Find returns the representative (root) of the set containing the element
// Returns the root and true if the element exists, false otherwise
// Uses path compression for optimization
func (ds *DisjointSet[T]) Find(element T) (T, bool) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	return ds.findWithCompression(element)
}

// findWithCompression finds the root with path compression
func (ds *DisjointSet[T]) findWithCompression(element T) (T, bool) {
	elem, exists := ds.elements[element]
	if !exists {
		var zero T
		return zero, false
	}

	if elem.parent != element {
		root, _ := ds.findWithCompression(elem.parent)
		ds.elements[element].parent = root
		return root, true
	}

	return element, true
}

// Union merges the sets containing element1 and element2
// Returns true if the sets were merged, false if they were already in the same set
// Uses union by rank for optimization
func (ds *DisjointSet[T]) Union(element1, element2 T) bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	// Ensure both elements exist
	if _, exists := ds.elements[element1]; !exists {
		ds.elements[element1] = &Element[T]{parent: element1, rank: 0}
	}
	if _, exists := ds.elements[element2]; !exists {
		ds.elements[element2] = &Element[T]{parent: element2, rank: 0}
	}

	root1, _ := ds.findWithoutCompression(element1)
	root2, _ := ds.findWithoutCompression(element2)

	if root1 == root2 {
		return false // Already in the same set
	}

	// Union by rank
	elem1 := ds.elements[root1]
	elem2 := ds.elements[root2]

	if elem1.rank < elem2.rank {
		elem1.parent = root2
	} else if elem1.rank > elem2.rank {
		elem2.parent = root1
	} else {
		elem2.parent = root1
		elem1.rank++
	}

	return true
}

// findWithoutCompression finds the root without path compression
// Used internally when we already have a lock
func (ds *DisjointSet[T]) findWithoutCompression(element T) (T, bool) {
	current := element
	for {
		elem, exists := ds.elements[current]
		if !exists {
			var zero T
			return zero, false
		}
		if elem.parent == current {
			return current, true
		}
		current = elem.parent
	}
}

// Connected checks if two elements are in the same set
func (ds *DisjointSet[T]) Connected(element1, element2 T) bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	root1, exists1 := ds.findWithoutCompression(element1)
	root2, exists2 := ds.findWithoutCompression(element2)
	return exists1 && exists2 && root1 == root2
}

// Size returns the total number of elements in the disjoint set
func (ds *DisjointSet[T]) Size() int {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	return len(ds.elements)
}

// IsEmpty returns true if the disjoint set is empty
func (ds *DisjointSet[T]) IsEmpty() bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	return len(ds.elements) == 0
}

// Clear removes all elements from the disjoint set
func (ds *DisjointSet[T]) Clear() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	ds.elements = make(map[T]*Element[T])
}

// Sets returns a slice of all disjoint sets, where each set is represented as a slice of elements
func (ds *DisjointSet[T]) Sets() [][]T {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	// Group elements by their root
	setsMap := make(map[T][]T)
	for element := range ds.elements {
		root, _ := ds.findWithoutCompression(element)
		setsMap[root] = append(setsMap[root], element)
	}

	// Convert to slice of slices
	sets := make([][]T, 0, len(setsMap))
	for _, set := range setsMap {
		sets = append(sets, set)
	}

	return sets
}

// SetCount returns the number of disjoint sets
func (ds *DisjointSet[T]) SetCount() int {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	roots := make(map[T]bool)
	for element := range ds.elements {
		root, _ := ds.findWithoutCompression(element)
		roots[root] = true
	}

	return len(roots)
}

// Elements returns a slice of all elements in the disjoint set
func (ds *DisjointSet[T]) Elements() []T {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	elements := make([]T, 0, len(ds.elements))
	for element := range ds.elements {
		elements = append(elements, element)
	}

	return elements
}

// ForEachSet iterates over each disjoint set, calling the function with all elements in the set
func (ds *DisjointSet[T]) ForEachSet(fn func(elements []T)) {
	sets := ds.Sets()
	for _, set := range sets {
		fn(set)
	}
}
