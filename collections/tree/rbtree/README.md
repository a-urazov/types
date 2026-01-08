# Red-Black Tree

A thread-safe Red-Black Tree implementation in Go with generic types.

## Overview

A Red-Black Tree is a self-balancing binary search tree where each node has an extra bit for denoting the color of the node, either red or black. The balancing of the tree is not perfect, but it guarantees searching, insertion, and deletion operations in O(log n) time.

## Properties

*   Every node is either red or black
*   The root is always black
*   No two adjacent red nodes (red node cannot have a red parent or red child)
*   Every path from root to null/leaf node has the same number of black nodes
*   All leaf nodes are black

## Usage

```
package main

import (
    "fmt"
    "types/collections/tree/rbtree"
)

func main() {
    // Create a new Red-Black Tree
    tree := rbtree.New[int, string]()

    // Insert key-value pairs
    tree.Set(10, "ten")
    tree.Set(5, "five")
    tree.Set(15, "fifteen")

    // Get values
    if val, ok := tree.Get(10); ok {
        fmt.Printf("Value for key 10: %s\n", val)
    }

    // Check if a key exists
    if tree.Contains(5) {
        fmt.Println("Key 5 exists in the tree")
    }

    // Get the size of the tree
    fmt.Printf("Tree size: %d\n", tree.Size())

    // Perform in-order traversal
    tree.InOrderTraversal(func(key int, value string) {
        fmt.Printf("%d: %s\n", key, value)
    })
}
```

## Methods

*   `New[K, V]() *RBTree[K, V]` - Creates a new empty Red-Black Tree
*   `Set(key K, value V)` - Inserts or updates a key-value pair
*   `Get(key K) (V, bool)` - Retrieves the value for a given key
*   `Contains(key K) bool` - Checks if a key exists in the tree
*   `Size() int` - Returns the number of elements in the tree
*   `IsEmpty() bool` - Checks if the tree is empty
*   `Delete(key K) bool` - Removes a key-value pair from the tree (placeholder implementation)
*   `InOrderTraversal(fn func(K, V))` - Traverses the tree in-order and calls the function for each node

## Thread Safety

This implementation is thread-safe using read-write mutexes for concurrent access. Multiple readers can access the tree simultaneously, but write operations are exclusive.
