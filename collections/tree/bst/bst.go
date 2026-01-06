package bst

import "cmp"

// Node представляет узел в двоичном дереве поиска.
type Node[T cmp.Ordered] struct {
	Value T
	Left  *Node[T]
	Right *Node[T]
}

// Tree представляет двоичное дерево поиска.
type Tree[T cmp.Ordered] struct {
	root *Node[T]
}

// New создает и возвращает новое пустое BST.
func New[T cmp.Ordered]() *Tree[T] {
	return &Tree[T]{}
}

// Insert добавляет значение в BST.
func (bst *Tree[T]) Insert(value T) {
	bst.root = insert(bst.root, value)
}

// insert - вспомогательная функция для рекурсивной вставки значения в дерево.
func insert[T cmp.Ordered](node *Node[T], value T) *Node[T] {
	if node == nil {
		return &Node[T]{Value: value}
	}

	if value < node.Value {
		node.Left = insert(node.Left, value)
	} else if value > node.Value {
		node.Right = insert(node.Right, value)
	}
	// Если value == node.Value, мы не вставляем дубликаты (стандартное поведение BST)

	return node
}

// Search возвращает true, если значение существует в BST, иначе false.
func (bst *Tree[T]) Search(value T) bool {
	return search(bst.root, value)
}

// search - вспомогательная функция для рекурсивного поиска значения в дереве.
func search[T cmp.Ordered](node *Node[T], value T) bool {
	if node == nil {
		return false
	}

	if value == node.Value {
		return true
	} else if value < node.Value {
		return search(node.Left, value)
	} else {
		return search(node.Right, value)
	}
}

// Delete удаляет значение из BST.
// Возвращает true, если значение было найдено и удалено, иначе false.
func (bst *Tree[T]) Delete(value T) bool {
	var deleted bool
	bst.root, deleted = delete(bst.root, value)
	return deleted
}

// delete - вспомогательная функция для рекурсивного удаления значения из дерева.
func delete[T cmp.Ordered](node *Node[T], value T) (*Node[T], bool) {
	if node == nil {
		return nil, false
	}

	if value < node.Value {
		newLeft, deleted := delete(node.Left, value)
		node.Left = newLeft
		return node, deleted
	} else if value > node.Value {
		newRight, deleted := delete(node.Right, value)
		node.Right = newRight
		return node, deleted
	} else {
		// Узел для удаления найден
		// Случай 1: Узел не имеет дочерних узлов (листовой узел)
		if node.Left == nil && node.Right == nil {
			return nil, true
		}
		// Случай 2: Узел имеет один дочерний узел
		if node.Left == nil {
			return node.Right, true
		}
		if node.Right == nil {
			return node.Left, true
		}
		// Случай 3: Узел имеет двух дочерних узлов
		// Найти упорядоченного по порядку преемника (наименьшее значение в правом поддереве)
		successor := findMin(node.Right)
		// Заменить значение узла значением преемника
		node.Value = successor.Value
		// Удалить преемника
		newRight, _ := delete(node.Right, successor.Value)
		node.Right = newRight
		return node, true
	}
}

// findMin находит узел с минимальным значением в поддереве.
func findMin[T cmp.Ordered](node *Node[T]) *Node[T] {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// InOrderTraversal выполняет упорядоченный обход дерева и применяет заданную функцию к значению каждого узла.
func (bst *Tree[T]) InOrderTraversal(fn func(T)) {
	inOrderTraversal(bst.root, fn)
}

// inOrderTraversal - вспомогательная функция для рекурсивного выполнения упорядоченного обхода.
func inOrderTraversal[T cmp.Ordered](node *Node[T], fn func(T)) {
	if node != nil {
		inOrderTraversal(node.Left, fn)
		fn(node.Value)
		inOrderTraversal(node.Right, fn)
	}
}

// PreOrderTraversal выполняет предварительный обход дерева и применяет заданную функцию к значению каждого узла.
func (bst *Tree[T]) PreOrderTraversal(fn func(T)) {
	preOrderTraversal(bst.root, fn)
}

// preOrderTraversal - вспомогательная функция для рекурсивного выполнения предварительного обхода.
func preOrderTraversal[T cmp.Ordered](node *Node[T], fn func(T)) {
	if node != nil {
		fn(node.Value)
		preOrderTraversal(node.Left, fn)
		preOrderTraversal(node.Right, fn)
	}
}

// PostOrderTraversal выполняет послепорядковый обход дерева и применяет заданную функцию к значению каждого узла.
func (bst *Tree[T]) PostOrderTraversal(fn func(T)) {
	postOrderTraversal(bst.root, fn)
}

// postOrderTraversal - вспомогательная функция для рекурсивного выполнения послепорядкового обхода.
func postOrderTraversal[T cmp.Ordered](node *Node[T], fn func(T)) {
	if node != nil {
		postOrderTraversal(node.Left, fn)
		postOrderTraversal(node.Right, fn)
		fn(node.Value)
	}
}

// Min возвращает минимальное значение в BST и булево значение, указывающее, не пусто ли дерево.
func (bst *Tree[T]) Min() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}
	return findMin(bst.root).Value, true
}

// Max возвращает максимальное значение в BST и булево значение, указывающее, не пусто ли дерево.
func (bst *Tree[T]) Max() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}
	node := bst.root
	for node.Right != nil {
		node = node.Right
	}
	return node.Value, true
}

// Size возвращает количество узлов в BST.
func (bst *Tree[T]) Size() int {
	return size(bst.root)
}

// size - вспомогательная функция для рекурсивного подсчета количества узлов в дереве.
func size[T cmp.Ordered](node *Node[T]) int {
	if node == nil {
		return 0
	}
	return 1 + size(node.Left) + size(node.Right)
}

// IsEmpty возвращает true, если BST пусто, иначе false.
func (bst *Tree[T]) IsEmpty() bool {
	return bst.root == nil
}
