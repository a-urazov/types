package skip

import (
	"cmp"
	"math/rand"
	"time"
)

const (
	maxLevel = 16   // Maximum number of levels
	p        = 0.25 // Probability for level promotion
)

// Node represents a node in the skip list.
type Node[T cmp.Ordered] struct {
	Value   T
	forward []*Node[T] // Array of pointers to next nodes at each level
}

// NewNode creates and returns a new node with the given value and level.
func NewNode[T cmp.Ordered](value T, level int) *Node[T] {
	return &Node[T]{
		Value:   value,
		forward: make([]*Node[T], level+1), // Level is 0-indexed, so size is level+1
	}
}

// List represents a skip list data structure.
type List[T cmp.Ordered] struct {
	header *Node[T]   // Header node with -inf value
	level  int        // Current maximum level (0-indexed)
	rand   *rand.Rand // Local random number generator
}

// New creates and returns a new empty SkipList.
func New[T cmp.Ordered]() *List[T] {
	var zeroValue T // Use zero value of type T
	header := NewNode(zeroValue, maxLevel)
	src := rand.NewSource(time.Now().UnixNano())
	return &List[T]{
		header: header,
		level:  0,
		rand:   rand.New(src),
	}
}

// randomLevel generates a random level for a new node.
func (sl *List[T]) randomLevel() int {
	level := 0
	for sl.rand.Float64() < p && level < maxLevel {
		level++
	}
	return level
}

// Insert adds a value to the skip list.
func (sl *List[T]) Insert(value T) {
	update := make([]*Node[T], maxLevel+1)

	x := sl.header
	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Value < value {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	// If the value already exists, we can choose to update it or ignore.
	// For this implementation, we'll ignore duplicates.
	if x != nil && x.Value == value {
		return
	}

	// Generate a random level for the new node
	newLevel := sl.randomLevel()

	if newLevel > sl.level {
		for i := sl.level + 1; i <= newLevel; i++ {
			update[i] = sl.header
		}
		sl.level = newLevel
	}

	// Создать новый узел
	newNode := NewNode(value, newLevel)

	// Insert the new node
	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

// Search returns true if the value exists in the skip list, false otherwise.
func (sl *List[T]) Search(value T) bool {
	x := sl.header
	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Value < value {
			x = x.forward[i]
		}
	}

	x = x.forward[0]
	return x != nil && x.Value == value
}

// Delete removes a value from the skip list.
// It returns true if the value was found and deleted, false otherwise.
func (sl *List[T]) Delete(value T) bool {
	update := make([]*Node[T], maxLevel+1)

	x := sl.header
	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Value < value {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	if x == nil || x.Value != value {
		return false
	}

	// Remove the node
	for i := 0; i <= sl.level; i++ {
		if update[i].forward[i] != x {
			break
		}
		update[i].forward[i] = x.forward[i]
	}

	// Update the current level if the top levels are empty
	for sl.level > 0 && sl.header.forward[sl.level] == nil {
		sl.level--
	}

	return true
}

// Size returns the number of elements in the skip list.
// This is an O(n) operation as we need to traverse the bottom level.
func (sl *List[T]) Size() int {
	count := 0
	x := sl.header.forward[0]
	for x != nil {
		count++
		x = x.forward[0]
	}
	return count
}

// IsEmpty returns true if the skip list is empty, false otherwise.
func (sl *List[T]) IsEmpty() bool {
	return sl.header.forward[0] == nil
}
