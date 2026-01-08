# Тип Option

Тип `Option` представляет значение, которое может присутствовать или отсутствовать. Это функциональная альтернатива nullable типам, помогающая избежать исключений указателя на ноль, делая отсутствие значения явным.

## Features

- Generic type with one type parameter: `Option[T any]`
- Methods to check if a value is present (`IsSome`, `IsNone`)
- Safe extraction of values with validity checking
- Helper methods like `GetOrElse` and `GetOrCall` for default values
- Functional programming methods like `Map`, `FlatMap`, and `Filter`
- JSON marshaling/unmarshaling support

## Usage

```go
package main

import (
    "fmt"
    "types/option"
)

func main() {
    // Create an Option with a value (Some)
    someValue := option.Some(42)

    if someValue.IsSome() {
        value, ok := someValue.Get()
        if ok {
            fmt.Println("Got value:", value) // Prints: "Got value: 42"
        }
    }

    // Create an Option without a value (None)
    noneValue := option.None[int]()

    if noneValue.IsNone() {
        fmt.Println("No value present")
    }

    // Using GetOrElse for default values
    result := noneValue.GetOrElse(100)
    fmt.Println("Result:", result) // Prints: "Result: 100"

    // Using Map to transform a value if present
    transformed := option.Map(someValue, func(x int) string {
        return fmt.Sprintf("Value is: %d", x)
    })

    if transformed.IsSome() {
        result, _ := transformed.Get()
        fmt.Println(result) // Prints: "Value is: 42"
    }
}
```

## API

- `Some[T](value T) Option[T]` - Creates an Option with a value present
- `None[T]() Option[T]` - Creates an Option with no value present
- `IsSome() bool` - Checks if the Option contains a value
- `IsNone() bool` - Checks if the Option does not contain a value
- `Get() (T, bool)` - Gets the value and validity flag
- `GetOrElse(defaultValue T) T` - Returns the value or a default
- `GetOrCall(f func() T) T` - Returns the value or result of calling a function
- `Map[T, U any](opt Option[T], fn func(T) U) Option[U]` - Maps a function over the value if present
- `FlatMap[T, U any](opt Option[T], fn func(T) Option[U]) Option[U]` - FlatMaps a function that returns an Option
- `Filter[T any](opt Option[T], predicate func(T) bool) Option[T]` - Filters the Option based on a predicate
- `String() string` - Returns a string representation
