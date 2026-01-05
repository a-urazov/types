# Канал

Пакет `channel` предоставляет обобщенную реализацию канала с расширенной функциональностью.

## Особенности

- **Обобщенный (Generic)**: Может работать с любым типом данных.
- **Буферизованный и небуферизованный**: Поддерживает как синхронные (небуферизованные), так и асинхронные (буферизованные) каналы.
- **Контекстно-зависимые операции**: Операции отправки и получения могут быть отменены с помощью `context.Context`.
- **Безопасность при закрытии**: Предотвращает панику при отправке в закрытый канал.
- **Итерация с помощью `Range`**: Удобный способ итерации по значениям в канале.
- **`Select`**: Базовая поддержка операции `select`.

## Установка

```bash
go get github.com/your-username/types/channel
```

## Использование

### Создание канала

Вы можете создать как буферизованный, так и небуферизованный канал.

```go
import "github.com/your-username/types/channel"

// Небуферизованный канал
ch1 := channel.New[int](0)

// Буферизованный канал
ch2 := channel.New[string](10)
```

### Отправка и получение

```go
import (
	"context"
	"fmt"
)

ch := channel.New[int](0)

go func() {
	err := ch.Send(context.Background(), 42)
	if err != nil {
		// Обработка ошибки
	}
}()

val, err := ch.Receive(context.Background())
if err != nil {
	// Обработка ошибки
}
fmt.Println(val) // Вывод: 42
```

### Закрытие канала

```go
ch.Close()

// Попытка отправки в закрытый канал вернет ошибку
err := ch.Send(context.Background(), 1)
fmt.Println(err) // Вывод: channel is closed
```

### Итерация

```go
ch := channel.New[int](3)
ch.Send(context.Background(), 1)
ch.Send(context.Background(), 2)
ch.Send(context.Background(), 3)
ch.Close()

ch.Range(func(value int) bool {
	fmt.Println(value)
	return true // Продолжить итерацию
})
```
