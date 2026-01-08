# Пакет Result

Пакет `result` предоставляет тип `Result[T]` для управления ошибками с использованием шаблона Result/Either. Этот тип может содержать либо успешное значение типа `T`, либо ошибку.

## Особенности

- Обобщенная реализация Result типа
- Методы для безопасной обработки ошибок
- Функции для монадических операций (map, andThen, orElse и т.д.)
- Паттерн-матчинг с помощью метода Match

## Установка

```bash
# Пакет является частью модуля types
```

## Использование

### Создание Result

```go
// Создание успешного результата
success := result.Ok(42)

// Создание результата с ошибкой
err := result.Err[int](errors.New("что-то пошло не так"))
```

### Проверка статуса

```go
if success.IsOk() {
    fmt.Println("Операция успешна")
}

if err.IsErr() {
    fmt.Println("Произошла ошибка")
}
```

### Извлечение значений

```go
// Безопасное извлечение значения
value, err := success.Unwrap()
if err == nil {
    fmt.Printf("Значение: %d\n", value)
}

// Извлечение со значением по умолчанию
value = success.UnwrapOr(0)
```

### Цепочки операций

```go
// Цепочка операций с обработкой ошибок
result := someOperation().
    AndThen(func(x int) result.Result[int] {
        return result.Ok(x * 2)
    }).
    Map(func(x int) int {
        return x + 1
    })
```

## API

- `Ok(value T) Result[T]` - Создает Result с успешным значением
- `Err(err error) Result[T]` - Создает Result с ошибкой
- `IsOk() bool` - Проверяет, содержит ли Result успешное значение
- `IsErr() bool` - Проверяет, содержит ли Result ошибку
- `Unwrap() (T, error)` - Возвращает значение и ошибку
- `UnwrapOr(defaultValue T) T` - Возвращает значение или значение по умолчанию
- `UnwrapOrElse(fn func(error) T) T` - Возвращает значение или результат функции
- `Expect(msg string) T` - Возвращает значение или паникует с сообщением
- `UnwrapErr() error` - Возвращает ошибку
- `And(other Result[T]) Result[T]` - Возвращает другой Result если текущий Ok
- `AndThen(fn func(T) Result[T]) Result[T]` - Выполняет функцию если текущий Ok
- `Or(other Result[T]) Result[T]` - Возвращает другой Result если текущий Err
- `OrElse(fn func(error) Result[T]) Result[T]` - Возвращает результат функции если текущий Err
- `Map(fn func(T) T) Result[T]` - Применяет функцию к значению
- `MapOr(defaultValue T, fn func(T) T) T` - Применяет функцию или возвращает defaultValue
- `MapOrElse(fallbackFn func(error) T, fn func(T) T) T` - Применяет функцию или фоллбэк
- `Match(okFn func(T), errFn func(error))` - Паттерн-матчинг
- `FromError(value T, err error) Result[T]` - Преобразует значение и ошибку в Result
- `FromFunc(fn func() (T, error)) Result[T]` - Выполняет функцию и оборачивает результат в Result
