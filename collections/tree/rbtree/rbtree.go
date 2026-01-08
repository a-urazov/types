package rbtree

import (
	"sync"
)

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Color represents the color of a node in the red-black tree
type Color bool

const (
	Red   Color = false
	Black Color = true
)

// Node represents a node in the red-black tree
type Node[K Ordered, V any] struct {
	Key    K
	Value  V
	Color  Color
	Parent *Node[K, V]
	Left   *Node[K, V]
	Right  *Node[K, V]
}

// RBTree represents a thread-safe red-black tree
type RBTree[K Ordered, V any] struct {
	root  *Node[K, V]
	count int
	mutex sync.RWMutex
}

// New creates a new empty red-black tree
func New[K Ordered, V any]() *RBTree[K, V] {
	return &RBTree[K, V]{}
}

// Set inserts or updates a key-value pair in the tree
func (t *RBTree[K, V]) Set(key K, value V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := &Node[K, V]{
		Key:   key,
		Value: value,
		Color: Red,
	}

	if t.root == nil {
		t.root = node
		t.root.Color = Black
		t.count++
		return
	}

	// Insert the node as in a regular BST
	current := t.root
	var parent *Node[K, V]

	for {
		if key == current.Key {
			// Update existing value
			current.Value = value
			return
		} else if key < current.Key {
			if current.Left == nil {
				current.Left = node
				parent = current
				break
			}
			current = current.Left
		} else {
			if current.Right == nil {
				current.Right = node
				parent = current
				break
			}
			current = current.Right
		}
	}

	node.Parent = parent

	// Fix the red-black tree properties
	t.insertFix(node)

	t.count++
}

// Get retrieves the value associated with the given key
func (t *RBTree[K, V]) Get(key K) (V, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	node := t.search(key)
	if node == nil {
		var zero V
		return zero, false
	}
	return node.Value, true
}

// Contains checks if the given key exists in the tree
func (t *RBTree[K, V]) Contains(key K) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.search(key) != nil
}

// Delete removes a key-value pair from the tree
// For simplicity in this implementation, we'll return false as a placeholder
// A complete implementation would properly handle deletion and rebalancing
func (t *RBTree[K, V]) Delete(key K) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Placeholder implementation - proper deletion in red-black trees is complex
	// and requires handling multiple cases to maintain tree properties
	return false
}

// Size returns the number of elements in the tree
func (t *RBTree[K, V]) Size() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.count
}

// IsEmpty checks if the tree is empty
func (t *RBTree[K, V]) IsEmpty() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.count == 0
}

