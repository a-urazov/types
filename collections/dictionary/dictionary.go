package dictionary

import "sync"

// Dictionary представляет собой универсальный словарь/карту.
type Dictionary[TKey comparable, TValue any] struct {
	items map[TKey]TValue
	mu    sync.RWMutex
}

// New создает новый словарь.
func New[TKey comparable, TValue any]() *Dictionary[TKey, TValue] {
	return &Dictionary[TKey, TValue]{
		items: make(map[TKey]TValue),
	}
}

// Set добавляет или обновляет пару ключ-значение.
func (d *Dictionary[TKey, TValue]) Set(key TKey, value TValue) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items[key] = value
}

// Get извлекает значение по его ключу.
func (d *Dictionary[TKey, TValue]) Get(key TKey) (TValue, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.items[key]
	return val, ok
}

// Remove удаляет пару ключ-значение. Возвращает true, если ключ существовал.
func (d *Dictionary[TKey, TValue]) Remove(key TKey) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.items[key]
	if ok {
		delete(d.items, key)
	}
	return ok
}

// ContainsKey проверяет, существует ли ключ в словаре.
func (d *Dictionary[TKey, TValue]) ContainsKey(key TKey) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.items[key]
	return ok
}

// Keys возвращает срез всех ключей.
func (d *Dictionary[TKey, TValue]) Keys() []TKey {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keys := make([]TKey, 0, len(d.items))
	for k := range d.items {
		keys = append(keys, k)
	}
	return keys
}

// Values возвращает срез всех значений.
func (d *Dictionary[TKey, TValue]) Values() []TValue {
	d.mu.RLock()
	defer d.mu.RUnlock()
	values := make([]TValue, 0, len(d.items))
	for _, v := range d.items {
		values = append(values, v)
	}
	return values
}

// Size возвращает количество элементов в словаре.
func (d *Dictionary[TKey, TValue]) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.items)
}

// IsEmpty возвращает true, если словарь пуст.
func (d *Dictionary[TKey, TValue]) IsEmpty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.items) == 0
}

// Clear очищает словарь.
func (d *Dictionary[TKey, TValue]) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = make(map[TKey]TValue)
}
