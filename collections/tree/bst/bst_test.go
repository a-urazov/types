package bst

import (
	"reflect"
	"testing"
)

func TestNewBST(t *testing.T) {
	bst := New[int]()
	if bst == nil {
		t.Error("New() should not return nil")
	}
	if !bst.IsEmpty() {
		t.Error("New BST should be empty")
	}
	if bst.Size() != 0 {
		t.Errorf("New BST size should be 0, got %d", bst.Size())
	}
}

func TestInsert(t *testing.T) {
	bst := New[int]()

	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		bst.Insert(v)
	}

	if bst.Size() != len(values) {
		t.Errorf("BST size should be %d after inserting %d values, got %d", len(values), len(values), bst.Size())
	}

	for _, v := range values {
		if !bst.Search(v) {
			t.Errorf("BST should contain inserted value %d", v)
		}
	}
}

func TestSearch(t *testing.T) {
	bst := New[string]()
	values := []string{"apple", "banana", "cherry", "date"}

	for _, v := range values {
		bst.Insert(v)
	}

	for _, v := range values {
		if !bst.Search(v) {
			t.Errorf("BST should contain inserted value %s", v)
		}
	}

	nonExistentValues := []string{"grape", "kiwi", "orange"}
	for _, v := range nonExistentValues {
		if bst.Search(v) {
			t.Errorf("BST should not contain non-existent value %s", v)
		}
	}
}

func TestDelete(t *testing.T) {
	bst := New[int]()
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		bst.Insert(v)
	}

	// Удалить листовой узел
	if !bst.Delete(2) {
		t.Error("Delete should return true for existing value")
	}
	if bst.Search(2) {
		t.Error("BST should not contain deleted value 2")
	}
	if bst.Size() != len(values)-1 {
		t.Errorf("BST size should be %d after deleting one value, got %d", len(values)-1, bst.Size())
	}

	// Удалить узел с одним дочерним узлом
	bst.Insert(1)       // Добавить 1 как левого потомка узла 2 (который был удален, поэтому 1 будет левым потомком узла 3)
	if !bst.Delete(2) { // Попытаться удалить 2 снова (должно возвращать false)
		t.Log("Попытка удалить уже удаленное значение 2, это должно возвращать false")
	}
	bst.Insert(2) // Повторно добавить 2
	if !bst.Delete(1) {
		t.Error("Delete should return true for existing value 1")
	}
	if bst.Search(1) {
		t.Error("BST should not contain deleted value 1")
	}

	// Удалить узел с двумя дочерними узлами (3)
	// В этот момент: [5, 3, 7, 4, 6, 8, 2] (размер 7 после предыдущих операций)
	// Удалить 3: 3 заменяется его упорядоченным преемником (4), а старый 4 удаляется.
	// Дерево становится: [5, 4, 7, 6, 8, 2] (размер 6)
	if !bst.Delete(3) {
		t.Error("Delete should return true for existing value 3")
	}
	if bst.Search(3) {
		t.Error("BST should not contain deleted value 3")
	}
	if bst.Size() != 6 { // Should be 6 after deleting 3
		t.Errorf("BST size should be 6 after deleting 3 (which had 2 children), got %d", bst.Size())
	}

	// Удалить несуществующее значение
	if bst.Delete(99) {
		t.Error("Delete should return false for non-existent value")
	}
}

func TestTraversals(t *testing.T) {
	bst := New[int]()
	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		bst.Insert(v)
	}

	// Тест InOrderTraversal (должна быть отсортирована)
	inOrderResult := []int{}
	bst.InOrderTraversal(func(v int) {
		inOrderResult = append(inOrderResult, v)
	})
	expectedInOrder := []int{2, 3, 4, 5, 6, 7, 8}
	if !reflect.DeepEqual(inOrderResult, expectedInOrder) {
		t.Errorf("InOrderTraversal result %v does not match expected %v", inOrderResult, expectedInOrder)
	}

	// Тест PreOrderTraversal (корень, левый, правый)
	preOrderResult := []int{}
	bst.PreOrderTraversal(func(v int) {
		preOrderResult = append(preOrderResult, v)
	})
	expectedPreOrder := []int{5, 3, 2, 4, 7, 6, 8} // Это ожидаемый предварительный обход для построенного дерева
	if !reflect.DeepEqual(preOrderResult, expectedPreOrder) {
		t.Errorf("PreOrderTraversal result %v does not match expected %v", preOrderResult, expectedPreOrder)
	}

	// Тест PostOrderTraversal (левый, правый, корень)
	postOrderResult := []int{}
	bst.PostOrderTraversal(func(v int) {
		postOrderResult = append(postOrderResult, v)
	})
	expectedPostOrder := []int{2, 4, 3, 6, 8, 7, 5} // Это ожидаемый послепорядковый обход для построенного дерева
	if !reflect.DeepEqual(postOrderResult, expectedPostOrder) {
		t.Errorf("PostOrderTraversal result %v does not match expected %v", postOrderResult, expectedPostOrder)
	}
}

func TestMin(t *testing.T) {
	bst := New[int]()
	_, ok := bst.Min()
	if ok {
		t.Error("Min should return false for empty tree")
	}

	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		bst.Insert(v)
	}

	min, ok := bst.Min()
	if !ok {
		t.Error("Min should return true for non-empty tree")
	}
	if min != 2 {
		t.Errorf("Min should return 2, got %d", min)
	}
}

func TestMax(t *testing.T) {
	bst := New[int]()
	_, ok := bst.Max()
	if ok {
		t.Error("Max should return false for empty tree")
	}

	values := []int{5, 3, 7, 2, 4, 6, 8}
	for _, v := range values {
		bst.Insert(v)
	}

	max, ok := bst.Max()
	if !ok {
		t.Error("Max should return true for non-empty tree")
	}
	if max != 8 {
		t.Errorf("Max should return 8, got %d", max)
	}
}

func TestSize(t *testing.T) {
	bst := New[int]()
	if bst.Size() != 0 {
		t.Errorf("Empty BST size should be 0, got %d", bst.Size())
	}

	values := []int{5, 3, 7, 2, 4, 6, 8}
	for i, v := range values {
		bst.Insert(v)
		if bst.Size() != i+1 {
			t.Errorf("BST size should be %d after inserting %d values, got %d", i+1, i+1, bst.Size())
		}
	}

	for i := len(values); i > 0; i-- {
		bst.Delete(values[len(values)-i])
		if bst.Size() != i-1 {
			t.Errorf("BST size should be %d after deleting one value, got %d", i-1, bst.Size())
		}
	}
}

func TestIsEmpty(t *testing.T) {
	bst := New[int]()
	if !bst.IsEmpty() {
		t.Error("New BST should be empty")
	}

	bst.Insert(1)
	if bst.IsEmpty() {
		t.Error("BST should not be empty after insertion")
	}

	bst.Delete(1)
	if !bst.IsEmpty() {
		t.Error("BST should be empty after deleting the last value")
	}
}
