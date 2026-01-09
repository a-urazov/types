# Пакет converter

Пакет `converter` предоставляет набор инструментов для преобразования значений одного типа в другой. Он поддерживает как встроенные типы (строки, числа, булевы значения), так и сложные структуры данных (карты, срезы).

## Основные возможности

- **Универсальное преобразование**: Функция `Convert[From, To]` позволяет преобразовывать значения между любыми типами
- **Специализированные конвертеры**: Предусмотрены оптимизированные конвертеры для типов: string, number, bool, map, slice
- **Интерфейс Converter**: Позволяет создавать собственные конвертеры
- **JSON-based преобразование**: Использует JSON маршалинг/анмаршалинг для сложных типов

## API

### Функция Convert

```cpp
#include <string>

void main() {
    printf("Hello, World!\n");
    return 0;
}

```

```go
func Convert[From, To any](value From) (To, error)

```

Преобразует значение типа `From` в значение типа `To`.

**Пример:**

```go
// Преобразование int в string
result, err := converter.Convert[int, string](42)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // "42"

```

### Интерфейс Converter

```go
type Converter interface {
    Convert(source any) (any, error)
}

```

Интерфейс для создания собственных конвертеров.

### Встроенные конвертеры

#### StringConverter
Преобразует значения в строку.

```go
converter := converter.StringConverter{}
result, _ := converter.Convert(42)
fmt.Println(result) // "42"

```

#### NumberConverter
Преобразует значения в число (float64).

```go
converter := converter.NumberConverter{}
result, _ := converter.Convert("123.45")
fmt.Println(result) // 123.45

```

#### BoolConverter
Преобразует значения в булево значение.

```go
converter := converter.BoolConverter{}
result, _ := converter.Convert("true")
fmt.Println(result) // true

```

#### MapConverter
Преобразует значения в map[string]any.

```go
converter := converter.MapConverter{}
result, _ := converter.Convert(struct{Name string}{Name: "John"})

```

#### SliceConverter
Преобразует значения в []any.

```go
converter := converter.SliceConverter{}
result, _ := converter.Convert([]int{1, 2, 3})

```

### Функция GetConverter

```go
func GetConverter(targetType string) Converter

```

Возвращает подходящий конвертер для целевого типа.

**Примеры типов:**
- `"string"` - StringConverter
- `"number"`, `"float"`, `"float64"`, `"int"` - NumberConverter
- `"bool"`, `"boolean"` - BoolConverter
- `"map"` - MapConverter
- `"slice"`, `"array"` - SliceConverter

**Пример:**

```go
converter := converter.GetConverter("string")
result, _ := converter.Convert(42)
fmt.Println(result) // "42"

```

### Функция ConvertTo

```go
func ConvertTo(source any, targetType reflect.Type) (any, error)

```

Преобразует значение к целевому типу, используя `reflect`.

## Примеры использования

### Базовое преобразование

```go
package main

import (
    "fmt"
    "log"
    "types/converter"
)

func main() {
    // int -> string
    result1, err := converter.Convert[int, string](42)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result1) // "42"

    // float64 -> string
    result2, err := converter.Convert[float64, string](3.14)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result2) // "3.14"
}

```

### Использование специализированных конвертеров

```go
package main

import (
    "fmt"
    "types/converter"
)

func main() {
    // StringConverter
    strConv := converter.StringConverter{}
    result, _ := strConv.Convert(42)
    fmt.Println(result) // "42"

    // NumberConverter
    numConv := converter.NumberConverter{}
    result, _ = numConv.Convert("123.45")
    fmt.Println(result) // 123.45

    // BoolConverter
    boolConv := converter.BoolConverter{}
    result, _ = boolConv.Convert("yes")
    fmt.Println(result) // true
}

```

### Использование GetConverter

```go
package main

import (
    "fmt"
    "types/converter"
)

func main() {
    conv := converter.GetConverter("number")
    result, _ := conv.Convert("123.45")
    fmt.Println(result) // 123.45
}

```

### Создание собственного конвертера

```go
package main

import (
    "fmt"
    "types/converter"
)

// CustomConverter преобразует значение в пользовательский формат
type CustomConverter struct{}

func (cc CustomConverter) Convert(source any) (any, error) {
    return fmt.Sprintf("[CUSTOM] %v", source), nil
}

func main() {
    conv := CustomConverter{}
    result, _ := conv.Convert(42)
    fmt.Println(result) // "[CUSTOM] 42"
}

```

## Поддерживаемые преобразования

### Числовые преобразования
- int/int8/int16/int32/int64 ↔ float32/float64
- uint/uint8/uint16/uint32/uint64 ↔ числовые типы
- string → float64 (парсинг)
- bool ↔ float64 (1.0/0.0)

### Строковые преобразования
- Любое значение → string (через fmt.Sprintf или JSON)
- Поддержка fmt.Stringer интерфейса

### Булевы преобразования
- "true", "1", "yes", "on" → true
- "false", "0", "no", "off", "" → false
- Числовые типы (0 → false, остальные → true)

### Сложные типы
- Struct → map[string]any (через JSON)
- Slice/Array → []any

## Обработка ошибок

Все функции преобразования возвращают `error` в качестве второго значения. Важно всегда проверять ошибки:

```go
result, err := converter.Convert[string, int]("not a number")
if err != nil {
    log.Printf("ошибка преобразования: %v", err)
}

```

## Производительность

- Для встроенных типов преобразование является быстрым и прямым
- Для сложных типов используется JSON маршалинг, что может быть медленнее
- Рекомендуется кэшировать конвертеры при интенсивном использовании

## Лицензия

Часть библиотеки типов Go