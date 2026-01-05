package trie

import (
	"testing"
)

func TestNewTrie(t *testing.T) {
	trie := New()
	if trie == nil {
		t.Error("New() should not return nil")
	}
	if !trie.IsEmpty() {
		t.Error("New trie should be empty")
	}
	if trie.Size() != 0 {
		t.Errorf("New trie size should be 0, got %d", trie.Size())
	}
}

func TestInsert(t *testing.T) {
	trie := New()
	words := []string{"apple", "app", "application", "apply", "banana", "band", "bandana"}

	for _, word := range words {
		trie.Insert(word)
		if !trie.Search(word) {
			t.Errorf("Trie should contain inserted word %s", word)
		}
		if trie.Size() != countInsertedWords(trie, words[:indexOf(words, word)+1]) {
			t.Errorf("Trie size should be %d after inserting %d words, got %d", countInsertedWords(trie, words[:indexOf(words, word)+1]), indexOf(words, word)+1, trie.Size())
		}
	}
}

func TestSearch(t *testing.T) {
	trie := New()
	words := []string{"hello", "world", "golang", "trie"}

	for _, word := range words {
		trie.Insert(word)
	}

	for _, word := range words {
		if !trie.Search(word) {
			t.Errorf("Trie should contain inserted word %s", word)
		}
	}

	nonExistentWords := []string{"help", "word", "go", "tree"}
	for _, word := range nonExistentWords {
		if trie.Search(word) {
			t.Errorf("Trie should not contain non-existent word %s", word)
		}
	}
}

func TestStartsWith(t *testing.T) {
	trie := New()
	words := []string{"cat", "car", "card", "care", "careful", "carefully"}

	for _, word := range words {
		trie.Insert(word)
	}

	prefixes := []string{"c", "ca", "car", "care", "caref", "carefu", "careful"}
	for _, prefix := range prefixes {
		if !trie.StartsWith(prefix) {
			t.Errorf("Trie should have a word starting with %s", prefix)
		}
	}

	nonPrefixes := []string{"dog", "caz", "carex", "carefuly"}
	for _, prefix := range nonPrefixes {
		if trie.StartsWith(prefix) {
			t.Errorf("Trie should not have a word starting with %s", prefix)
		}
	}
}

func TestDelete(t *testing.T) {
	trie := New()
	words := []string{"test", "testing", "tester", "tea", "team", "tear"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Удалить слово, которое является префиксом другого
	if !trie.Delete("test") {
		t.Error("Delete should return true for existing word 'test'")
	}
	if trie.Search("test") {
		t.Error("Trie should not contain deleted word 'test'")
	}
	if !trie.Search("testing") {
		t.Error("Trie should still contain 'testing' after deleting 'test'")
	}
	if trie.Size() != 5 { // testing, tester, tea, team, tear
		t.Errorf("Trie size should be 5 after deleting 'test', got %d", trie.Size())
	}

	// Удалить слово, которое не является префиксом другого
	if !trie.Delete("tea") {
		t.Error("Delete should return true for existing word 'tea'")
	}
	if trie.Search("tea") {
		t.Error("Trie should not contain deleted word 'tea'")
	}
	if trie.Size() != 4 { // testing, tester, team, tear
		t.Errorf("Trie size should be 4 after deleting 'tea', got %d", trie.Size())
	}

	// Удалить слово, которое имеет общий префикс с другими
	if !trie.Delete("team") {
		t.Error("Delete should return true for existing word 'team'")
	}
	if trie.Search("team") {
		t.Error("Trie should not contain deleted word 'team'")
	}
	if !trie.Search("tear") {
		t.Error("Trie should still contain 'tear' after deleting 'team'")
	}
	if trie.Size() != 3 { // testing, tester, tear
		t.Errorf("Trie size should be 3 after deleting 'team', got %d", trie.Size())
	}

	// Удалить несуществующее слово
	if trie.Delete("nonexistent") {
		t.Error("Delete should return false for non-existent word")
	}

	// Удалить последнее слово
	remainingWords := []string{"testing", "tester", "tear"}
	for i, word := range remainingWords {
		if !trie.Delete(word) {
			t.Errorf("Delete should return true for existing word %s", word)
		}
		if trie.Search(word) {
			t.Errorf("Trie should not contain deleted word %s", word)
		}
		expectedSize := len(remainingWords) - i - 1
		if trie.Size() != expectedSize {
			t.Errorf("Trie size should be %d after deleting %s, got %d", expectedSize, word, trie.Size())
		}
	}

	if !trie.IsEmpty() {
		t.Error("Trie should be empty after deleting all words")
	}
}

func TestSize(t *testing.T) {
	trie := New()
	words := []string{"a", "aa", "aaa", "b", "bb", "bbb"}

	if trie.Size() != 0 {
		t.Errorf("Empty trie size should be 0, got %d", trie.Size())
	}

	for i, word := range words {
		trie.Insert(word)
		if trie.Size() != i+1 {
			t.Errorf("Trie size should be %d after inserting %d words, got %d", i+1, i+1, trie.Size())
		}
	}

	// Delete some words and check size
	trie.Delete("aa")
	if trie.Size() != 5 {
		t.Errorf("Trie size should be 5 after deleting 'aa', got %d", trie.Size())
	}

	trie.Delete("bb")
	if trie.Size() != 4 {
		t.Errorf("Trie size should be 4 after deleting 'bb', got %d", trie.Size())
	}

	// Удалить все
	for _, word := range words {
		trie.Delete(word)
	}
	if trie.Size() != 0 {
		t.Errorf("Trie size should be 0 after deleting all words, got %d", trie.Size())
	}
}

func TestIsEmpty(t *testing.T) {
	trie := New()
	if !trie.IsEmpty() {
		t.Error("New trie should be empty")
	}

	trie.Insert("word")
	if trie.IsEmpty() {
		t.Error("Trie should not be empty after insertion")
	}

	trie.Delete("word")
	if !trie.IsEmpty() {
		t.Error("Trie should be empty after deleting the last word")
	}
}

// Вспомогательная функция для поиска индекса строки в срезе
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// Вспомогательная функция для подсчета количества уникальных вставленных слов
func countInsertedWords(trie *Tree, words []string) int {
	count := 0
	seen := make(map[string]bool)
	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			count++
		}
	}
	return count
}
