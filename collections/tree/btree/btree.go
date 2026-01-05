package btree

import "cmp"

const t = 3 // Минимальная степень (определяет диапазон количества ключей в узле)

// Node представляет узел в B-дереве.
type Node[T cmp.Ordered] struct {
	keys     []T        // Массив ключей
	children []*Node[T] // Массив указателей на дочерние узлы
	leaf     bool       // True, если узел является листом
}

// NewNode создает и возвращает новый узел с заданными свойствами.
func NewNode[T cmp.Ordered](leaf bool) *Node[T] {
	return &Node[T]{
		keys:     make([]T, 0, 2*t-1),      // Емкость для 2t-1 ключей
		children: make([]*Node[T], 0, 2*t), // Емкость для 2t дочерних узлов
		leaf:     leaf,
	}
}

// Tree представляет структуру данных B-дерева.
type Tree[T cmp.Ordered] struct {
	root *Node[T]
}

// New создает и возвращает новое пустое B-дерево.
func New[T cmp.Ordered]() *Tree[T] {
	return &Tree[T]{
		root: NewNode[T](true), // Корень изначально является листом
	}
}

// Insert добавляет значение в B-дерево.
func (bt *Tree[T]) Insert(value T) {
	root := bt.root
	// Если корень полный, дерево увеличивается в высоту
	if len(root.keys) == 2*t-1 {
		newRoot := NewNode[T](false)
		newRoot.children = append(newRoot.children, root)
		bt.splitChild(newRoot, 0)
		bt.insertNonFull(newRoot, value)
		bt.root = newRoot
	} else {
		bt.insertNonFull(root, value)
	}
}

