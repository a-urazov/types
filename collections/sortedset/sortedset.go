package sortedset

import (
	"sync"
)

// Node represents a node in the red-black tree
type Node[T any] struct {
	item   T
	color  bool // true = red, false = black
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]
}

// SortedSet represents a sorted set of unique elements
type SortedSet[T any] struct {
	root    *Node[T]
	size    int
	compare func(a, b T) bool // less function: returns true if a < b
	mutex   sync.RWMutex
}

// New creates a new SortedSet with the given comparison function
func New[T any](compare func(a, b T) bool) *SortedSet[T] {
	return &SortedSet[T]{
		compare: compare,
	}
}

// Add inserts an element into the set
// Returns true if the element was added, false if it already exists
func (s *SortedSet[T]) Add(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.root == nil {
		s.root = &Node[T]{item: item, color: false} // root is always black
		s.size = 1
		return true
	}

	added := s.insert(s.root, item)
	if added {
		s.size++
	}
	return added
}

// insert recursively inserts an item into the tree
func (s *SortedSet[T]) insert(node *Node[T], item T) bool {
	if s.compare(item, node.item) {
		// item < node.item, go left
		if node.left == nil {
			node.left = &Node[T]{item: item, color: true, parent: node}
			s.fixInsert(node.left)
			return true
		}
		return s.insert(node.left, item)
	} else if s.compare(node.item, item) {
		// node.item < item, go right
		if node.right == nil {
			node.right = &Node[T]{item: item, color: true, parent: node}
			s.fixInsert(node.right)
			return true
		}
		return s.insert(node.right, item)
	} else {
		// item == node.item, already exists
		return false
	}
}

// fixInsert fixes the red-black tree properties after insertion
func (s *SortedSet[T]) fixInsert(node *Node[T]) {
	for node != s.root && node.parent.color {
		if node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle != nil && uncle.color {
				// Case 1: Uncle is red
				node.parent.color = false
				uncle.color = false
				node.parent.parent.color = true
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					// Case 2: Uncle is black and node is right child
					node = node.parent
					s.rotateLeft(node)
				}
				// Case 3: Uncle is black and node is left child
				node.parent.color = false
				node.parent.parent.color = true
				s.rotateRight(node.parent.parent)
			}
		} else {
			// Mirror of above cases
			uncle := node.parent.parent.left
			if uncle != nil && uncle.color {
				node.parent.color = false
				uncle.color = false
				node.parent.parent.color = true
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					s.rotateRight(node)
				}
				node.parent.color = false
				node.parent.parent.color = true
				s.rotateLeft(node.parent.parent)
			}
		}
	}
	s.root.color = false // root is always black
}

// rotateLeft performs a left rotation
func (s *SortedSet[T]) rotateLeft(x *Node[T]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		s.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

// rotateRight performs a right rotation
func (s *SortedSet[T]) rotateRight(y *Node[T]) {
	x := y.left
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == nil {
		s.root = x
	} else if y == y.parent.right {
		y.parent.right = x
	} else {
		y.parent.left = x
	}
	x.right = y
	y.parent = x
}

// Contains checks if an element exists in the set
func (s *SortedSet[T]) Contains(item T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.search(s.root, item) != nil
}

// search finds a node with the given item
func (s *SortedSet[T]) search(node *Node[T], item T) *Node[T] {
	if node == nil {
		return nil
	}

	if s.compare(item, node.item) {
		return s.search(node.left, item)
	} else if s.compare(node.item, item) {
		return s.search(node.right, item)
	} else {
		return node
	}
}

// Remove deletes an element from the set
// Returns true if the element was removed, false if it didn't exist
func (s *SortedSet[T]) Remove(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	node := s.search(s.root, item)
	if node == nil {
		return false
	}

	s.deleteNode(node)
	s.size--
	return true
}

// deleteNode removes a node from the tree
func (s *SortedSet[T]) deleteNode(z *Node[T]) {
	var y, x *Node[T]

	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = s.successor(z)
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		s.root = x
	} else {
		if y == y.parent.left {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
	}

	if y != z {
		z.item = y.item
	}

	if !y.color { // if y is black, fix the tree
		s.fixDelete(x, y.parent)
	}
}

// successor finds the in-order successor of a node
func (s *SortedSet[T]) successor(node *Node[T]) *Node[T] {
	if node.right != nil {
		return s.minimum(node.right)
	}

	y := node.parent
	for y != nil && node == y.right {
		node = y
		y = y.parent
	}
	return y
}

// minimum finds the node with minimum value in subtree
func (s *SortedSet[T]) minimum(node *Node[T]) *Node[T] {
	for node.left != nil {
		node = node.left
	}
	return node
}

// fixDelete fixes the red-black tree properties after deletion
func (s *SortedSet[T]) fixDelete(x *Node[T], parent *Node[T]) {
	for x != s.root && (x == nil || !x.color) {
		if x == parent.left {
			w := parent.right
			if w != nil && w.color {
				w.color = false
				parent.color = true
				s.rotateLeft(parent)
				w = parent.right
			}
			if w != nil && (w.left == nil || !w.left.color) && (w.right == nil || !w.right.color) {
				w.color = true
				x = parent
				parent = x.parent
			} else {
				if w.right == nil || !w.right.color {
					if w.left != nil {
						w.left.color = false
					}
					w.color = true
					s.rotateRight(w)
					w = parent.right
				}
				w.color = parent.color
				parent.color = false
				if w.right != nil {
					w.right.color = false
				}
				s.rotateLeft(parent)
				x = s.root
			}
		} else {
			// Mirror of above cases
			w := parent.left
			if w != nil && w.color {
				w.color = false
				parent.color = true
				s.rotateRight(parent)
				w = parent.left
			}
			if w != nil && (w.right == nil || !w.right.color) && (w.left == nil || !w.left.color) {
				w.color = true
				x = parent
				parent = x.parent
			} else {
				if w.left == nil || !w.left.color {
					if w.right != nil {
						w.right.color = false
					}
					w.color = true
					s.rotateLeft(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = false
				if w.left != nil {
					w.left.color = false
				}
				s.rotateRight(parent)
				x = s.root
			}
		}
	}
	if x != nil {
		x.color = false
	}
}

// Size returns the number of elements in the set
func (s *SortedSet[T]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.size
}

// IsEmpty returns true if the set is empty
func (s *SortedSet[T]) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.size == 0
}

// Clear removes all elements from the set
func (s *SortedSet[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.root = nil
	s.size = 0
}

// ToSlice returns a slice containing all elements in sorted order
func (s *SortedSet[T]) ToSlice() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]T, 0, s.size)
	s.inOrderTraversal(s.root, &result)
	return result
}

