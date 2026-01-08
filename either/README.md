# Тип Either

Тип `Either` представляет значение одного из двух возможных типов - `Left` или `Right`. Это часто используется в функциональном программировании для представления вычислений, которые могут привести либо к значению успеха (`Right`), либо к значению ошибки/сбоя (`Left`).

## Features

- Generic type with two type parameters: `Either[L any, R any]`
- Methods to check if the value is `Left` or `Right`
- Safe extraction of values with validity checking
- Helper methods like `OrLeft` and `OrRight` for default values
- Functional programming methods like `Map` and `FlatMap`
- JSON marshaling/unmarshaling support

## Usage

```go
package main

import (
    "fmt"
    "types/either"
)

func main() {
    // Create a Left value (typically used for errors)
    leftValue := either.Left[string, int]("Something went wrong")

    if leftValue.IsLeft() {
        fmt.Println("Got an error:", leftValue.OrLeft("default"))
    }

    // Create a Right value (typically used for success)
    rightValue := either.Right[string, int](42)

    if rightValue.IsRight() {
        fmt.Println("Got a success value:", rightValue.OrRight(0))
    }

    // Using Map to transform a Right value
    transformed := rightValue.Map(func(x int) string {
        return fmt.Sprintf("Value is: %d", x)
    })

    if transformed.IsRight() {
        result, _ := transformed.GetRight()
        fmt.Println(result) // Prints: "Value is: 42"
    }
}
```

## API

- `Left[L, R](value L) Either[L, R]` - Creates an Either with a Left value
- `Right[L, R](value R) Either[L, R]` - Creates an Either with a Right value
- `IsLeft() bool` - Checks if the Either contains a Left value
- `IsRight() bool` - Checks if the Either contains a Right value
- `GetLeft() (L, bool)` - Gets the Left value and validity flag
- `GetRight() (R, bool)` - Gets the Right value and validity flag
- `OrLeft(defaultValue L) L` - Returns the Left value or a default
- `OrRight(defaultValue R) R` - Returns the Right value or a default
- `String() string` - Returns a string representation
- `MapRight[L, R, T any](e Either[L, R], fn func(R) T) Either[L, T]` - Maps a function over the Right value
- `FlatMapRight[L, R, T any](e Either[L, R], fn func(R) Either[L, T]) Either[L, T]` - FlatMaps a function that returns an Either over the Right value
- `MapLeft[L, R, T any](e Either[L, R], fn func(L) T) Either[T, R]` - Maps a function over the Left value
- `FlatMapLeft[L, R, L2, R2 any](e Either[L, R], fn func(L) Either[L2, R2]) Either[L2, R2]` - FlatMaps a function that returns an Either over the Left value
