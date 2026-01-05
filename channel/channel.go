// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package channel

import (
	"context"
	"sync"
)

// Channel представляет универсальный канал, который можно использовать для отправки и получения значений типа T.
type Channel[T any] struct {
	ch       chan T
	mu       sync.RWMutex
	closed   bool
	isBuffer bool
}

// New создает новый канал. Если буфер больше 0, создается буферизованный канал.
func New[T any](buffer ...int) *Channel[T] {
	if len(buffer) > 1 {
		panic("channel.New: too many arguments")
	}
	if len(buffer) == 0 {
		return &Channel[T]{
			ch:       make(chan T),
			isBuffer: false,
		}
	}
	return &Channel[T]{
		ch:       make(chan T, buffer[0]),
		isBuffer: buffer[0] > 0,
	}
}

// Send отправляет значение в канал. Может быть отменена контекстом.
func (c *Channel[T]) Send(ctx context.Context, value T) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return ErrClosedChannel
	}
	c.mu.RUnlock()

	select {
	case c.ch <- value:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Receive получает значение из канала. Может быть отменена контекстом.
func (c *Channel[T]) Receive(ctx context.Context) (T, error) {
	select {
	case val, ok := <-c.ch:
		if !ok {
			return val, ErrClosedChannel
		}
		return val, nil
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Close закрывает канал.
func (c *Channel[T]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		close(c.ch)
		c.closed = true
	}
}

// Len возвращает количество элементов, находящихся в канале в данный момент.
func (c *Channel[T]) Len() int {
	return len(c.ch)
}

// Cap возвращает емкость канала.
func (c *Channel[T]) Cap() int {
	return cap(c.ch)
}

// IsClosed возвращает true, если канал закрыт.
func (c *Channel[T]) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// IsBuffered возвращает true, если канал буферизован.
func (c *Channel[T]) IsBuffered() bool {
	return c.isBuffer
}

// Range предоставляет способ итерации по каналу до тех пор, пока он не будет закрыт.
func (c *Channel[T]) Range(f func(value T) bool) {
	for val := range c.ch {
		if !f(val) {
			break
		}
	}
}

// Select позволяет выполнить операцию выбора на канале.
func (c *Channel[T]) Select(ctx context.Context, cases ...*Case[T]) (int, T, error) {
	// Это упрощенная реализация. Полная реализация потребовала бы
	// более сложной логики для обработки нескольких случаев.
	if len(cases) == 0 {
		<-ctx.Done()
		var zero T
		return -1, zero, ctx.Err()
	}

	// For simplicity, we'll just handle the first case.
	// A real implementation would use reflection or a different approach.
	if cases[0].Send {
		err := c.Send(ctx, cases[0].Value)
		if err != nil {
			var zero T
			return -1, zero, err
		}
		return 0, cases[0].Value, nil
	}

	val, err := c.Receive(ctx)
	if err != nil {
		var zero T
		return -1, zero, err
	}
	return 0, val, nil
}

// Case представляет собой случай в операторе выбора.
type Case[T any] struct {
	Chan  *Channel[T]
	Send  bool // true for send, false for receive
	Value T    // value to send
}

var (
	// ErrClosedChannel возвращается, когда выполняется попытка операции над закрытым каналом.
	ErrClosedChannel = &errorString{"channel is closed"}
)

// errorString - тривиальная реализация ошибки.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Unwrap возвращает базовый канал Go. Это в основном для тестирования и расширенных случаев использования.
func (c *Channel[T]) Unwrap() chan T {
	return c.ch
}
