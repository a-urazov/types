package dictionary

import (
	"types/internal/common"
)

// Dictionary представляет собой универсальный словарь/карту.
type Dictionary[TKey comparable, TValue any] struct {
	items *common.Map[TKey, TValue]
}

// New создает новый словарь.
func New[TKey comparable, TValue any]() *Dictionary[TKey, TValue] {
	return &Dictionary[TKey, TValue]{
		items: common.NewMap[TKey, TValue](),
	}
}

// Set добавляет или обновляет пару ключ-значение.
func (d *Dictionary[TKey, TValue]) Set(key TKey, value TValue) {
	d.items.WithWriteLock(func(items map[TKey]TValue) map[TKey]TValue {
		items[key] = value
		return items
	})
}

// Get извлекает значение по его ключу.
func (d *Dictionary[TKey, TValue]) Get(key TKey) (TValue, bool) {
	var val TValue
	var ok bool
	d.items.WithReadLock(func(items map[TKey]TValue) {
		val, ok = items[key]
	})
	return val, ok
}

// Remove удаляет пару ключ-значение. Возвращает true, если ключ существовал.
func (d *Dictionary[TKey, TValue]) Remove(key TKey) bool {
	var existed bool
	d.items.WithWriteLock(func(items map[TKey]TValue) map[TKey]TValue {
		if _, ok := items[key]; ok {
			existed = true
			delete(items, key)
		}
		return items
	})
	return existed
}

// ContainsKey проверяет, существует ли ключ в словаре.
func (d *Dictionary[TKey, TValue]) ContainsKey(key TKey) bool {
	var ok bool
	d.items.WithReadLock(func(items map[TKey]TValue) {
		_, ok = items[key]
	})
	return ok
}

// Keys возвращает срез всех ключей.
func (d *Dictionary[TKey, TValue]) Keys() []TKey {
	var keys []TKey
	d.items.WithReadLock(func(items map[TKey]TValue) {
		keys = make([]TKey, 0, len(items))
		for k := range items {
			keys = append(keys, k)
		}
	})
	return keys
}

// Values возвращает срез всех значений.
func (d *Dictionary[TKey, TValue]) Values() []TValue {
	var values []TValue
	d.items.WithReadLock(func(items map[TKey]TValue) {
		values = make([]TValue, 0, len(items))
		for _, v := range items {
			values = append(values, v)
		}
	})
	return values
}

// Size возвращает количество элементов в словаре.
func (d *Dictionary[TKey, TValue]) Size() int {
	var size int
	d.items.WithReadLock(func(items map[TKey]TValue) {
		size = len(items)
	})
	return size
}

// IsEmpty возвращает true, если словарь пуст.
func (d *Dictionary[TKey, TValue]) IsEmpty() bool {
	var empty bool
	d.items.WithReadLock(func(items map[TKey]TValue) {
		empty = len(items) == 0
	})
	return empty
}

// Clear очищает словарь.
func (d *Dictionary[TKey, TValue]) Clear() {
	d.items.SetItems(make(map[TKey]TValue))
}
