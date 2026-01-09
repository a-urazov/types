# Context - DI контейнер для Go

Пакет `context` предоставляет легкий и типобезопасный DI (Dependency Injection) контейнер для управления зависимостями в Go приложениях.

Кроме того, наш `Context` полностью реализует интерфейс `context.Context` из стандартной библиотеки, что позволяет использовать его как обычный контекст в функциях, которые ожидают `context.Context`.

## Особенности

- **Типобезопасность**: Использует Go generics для полной типобезопасности во время компиляции
- **Поддержка Lifetime**: Singleton, Transient и подготовка к Scoped
- **Thread-safe**: Безопасен для использования в многопоточных приложениях
- **Простой API**: Интуитивный и легкий в использовании интерфейс
- **Flexible**: Поддерживает регистрацию через конструкторы и экземпляры
- **Полная реализация context.Context**: Поддерживает Deadline, Cancel, Value, Done и все методы стандартного контекста

## Использование

### Базовая регистрация Singleton

```go
ctx := context.New()

// Регистрируем Logger как Singleton (один экземпляр на приложение)
context.RegisterSingleton[Logger](ctx, func(c *context.Context) (any, error) {
    return &SimpleLogger{name: "app"}, nil
})

// Разрешаем Logger
logger, err := context.Resolve[Logger](ctx)
if err != nil {
    log.Fatal(err)
}

```

```go

```

### Регистрация Transient

```go
// Регистрируем Logger как Transient (новый экземпляр каждый раз)
context.RegisterTransient[Logger](ctx, func(c *context.Context) (any, error) {
    return &SimpleLogger{name: "app"}, nil
})

// Каждый раз создается новый экземпляр
logger1, _ := context.Resolve[Logger](ctx)
logger2, _ := context.Resolve[Logger](ctx)
// logger1 != logger2

```

### Регистрация экземпляра

```go
logger := &SimpleLogger{name: "singleton"}
context.RegisterInstance[Logger](ctx, logger)

// Всегда возвращает точно этот экземпляр
resolved, _ := context.Resolve[Logger](ctx)
// resolved == logger

```

### Вложенные зависимости

```go
// Регистрируем Logger
context.RegisterSingleton[Logger](ctx, func(c *context.Context) (any, error) {
    return &SimpleLogger{}, nil
})

// Регистрируем Database
context.RegisterSingleton[Database](ctx, func(c *context.Context) (any, error) {
    return &MockDatabase{}, nil
})

// Регистрируем Service с зависимостями
context.RegisterSingleton[*Service](ctx, func(c *context.Context) (any, error) {
    logger, _ := context.Resolve[Logger](c)
    db, _ := context.Resolve[Database](c)

    return &Service{
        logger: logger,
        db:     db,
    }, nil
})

// Разрешаем Service с автоматическим разрешением его зависимостей
service, _ := context.Resolve[*Service](ctx)

```

## API

### `New() *Context`

Создает новый пустой DI контейнер.

```go
ctx := context.New()

```

### `Register[T](ctx *Context, constructor ConstructorFunc, lifetime Lifetime) error`

Регистрирует зависимость с указанным конструктором и временем жизни.

```go
context.Register[Logger](ctx, constructor, context.Singleton)

```

### `RegisterSingleton[T](ctx *Context, constructor ConstructorFunc) error`

Сокращение для `Register` с `Singleton` lifetime.

```go
context.RegisterSingleton[Logger](ctx, func(c *context.Context) (any, error) {
    return &SimpleLogger{}, nil
})

```

### `RegisterTransient[T](ctx *Context, constructor ConstructorFunc) error`

Сокращение для `Register` с `Transient` lifetime.

```go
context.RegisterTransient[Logger](ctx, func(c *context.Context) (any, error) {
    return &SimpleLogger{}, nil
})

```

### `RegisterInstance[T](ctx *Context, instance T) error`

Регистрирует конкретный экземпляр как Singleton.

```go
logger := &SimpleLogger{name: "app"}
context.RegisterInstance[Logger](ctx, logger)

```

### `Resolve[T](ctx *Context) (T, error)`

Получает или создает экземпляр зарегистрированной зависимости.

```go
logger, err := context.Resolve[Logger](ctx)
if err != nil {
    log.Fatal(err)
}

```

### `Contains[T](ctx *Context) bool`

Проверяет, зарегистрирована ли зависимость для указанного типа.

```go
if context.Contains[Logger](ctx) {
    logger, _ := context.Resolve[Logger](ctx)
}

```

### `GetServices(ctx *Context) int`

Возвращает количество зарегистрированных зависимостей (для отладки).

```go
count := context.GetServices(ctx)

```

## Lifetime

- **Singleton**: Один экземпляр на всё время жизни контейнера. Идеален для stateless сервисов.
- **Transient**: Новый экземпляр создается при каждом разрешении. Хорош для stateful объектов.
- **Scoped**: (Планируется) Один экземпляр в пределах scope. Полезно для веб-приложений.

## Thread-safety

Контейнер полностью thread-safe благодаря использованию `sync.RWMutex` и `sync.Mutex`. Можно безопасно использовать в многопоточных приложениях.

## Пример полного приложения

