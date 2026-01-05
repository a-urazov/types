package trie

// TrieNode представляет узел в trie.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// NewTrieNode создает и возвращает новый узел trie.
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

// Tree представляет структуру данных trie (префиксное дерево).
type Tree struct {
	root *TrieNode
}

// New создает и возвращает новый пустой Trie.
func New() *Tree {
	return &Tree{
		root: NewTrieNode(),
	}
}

// Insert добавляет слово в trie.
func (t *Tree) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Search возвращает true, если слово существует в trie, иначе false.
func (t *Tree) Search(word string) bool {
	node := t.root
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

// StartsWith возвращает true, если в trie существует слово, начинающееся с данного префикса, иначе false.
func (t *Tree) StartsWith(prefix string) bool {
	node := t.root
	for _, char := range prefix {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}
	return true
}

// deleteTrieNode - вспомогательная функция для рекурсивного удаления слова из trie.
// Возвращает (deleted bool, shouldDeleteNode bool).
// `deleted` равно true, если слово было найдено и помечено как удаленное.
// `shouldDeleteNode` равно true, если текущий узел может быть физически удален из trie.
func deleteTrieNode(node *TrieNode, runes []rune, index int) (bool, bool) {
	if node == nil {
		return false, false
	}

	if index == len(runes) {
		// Достигнут конец слова
		if !node.isEnd {
			// Слово не существует в trie
			return false, false
		}
		node.isEnd = false // Отметить как не конец слова
		deleted := true

		// Если этот узел не имеет дочерних узлов, он может быть удален
		if len(node.children) == 0 {
			return deleted, true
		}
		return deleted, false // Узел все еще является частью других слов
	}

	char := runes[index]
	child, exists := node.children[char]
	if !exists {
		// Слово не существует в trie
		return false, false
	}

	deleted, shouldDeleteChild := deleteTrieNode(child, runes, index+1)

	if shouldDeleteChild {
		delete(node.children, char)
		// Если этот узел не имеет других дочерних узлов и не является концом другого слова, он может быть удален
		if len(node.children) == 0 && !node.isEnd {
			return deleted, true
		}
	}

	return deleted, false
}

// Delete удаляет слово из trie.
// Возвращает true, если слово было найдено и удалено, иначе false.
func (t *Tree) Delete(word string) bool {
	deleted, _ := deleteTrieNode(t.root, []rune(word), 0)
	return deleted
}

// Size возвращает количество слов в trie.
func (t *Tree) Size() int {
	return size(t.root)
}

// size - вспомогательная функция для рекурсивного подсчета количества слов в trie.
func size(node *TrieNode) int {
	if node == nil {
		return 0
	}

	count := 0
	if node.isEnd {
		count = 1
	}

	for _, child := range node.children {
		count += size(child)
	}

	return count
}

// IsEmpty возвращает true, если trie пуст (не имеет слов), иначе false.
func (t *Tree) IsEmpty() bool {
	return t.Size() == 0
}