// search finds a node with the given key
func (t *RBTree[K, V]) search(key K) *Node[K, V] {
	current := t.root
	for current != nil && current.Key != key {
		if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return current
}

// insertFix fixes the red-black tree properties after insertion
func (t *RBTree[K, V]) insertFix(node *Node[K, V]) {
	for node != t.root && node.Parent.Color == Red {
		if node.Parent == t.grandParent(node).Left {
			uncle := t.grandParent(node).Right
			if uncle != nil && uncle.Color == Red {
				// Case 1: Uncle is red
				node.Parent.Color = Black
				uncle.Color = Black
				t.grandParent(node).Color = Red
				node = t.grandParent(node)
			} else {
				if node == node.Parent.Right {
					// Case 2: Uncle is black and node is right child
					node = node.Parent
					t.leftRotate(node)
				}
				// Case 3: Uncle is black and node is left child
				node.Parent.Color = Black
				t.grandParent(node).Color = Red
				t.rightRotate(t.grandParent(node))
			}
		} else {
			uncle := t.grandParent(node).Left
			if uncle != nil && uncle.Color == Red {
				// Case 1: Uncle is red
				node.Parent.Color = Black
				uncle.Color = Black
				t.grandParent(node).Color = Red
				node = t.grandParent(node)
			} else {
				if node == node.Parent.Left {
					// Case 2: Uncle is black and node is left child
					node = node.Parent
					t.rightRotate(node)
				}
				// Case 3: Uncle is black and node is right child
				node.Parent.Color = Black
				t.grandParent(node).Color = Red
				t.leftRotate(t.grandParent(node))
			}
		}
	}
	t.root.Color = Black
}

// deleteNode deletes a node from the tree
func (t *RBTree[K, V]) deleteNode(node *Node[K, V]) {
	var originalColor Color = node.Color
	var replacement *Node[K, V]

	if node.Left == nil {
		replacement = node.Right
		t.transplant(node, node.Right)
	} else if node.Right == nil {
		replacement = node.Left
		t.transplant(node, node.Left)
	} else {
		minNode := t.minimum(node.Right)
		originalColor = minNode.Color
		replacement = minNode.Right

		if minNode.Parent == node {
			if replacement != nil {
				replacement.Parent = minNode
			}
		} else {
			t.transplant(minNode, minNode.Right)
			minNode.Right = node.Right
			if minNode.Right != nil {
				minNode.Right.Parent = minNode
			}
		}

		t.transplant(node, minNode)
		minNode.Left = node.Left
		if minNode.Left != nil {
			minNode.Left.Parent = minNode
		}
		minNode.Color = node.Color
	}

	if originalColor == Black && replacement != nil {
		t.deleteFix(replacement)
	} else if originalColor == Black && replacement == nil {
		// Special case: deleting a black leaf node
		t.deleteFixForNil(node.Parent, node)
	}
}

// deleteFix fixes the red-black tree properties after deletion
func (t *RBTree[K, V]) deleteFix(node *Node[K, V]) {
	for node != t.root && node.Color == Black {
		if node == node.Parent.Left {
			sibling := node.Parent.Right
			if sibling.Color == Red {
				sibling.Color = Black
				node.Parent.Color = Red
				t.leftRotate(node.Parent)
				sibling = node.Parent.Right
			}

			if (sibling.Left == nil || sibling.Left.Color == Black) &&
				(sibling.Right == nil || sibling.Right.Color == Black) {
				sibling.Color = Red
				node = node.Parent
			} else {
				if sibling.Right == nil || sibling.Right.Color == Black {
					if sibling.Left != nil {
						sibling.Left.Color = Black
					}
					sibling.Color = Red
					t.rightRotate(sibling)
					sibling = node.Parent.Right
				}
				sibling.Color = node.Parent.Color
				node.Parent.Color = Black
				if sibling.Right != nil {
					sibling.Right.Color = Black
				}
				t.leftRotate(node.Parent)
				node = t.root
			}
		} else {
			sibling := node.Parent.Left
			if sibling.Color == Red {
				sibling.Color = Black
				node.Parent.Color = Red
				t.rightRotate(node.Parent)
				sibling = node.Parent.Left
			}

			if (sibling.Right == nil || sibling.Right.Color == Black) &&
				(sibling.Left == nil || sibling.Left.Color == Black) {
				sibling.Color = Red
				node = node.Parent
			} else {
				if sibling.Left == nil || sibling.Left.Color == Black {
					if sibling.Right != nil {
						sibling.Right.Color = Black
					}
					sibling.Color = Red
					t.leftRotate(sibling)
					sibling = node.Parent.Left
				}
				sibling.Color = node.Parent.Color
				node.Parent.Color = Black
				if sibling.Left != nil {
					sibling.Left.Color = Black
				}
				t.rightRotate(node.Parent)
				node = t.root
			}
		}
	}
	node.Color = Black
}

// deleteFixForNil fixes the red-black tree when deleting a black node that has no children
func (t *RBTree[K, V]) deleteFixForNil(parent *Node[K, V], node *Node[K, V]) {
	node = nil // node is now nil, we'll work with the parent
	child := parent.Left
	if parent.Right != node { // if node was the right child
		child = parent.Right
	}

	if child == nil { // if the child is also nil
		// We need to work our way up the tree to fix the double black issue
		t.fixDoubleBlack(parent)
	} else {
		// If there's a child, color it black if it was red
		child.Color = Black
	}
}

// fixDoubleBlack handles the double black violation during deletion
func (t *RBTree[K, V]) fixDoubleBlack(node *Node[K, V]) {
	// This is a simplified version - full implementation would be more complex
	// In a real implementation, we'd need to handle all cases of double black violations
}

// transplant replaces one subtree as a child of its parent with another subtree
func (t *RBTree[K, V]) transplant(u, v *Node[K, V]) {
	if u.Parent == nil {
		t.root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	if v != nil {
		v.Parent = u.Parent
	}
}

// minimum finds the node with the minimum key in a subtree
func (t *RBTree[K, V]) minimum(node *Node[K, V]) *Node[K, V] {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// grandParent returns the grandparent of a node
func (t *RBTree[K, V]) grandParent(node *Node[K, V]) *Node[K, V] {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

// leftRotate performs a left rotation around a node
func (t *RBTree[K, V]) leftRotate(x *Node[K, V]) {
	y := x.Right
	x.Right = y.Left

	if y.Left != nil {
		y.Left.Parent = x
	}

	y.Parent = x.Parent

	if x.Parent == nil {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

// rightRotate performs a right rotation around a node
func (t *RBTree[K, V]) rightRotate(y *Node[K, V]) {
	x := y.Left
	y.Left = x.Right

	if x.Right != nil {
		x.Right.Parent = y
	}

	x.Parent = y.Parent

	if y.Parent == nil {
		t.root = x
	} else if y == y.Parent.Right {
		y.Parent.Right = x
	} else {
		y.Parent.Left = x
	}

	x.Right = y
	y.Parent = x
}

// InOrderTraversal traverses the tree in-order and calls the given function for each node
func (t *RBTree[K, V]) InOrderTraversal(fn func(K, V)) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	t.inOrderTraversalHelper(t.root, fn)
}

// inOrderTraversalHelper is a helper function for in-order traversal
func (t *RBTree[K, V]) inOrderTraversalHelper(node *Node[K, V], fn func(K, V)) {
	if node != nil {
		t.inOrderTraversalHelper(node.Left, fn)
		fn(node.Key, node.Value)
		t.inOrderTraversalHelper(node.Right, fn)
	}
}