```go
package main

import (
    "log"
    "types/context"
)

type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(msg string) {
    log.Println(msg)
}

type Database interface {
    Query(sql string) ([]string, error)
}

type MockDB struct{}

func (db *MockDB) Query(sql string) ([]string, error) {
    return []string{"result1", "result2"}, nil
}

type UserService struct {
    logger Logger
    db     Database
}

func (s *UserService) GetUsers() {
    s.logger.Log("Получение пользователей...")
    results, _ := s.db.Query("SELECT * FROM users")
    s.logger.Log(string(len(results)))
}

func main() {
    ctx := context.New()

    // Регистрируем зависимости
    context.RegisterSingleton[Logger](ctx, func(c *context.Context) (any, error) {
        return &ConsoleLogger{}, nil
    })

    context.RegisterSingleton[Database](ctx, func(c *context.Context) (any, error) {
        return &MockDB{}, nil
    })

    context.RegisterSingleton[*UserService](ctx, func(c *context.Context) (any, error) {
        logger, _ := context.Resolve[Logger](c)
        db, _ := context.Resolve[Database](c)
        return &UserService{logger: logger, db: db}, nil
    })

    // Разрешаем и используем сервис
    service, _ := context.Resolve[*UserService](ctx)
    service.GetUsers()
}

```

## Производительность

Контейнер оптимизирован для производительности с кэшированием Singleton экземпляров и минимальными блокировками для операций разрешения.

## context.Context методы

Наш `Context` реализует полный интерфейс `context.Context` из стандартной библиотеки:

### Базовые методы

#### `Deadline() (time.Time, bool)`

Возвращает deadline для контекста.

```go
ctx := context.New()
deadline, ok := ctx.Deadline()
if ok {
    fmt.Printf("Deadline: %v\n", deadline)
}

```

#### `Done() <-chan struct{}`

Возвращает канал, который закрывается при отмене контекста.

```go
ctx := context.New()
go func() {
    time.Sleep(time.Second)
    ctx.Cancel()
}()

<-ctx.Done()
fmt.Println("Context canceled")

```

#### `Err() error`

Возвращает ошибку отмены (если контекст отменен).

```go
ctx := context.New()
ctx.Cancel()
if err := ctx.Err(); err != nil {
    fmt.Printf("Error: %v\n", err)
}

```

#### `Value(key any) any`

Получает значение, связанное с ключом.

```go
ctx := context.New()
ctx.SetValue("user", "alice")
user := ctx.Value("user")

```

### Методы отмены

#### `Cancel()`

Отменяет контекст с ошибкой `context.Canceled`.

```go
ctx := context.New()
ctx.Cancel()

```

#### `CancelWithError(err error)`

Отменяет контекст с указанной ошибкой.

```go
ctx := context.New()
ctx.CancelWithError(fmt.Errorf("custom error"))

```

### Методы создания дочерних контекстов

#### `WithCancel() (*Context, context.CancelFunc)`

Создает новый контекст с возможностью отмены.

```go
ctx := context.New()
newCtx, cancel := ctx.WithCancel()
defer cancel()

// Использование newCtx...

```

#### `WithDeadline(deadline time.Time) (*Context, context.CancelFunc)`

Создает новый контекст с deadline.

```go
ctx := context.New()
deadline := time.Now().Add(5 * time.Second)
newCtx, cancel := ctx.WithDeadline(deadline)
defer cancel()

```

#### `WithTimeout(timeout time.Duration) (*Context, context.CancelFunc)`

Создает новый контекст с timeout.

```go
ctx := context.New()
newCtx, cancel := ctx.WithTimeout(5 * time.Second)
defer cancel()

```

#### `WithValue(key any, value any) *Context`

Создает новый контекст с добавленным значением.

```go
ctx := context.New()
newCtx := ctx.WithValue("user", "alice")
user := newCtx.Value("user")

```

### Вспомогательные методы

#### `SetValue(key any, value any)`

Устанавливает значение в текущем контексте.

```go
ctx := context.New()
ctx.SetValue("request_id", "12345")

```

#### `NewWithContext(parent context.Context) *Context`

Создает новый DI контейнер с parent контекстом из стандартной библиотеки.

```go
parentCtx := context.Background()
ctx := context.NewWithContext(parentCtx)

```

## Примеры

### Пример 1: DI контейнер с контекстом

```go
package main

import (
    "fmt"
    "types/context"
)

type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(msg string) {
    fmt.Println(msg)
}

func main() {
    // Создаем контекст
    ctx := context.New()

    // Устанавливаем значение контекста
    ctx.SetValue("appName", "MyApp")

    // Регистрируем Logger
    context.RegisterSingleton[Logger](ctx, func(c *context.Context) (any, error) {
        appName := c.Value("appName").(string)
        fmt.Printf("Creating logger for %s\n", appName)
        return &ConsoleLogger{}, nil
    })

    // Разрешаем Logger
    logger, _ := context.Resolve[Logger](ctx)
    logger.Log("Hello, World!")
}

```

### Пример 2: Работа с Deadline

```go
ctx := context.New()
newCtx, cancel := ctx.WithTimeout(2 * time.Second)
defer cancel()

// Использование newCtx...
if err := newCtx.Err(); err == context.DeadlineExceeded {
    fmt.Println("Timeout exceeded")
}

```

### Пример 3: Наследование значений

```go
ctx := context.New()
ctx.SetValue("request_id", "123")

// Дочерний контекст наследует значение
childCtx := ctx.WithValue("user_id", "456")

fmt.Println(childCtx.Value("request_id")) // "123"
fmt.Println(childCtx.Value("user_id"))    // "456"

```

### Пример 4: Использование с обычным context.Context API

```go
// Наш Context реализует context.Context, поэтому его можно использовать везде
func handleRequest(ctx context.Context, userID string) {
    select {
    case <-ctx.Done():
        fmt.Println("Request canceled")
    case <-time.After(5 * time.Second):
        fmt.Println("Request completed")
    }
}

func main() {
    ctx := context.New()
    handleRequest(ctx, "user123")
}

```