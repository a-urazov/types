package sorted

import (
	"cmp"
	"sort"
	"sync"
)

// Dictionary представляет собой словарь, отсортированный по ключам.
type Dictionary[TKey cmp.Ordered, TValue any] struct {
	keys   []TKey
	values []TValue
	mu     sync.RWMutex
}

// New создает новый SortedDictionary.
func New[TKey cmp.Ordered, TValue any]() *Dictionary[TKey, TValue] {
	return &Dictionary[TKey, TValue]{
		keys:   make([]TKey, 0),
		values: make([]TValue, 0),
	}
}

// findIndex выполняет двоичный поиск для нахождения индекса ключа.
// Он возвращает индекс и логическое значение, указывающее, был ли найден ключ.
func (d *Dictionary[TKey, TValue]) findIndex(key TKey) (int, bool) {
	// sort.Search использует двоичный поиск
	i := sort.Search(len(d.keys), func(i int) bool {
		return d.keys[i] >= key
	})
	if i < len(d.keys) && d.keys[i] == key {
		return i, true
	}
	return i, false
}

// Set добавляет или обновляет пару ключ-значение, сохраняя порядок сортировки.
func (d *Dictionary[TKey, TValue]) Set(key TKey, value TValue) {
	d.mu.Lock()
	defer d.mu.Unlock()

	i, found := d.findIndex(key)
	if found {
		d.values[i] = value
		return
	}

	// Вставить новый ключ и значение в правильной позиции
	d.keys = append(d.keys, *new(TKey))
	copy(d.keys[i+1:], d.keys[i:])
	d.keys[i] = key

	d.values = append(d.values, *new(TValue))
	copy(d.values[i+1:], d.values[i:])
	d.values[i] = value
}

// Get извлекает значение по его ключу.
func (d *Dictionary[TKey, TValue]) Get(key TKey) (TValue, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if i, found := d.findIndex(key); found {
		return d.values[i], true
	}
	var zero TValue
	return zero, false
}

// Remove удаляет пару ключ-значение.
func (d *Dictionary[TKey, TValue]) Remove(key TKey) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	i, found := d.findIndex(key)
	if !found {
		return false
	}

	d.keys = append(d.keys[:i], d.keys[i+1:]...)
	d.values = append(d.values[:i], d.values[i+1:]...)
	return true
}

// ContainsKey проверяет, существует ли ключ.
func (d *Dictionary[TKey, TValue]) ContainsKey(key TKey) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, found := d.findIndex(key)
	return found
}

// Keys возвращает срез всех ключей в отсортированном порядке.
func (d *Dictionary[TKey, TValue]) Keys() []TKey {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keysCopy := make([]TKey, len(d.keys))
	copy(keysCopy, d.keys)
	return keysCopy
}

// Values возвращает срез всех значений в порядке ключей.
func (d *Dictionary[TKey, TValue]) Values() []TValue {
	d.mu.RLock()
	defer d.mu.RUnlock()
	valuesCopy := make([]TValue, len(d.values))
	copy(valuesCopy, d.values)
	return valuesCopy
}

// Size возвращает количество элементов в словаре.
func (d *Dictionary[TKey, TValue]) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.keys)
}

// IsEmpty возвращает true, если словарь пуст.
func (d *Dictionary[TKey, TValue]) IsEmpty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.keys) == 0
}

// Clear очищает словарь.
func (d *Dictionary[TKey, TValue]) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.keys = make([]TKey, 0)
	d.values = make([]TValue, 0)
}
