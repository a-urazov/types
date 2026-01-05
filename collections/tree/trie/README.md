# Trie

Пакет `trie` предоставляет реализацию префиксного дерева (Trie), также известного как дерево поиска. Trie используется для эффективного хранения и поиска строк, где есть общие префиксы.

## Основные методы

- `New() *Tree` - создает новое пустое Trie
- `Insert(word string)` - вставляет слово в Trie
- `Search(word string) bool` - ищет точное совпадение слова
- `StartsWith(prefix string) bool` - проверяет, есть ли слова с заданным префиксом
- `Delete(word string) bool` - удаляет слово из Trie
- `GetAllWordsWithPrefix(prefix string) []string` - возвращает все слова с заданным префиксом

## Пример использования

```go
package main

import (
    "fmt"
    "types/collections/tree/trie"
)

func main() {
    trie := trie.New()
    
    trie.Insert("apple")
    trie.Insert("app")
    trie.Insert("application")
    trie.Insert("apply")
    trie.Insert("banana")
    
    // Поиск точного совпадения
    if trie.Search("apple") {
        fmt.Println("'apple' found")
    }
    
    // Проверка префикса
    if trie.StartsWith("app") {
        fmt.Println("Words starting with 'app' exist")
    }
    
    // Получить все слова с префиксом
    words := trie.GetAllWordsWithPrefix("app")
    fmt.Println("Words with prefix 'app':", words)
    // Output: [app apple application apply]
}
```

## Особенности

- **Быстрая автодополнение**: поиск всех слов с префиксом O(k), где k - количество результатов
- **Экономия памяти**: общие префиксы хранятся один раз
- **Последовательный поиск**: поиск слова O(m), где m - длина слова
- **Простота реализации**: понятная и прямолинейная структура

## Производительность

| Операция | Временная сложность |
|----------|---|
| Вставка | O(m), где m - длина слова |
| Поиск | O(m) |
| Проверка префикса | O(m) |
| Получение всех с префиксом | O(n), где n - количество узлов в поддереве |

## Применение

- Автодополнение в текстовых редакторах и поисковых системах
- Проверка орфографии и исправление опечаток
- Поиск слов в словаре
- IP маршрутизация
- Системы быстрого поиска
- Реализация фильтров и классификаторов
