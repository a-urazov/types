# Словарь

Пакет `dictionary` предоставляет реализацию словаря (ассоциативного массива или отображения). Словарь хранит пары ключ-значение и позволяет быстро найти значение по его ключу.

## Основные методы

- `New[TKey, TValue]() *Dictionary[TKey, TValue]` - создает новый пустой словарь
- `Set(key TKey, value TValue)` - добавляет или обновляет пару ключ-значение
- `Get(key TKey) (TValue, bool)` - получает значение по ключу
- `Remove(key TKey) bool` - удаляет пару по ключу
- `Contains(key TKey) bool` - проверяет наличие ключа
- `Size() int` - возвращает количество пар
- `IsEmpty() bool` - проверяет, пуст ли словарь
- `Clear()` - очищает словарь
- `Keys() []TKey` - возвращает все ключи
- `Values() []TValue` - возвращает все значения
- `ForEach(fn func(key TKey, value TValue))` - итерирует по всем парам

## Пример использования

```go
package main

import (
    "fmt"
    "types/collections/dictionary"
)

func main() {
    dict := dictionary.New[string, int]()

    dict.Set("one", 1)
    dict.Set("two", 2)
    dict.Set("three", 3)

    if val, ok := dict.Get("two"); ok {
        fmt.Println("Value for 'two':", val)
    }

    dict.ForEach(func(key string, val int) {
        fmt.Printf("%s: %d\n", key, val)
    })

    fmt.Println("Size:", dict.Size())
}
```

## Специализированные варианты

- **[Sorted Dictionary](./sorted/)** - словарь, который хранит ключи в отсортированном порядке

## Особенности

- **Быстрый поиск**: O(1) в среднем случае благодаря хеширванию
- **Потокобезопасность**: встроенная синхронизация
- **Гибкость**: поддерживает любые сравнимые типы ключей и любые типы значений

## Применение

- Кэширование данных
- Конфигурационные параметры
- Индексирование данных
- Счетчики частоты элементов
- Преобразование между разными идентификаторами
