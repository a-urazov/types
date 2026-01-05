# List

Пакет `list` предоставляет универсальную реализацию списка с поддержкой добавления, вставки и удаления элементов. Список потокобезопасен благодаря встроенному мутексу.

## Основные методы

- `New[T]() *List[T]` - создает новый пустой список
- `Add(item T)` - добавляет элемент в конец списка
- `Insert(index int, item T) bool` - вставляет элемент по указанному индексу
- `Remove(index int) (T, bool)` - удаляет элемент по индексу
- `Get(index int) (T, bool)` - получает элемент по индексу
- `Size() int` - возвращает количество элементов
- `Clear()` - очищает список
- `Contains(item T) bool` - проверяет наличие элемента
- `ForEach(fn func(item T))` - итерирует по всем элементам

## Пример использования

```go
package main

import (
    "fmt"
    "types/collections/list"
)

func main() {
    l := list.New[string]()
    
    l.Add("first")
    l.Add("second")
    l.Add("third")
    
    if item, ok := l.Get(0); ok {
        fmt.Println("First item:", item)
    }
    
    l.Insert(1, "inserted")
    fmt.Println("Size:", l.Size())
}
```

## Специализированные варианты

- **[LinkedList](./linked/)** - двусвязный список
- **[SkipList](./skip/)** - список с перепрыгиванием