// inOrderTraversal performs in-order traversal to collect elements
func (s *SortedSet[T]) inOrderTraversal(node *Node[T], result *[]T) {
	if node != nil {
		s.inOrderTraversal(node.left, result)
		*result = append(*result, node.item)
		s.inOrderTraversal(node.right, result)
	}
}

// ForEach iterates over all elements in sorted order
func (s *SortedSet[T]) ForEach(action func(item T)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.inOrderForEach(s.root, action)
}

// inOrderForEach performs in-order traversal with action
func (s *SortedSet[T]) inOrderForEach(node *Node[T], action func(item T)) {
	if node != nil {
		s.inOrderForEach(node.left, action)
		action(node.item)
		s.inOrderForEach(node.right, action)
	}
}

// First returns the smallest element in the set
func (s *SortedSet[T]) First() (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.root == nil {
		var zero T
		return zero, false
	}

	node := s.minimum(s.root)
	return node.item, true
}

// Last returns the largest element in the set
func (s *SortedSet[T]) Last() (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.root == nil {
		var zero T
		return zero, false
	}

	node := s.maximum(s.root)
	return node.item, true
}

// maximum finds the node with maximum value in subtree
func (s *SortedSet[T]) maximum(node *Node[T]) *Node[T] {
	for node.right != nil {
		node = node.right
	}
	return node
}

// Ceiling returns the smallest element >= given item
func (s *SortedSet[T]) Ceiling(item T) (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.root == nil {
		var zero T
		return zero, false
	}

	result := s.ceiling(s.root, item)
	if result == nil {
		var zero T
		return zero, false
	}
	return result.item, true
}

// ceiling finds the ceiling node
func (s *SortedSet[T]) ceiling(node *Node[T], item T) *Node[T] {
	if node == nil {
		return nil
	}

	if !s.compare(node.item, item) { // node.item >= item
		leftCeiling := s.ceiling(node.left, item)
		if leftCeiling != nil {
			return leftCeiling
		}
		return node
	}

	return s.ceiling(node.right, item)
}

// Floor returns the largest element <= given item
func (s *SortedSet[T]) Floor(item T) (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.root == nil {
		var zero T
		return zero, false
	}

	result := s.floor(s.root, item)
	if result == nil {
		var zero T
		return zero, false
	}
	return result.item, true
}

// floor finds the floor node
func (s *SortedSet[T]) floor(node *Node[T], item T) *Node[T] {
	if node == nil {
		return nil
	}

	if !s.compare(item, node.item) { // item >= node.item
		rightFloor := s.floor(node.right, item)
		if rightFloor != nil {
			return rightFloor
		}
		return node
	}

	return s.floor(node.left, item)
}