// insertNonFull вставляет значение в узел, который не полный.
func (bt *Tree[T]) insertNonFull(node *Node[T], value T) {
	i := len(node.keys) - 1

	if node.leaf {
		// Если узел является листом, вставить значение напрямую
		node.keys = append(node.keys, *new(T)) // Добавить место для нового ключа
		for i >= 0 && node.keys[i] > value {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = value
	} else {
		// Если узел не листе, найти дочерний узел для вставки
		for i >= 0 && node.keys[i] > value {
			i--
		}
		i++
		child := node.children[i]

		// Если дочерний узел полный, разделить его
		if len(child.keys) == 2*t-1 {
			bt.splitChild(node, i)
			// После разделения средний ключ перемещается вверх, и у нас есть 2 дочерних узла
			// Решить, в какой из двух дочерних узлов вставить
			if node.keys[i] < value {
				i++
			}
			child = node.children[i]
		}

		bt.insertNonFull(child, value)
	}
}

// splitChild разделяет дочерний узел заданного узла.
// i - индекс дочерного узла в node.children, который будет разделен.
func (bt *Tree[T]) splitChild(node *Node[T], i int) {
	z := node.children[i]
	y := z

	// Создать новый узел для хранения второй половины ключей и дочерних узлов y
	x := NewNode[T](y.leaf)
	node.children = append(node.children, nil) // Освободить место для нового дочернего узла
	// Сдвинуть дочерние узлы вправо, чтобы освободить место для x
	copy(node.children[i+2:], node.children[i+1:])
	node.children[i+1] = x

	// Переместить вторую половину ключей y в x
	midIndex := t - 1
	midKey := y.keys[midIndex] // Ключ для перемещения в родительский узел
	x.keys = y.keys[midIndex+1:]
	y.keys = y.keys[:midIndex] // Хранить только первую половину в y

	// Переместить вторую половину дочерних узлов y в x (если y не лист)
	if !y.leaf {
		x.children = y.children[midIndex+1:]
		y.children = y.children[:midIndex+1]
	}

	// Переместить средний ключ y вверх в узел
	node.keys = append(node.keys, *new(T)) // Освободить место для нового ключа
	// Сдвинуть ключи вправо, чтобы освободить место для midKey
	copy(node.keys[i+1:], node.keys[i:])
	node.keys[i] = midKey
}

// Search возвращает true, если значение существует в B-дереве, иначе false.
func (bt *Tree[T]) Search(value T) bool {
	return bt.searchHelper(bt.root, value)
}

// searchHelper - вспомогательная рекурсивная функция для поиска.
func (bt *Tree[T]) searchHelper(node *Node[T], value T) bool {
	i := 0
	for i < len(node.keys) && value > node.keys[i] {
		i++
	}

	if i < len(node.keys) && value == node.keys[i] {
		return true
	}

	if node.leaf {
		return false
	}

	return bt.searchHelper(node.children[i], value)
}

// Delete удаляет значение из B-дерева.
// Возвращает true, если значение было найдено и удалено, иначе false.
func (bt *Tree[T]) Delete(value T) bool {
	deleted := bt.deleteHelper(bt.root, value)
	if len(bt.root.keys) == 0 && !bt.root.leaf {
		// Если корень становится пустым и имеет дочерний узел, сделать его новым корнем
		bt.root = bt.root.children[0]
	}
	return deleted
}

// deleteHelper - вспомогательная рекурсивная функция для удаления.
// Возвращает true, если значение было найдено и удалено.
func (bt *Tree[T]) deleteHelper(node *Node[T], value T) bool {
	i := 0
	for i < len(node.keys) && value > node.keys[i] {
		i++
	}

	// Если ключ найден в этом узле
	if i < len(node.keys) && value == node.keys[i] {
		return bt.deleteKeyFromNode(node, i)
	}

	// Если ключ не найден в этом узле
	if node.leaf {
		// Ключ не существует в дереве
		return false
	}

	// Ключ находится в поддереве
	// Перед тем как спуститься, убедитесь, что дочерний узел, в который мы идем, имеет как минимум t ключей
	// (чтобы мы могли безопасно удалить из него, если нужно)
	child := node.children[i]
	if len(child.keys) < t {
		bt.fill(node, i)
		// После заполнения дочерний узел мог измениться, поэтому нам нужно переоценить
		// Если ключ, который мы ищем, был перемещен вверх во время заполнения, нам нужно это обработать
		// Но для простоты мы просто снова вызовем deleteHelper с тем же значением
		// и позволим ему найти ключ в новой позиции.
		// Это упрощение; более эффективный подход прямо обрабатывал бы новую структуру.
		// Однако для корректности этот рекурсивный вызов приемлем.
		return bt.deleteHelper(node, value)
	}

	return bt.deleteHelper(child, value)
}

// deleteKeyFromNode удаляет ключ с индексом i из узла.
func (bt *Tree[T]) deleteKeyFromNode(node *Node[T], i int) bool {
	if node.leaf {
		// Случай 1: Ключ находится в листовом узле
		// Просто удалить ключ
		copy(node.keys[i:], node.keys[i+1:])
		node.keys[len(node.keys)-1] = *new(T) // Очистить последний элемент
		node.keys = node.keys[:len(node.keys)-1]
		return true
	}

	// Случай 2: Ключ находится во внутреннем узле
	leftChild := node.children[i]
	rightChild := node.children[i+1]

	if len(leftChild.keys) >= t {
		// Случай 2a: Дочерний узел-предшественник имеет как минимум t ключей
		// Найти предшественника (самый правый ключ в левом поддереве)
		pred := bt.getRightmostKey(leftChild)
		node.keys[i] = pred
		return bt.deleteHelper(leftChild, pred)
	} else if len(rightChild.keys) >= t {
		// Случай 2b: Дочерний узел-преемник имеет как минимум t ключей
		// Найти преемника (самый левый ключ в правом поддереве)
		succ := bt.getLeftmostKey(rightChild)
		node.keys[i] = succ
		return bt.deleteHelper(rightChild, succ)
	} else {
		// Случай 2c: Оба дочерних узла имеют t-1 ключей
		// Объединить node[i] и node[i+1] и ключ узла на позиции i
		bt.merge(node, i)
		// После объединения ключ, который был на позиции node.keys[i], теперь находится в объединенном дочернем узле
		// Поэтому мы рекурсивно удаляем исходный ключ из объединенного дочернего узла
		return bt.deleteHelper(node.children[i], node.keys[i])
	}
}

// getRightmostKey возвращает самый правый (наибольший) ключ в поддереве, корнем которого является узел.
func (bt *Tree[T]) getRightmostKey(node *Node[T]) T {
	current := node
	for !current.leaf {
		current = current.children[len(current.children)-1]
	}
	return current.keys[len(current.keys)-1]
}

// getLeftmostKey возвращает самый левый (наименьший) ключ в поддереве, корнем которого является узел.
func (bt *Tree[T]) getLeftmostKey(node *Node[T]) T {
	current := node
	for !current.leaf {
		current = current.children[0]
	}
	return current.keys[0]
}

// fill гарантирует, что дочерний узел на позиции i имеет как минимум t ключей.
func (bt *Tree[T]) fill(node *Node[T], i int) {
	if i != 0 && len(node.children[i-1].keys) >= t {
		// Если левый сосед имеет более t-1 ключей, одолжить от него
		bt.borrowFromPrev(node, i)
	} else if i != len(node.children)-1 && len(node.children[i+1].keys) >= t {
		// Если правый сосед имеет более t-1 ключей, одолжить от него
		bt.borrowFromNext(node, i)
	} else {
		// Иначе объединить с соседом
		if i != len(node.children)-1 {
			bt.merge(node, i)
		} else {
			bt.merge(node, i-1)
		}
	}
}

// borrowFromPrev перемещает ключ из node.children[i-1] в node.children[i].
func (bt *Tree[T]) borrowFromPrev(node *Node[T], i int) {
	child := node.children[i]
	sibling := node.children[i-1]

	// Переместить ключ из родителя в дочерний узел
	child.keys = append([]T{node.keys[i-1]}, child.keys...)
	node.keys[i-1] = sibling.keys[len(sibling.keys)-1]

	// Переместить дочерний узел из соседа в дочерний узел (если не лист)
	if !sibling.leaf {
		child.children = append([]*Node[T]{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	// Удалить заимствованный ключ из соседа
	sibling.keys = sibling.keys[:len(sibling.keys)-1]
}

// borrowFromNext перемещает ключ из node.children[i+1] в node.children[i].
func (bt *Tree[T]) borrowFromNext(node *Node[T], i int) {
	child := node.children[i]
	sibling := node.children[i+1]

	// Переместить ключ из родителя в дочерний узел
	child.keys = append(child.keys, node.keys[i])

	// Переместить ключ из соседа в родителя
	node.keys[i] = sibling.keys[0]

	// Переместить дочерний узел из соседа в дочерний узел (если не лист)
	if !sibling.leaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}

	// Удалить перемещенный ключ из соседа
	sibling.keys = sibling.keys[1:]
}

// merge объединяет node.children[i] с node.children[i+1].
// node.children[i+1] освобождается после объединения.
func (bt *Tree[T]) merge(node *Node[T], i int) {
	child := node.children[i]
	sibling := node.children[i+1]

	// Переместить ключ из узла в дочерний узел
	child.keys = append(child.keys, node.keys[i])

	// Переместить все ключи соседа в дочерний узел
	child.keys = append(child.keys, sibling.keys...)

	// Переместить все дочерние узлы соседа в дочерний узел (если не лист)
	if !sibling.leaf {
		child.children = append(child.children, sibling.children...)
	}

	// Удалить ключ из узла и отрегулировать дочерние узлы
	copy(node.keys[i:], node.keys[i+1:])
	node.keys[len(node.keys)-1] = *new(T) // Очистить последний элемент
	node.keys = node.keys[:len(node.keys)-1]

	copy(node.children[i+1:], node.children[i+2:])
	node.children[len(node.children)-1] = nil // Очистить последний элемент
	node.children = node.children[:len(node.children)-1]
}

// Size возвращает количество элементов в B-дереве.
// Это операция O(n), так как нам нужно пройти по всем узлам.
func (bt *Tree[T]) Size() int {
	return bt.sizeHelper(bt.root)
}

// sizeHelper - вспомогательная рекурсивная функция для подсчета размера.
func (bt *Tree[T]) sizeHelper(node *Node[T]) int {
	if node == nil {
		return 0
	}

	size := len(node.keys)
	if node.leaf {
		return size
	}

	for _, child := range node.children {
		size += bt.sizeHelper(child)
	}
	return size
}

// IsEmpty возвращает true, если B-дерево пусто, иначе false.
func (bt *Tree[T]) IsEmpty() bool {
	return bt.root == nil || len(bt.root.keys) == 0
}
