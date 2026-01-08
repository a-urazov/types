// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package channel

import (
	"context"
	"reflect"
	"sync"
	"time"
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

// Len возвращает количество элементов, находящихихся в канале в данный момент.
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
	if len(cases) == 0 {
		<-ctx.Done()
		var zero T
		return -1, zero, ctx.Err()
	}

	// Convert our cases to reflect.SelectCase
	reflectCases := make([]reflect.SelectCase, len(cases)+1) // +1 for ctx.Done()

	// Add context done case
	reflectCases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	}

	// Add each case
	for i, case_ := range cases {
		if case_.Send {
			reflectCases[i+1] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(case_.Chan.ch),
				Send: reflect.ValueOf(case_.Value),
			}
		} else {
			reflectCases[i+1] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(case_.Chan.ch),
			}
		}
	}

	chosen, value, recvOK := reflect.Select(reflectCases)

	// If context was chosen (index 0), return context error
	if chosen == 0 {
		var zero T
		return -1, zero, ctx.Err()
	}

	// Adjust index since we had the ctx.Done() case at index 0
	caseIndex := chosen - 1
	selectedCase := cases[caseIndex]

	if selectedCase.Send {
		// For send operations, we return the value that was sent
		return caseIndex, selectedCase.Value, nil
	}
	// For receive operations, we need to convert the received value
	if !recvOK {
		// Channel was closed
		var zero T
		return caseIndex, zero, ErrClosedChannel
	}

	// Convert the received value back to type T
	receivedValue, ok := value.Interface().(T)
	if !ok {
		var zero T
		return caseIndex, zero, &errorString{"type assertion failed"}
	}

	return caseIndex, receivedValue, nil
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

// TrySend пытается отправить значение в канал без блокировки.
// Возвращает true, если значение было успешно отправлено, иначе false.
func (c *Channel[T]) TrySend(value T) bool {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return false
	}
	c.mu.RUnlock()

	select {
	case c.ch <- value:
		return true
	default:
		return false
	}
}

// TryReceive пытается получить значение из канала без блокировки.
// Возвращает полученное значение и true, если значение было доступно, иначе zero value и false.
func (c *Channel[T]) TryReceive() (T, bool) {
	select {
	case val, ok := <-c.ch:
		if !ok {
			var zero T
			return zero, false
		}
		return val, true
	default:
		var zero T
		return zero, false
	}
}

// Drain извлекает и отбрасывает все доступные значения из канала.
// Возвращает количество извлеченных значений.
func (c *Channel[T]) Drain() int {
	count := 0
	for {
		select {
		case _, ok := <-c.ch:
			if !ok {
				return count
			}
			count++
		default:
			return count
		}
	}
}

// Done возвращает канал, который будет закрыт при закрытии этого канала.
// Полезно для сигнализации о завершении.
func (c *Channel[T]) Done() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range c.ch { // читаем до закрытия канала
		}
	}()
	return done
}

// Merge объединяет несколько каналов в один, передавая значения из каждого входного канала.
// Результат - новый канал, из которого будут поступать значения из всех входных каналов.
func Merge[T any](channels ...*Channel[T]) *Channel[T] {
	out := New[T]()

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(inputCh *Channel[T]) {
			defer wg.Done()
			for {
				val, err := inputCh.Receive(context.Background())
				if err != nil {
					return // канал закрыт или произошла ошибка
				}

				if err := out.Send(context.Background(), val); err != nil {
					return // выходной канал закрыт
				}
			}
		}(ch)
	}

	// Закрываем выходной канал после завершения всех горутин
	go func() {
		wg.Wait()
		out.Close()
	}()

	return out
}

// FanOut создает несколько выходных каналов и рассылает каждое полученное значение во все выходные каналы.
// Полезно для дублирования потока данных.
func FanOut[T any](input *Channel[T], outputs ...*Channel[T]) {
	go func() {
		for {
			value, err := input.Receive(context.Background())
			if err != nil {
				// Входной канал закрыт, закрываем все выходные
				closeAllOutputs(outputs)
				return
			}

			// Отправляем значение во все выходные каналы
			sendToAllOutputs(value, outputs)
		}
	}()
}

// вспомогательная функция для закрытия всех выходных каналов
func closeAllOutputs[T any](outputs []*Channel[T]) {
	for _, out := range outputs {
		out.Close()
	}
}

// вспомогательная функция для отправки значения во все выходные каналы
func sendToAllOutputs[T any](value T, outputs []*Channel[T]) {
	for _, out := range outputs {
		if err := out.Send(context.Background(), value); err != nil {
			// Если не можем отправить в один из каналов, продолжаем с остальными
			continue
		}
	}
}

// Unwrap возвращает базовый канал Go. Это в основном для тестирования и расширенных случаев использования.
func (c *Channel[T]) Unwrap() chan T {
	return c.ch
}

// BatchSend отправляет несколько значений в канал за одну операцию.
// Возвращает количество успешно отправленных значений.
func (c *Channel[T]) BatchSend(ctx context.Context, values []T) int {
	count := 0
	for _, value := range values {
		err := c.Send(ctx, value)
		if err != nil {
			return count // прерываем при ошибке
		}
		count++
	}
	return count
}

// BatchReceive получает несколько значений из канала за одну операцию.
// Параметр maxCount ограничивает максимальное количество получаемых значений.
func (c *Channel[T]) BatchReceive(ctx context.Context, count int) ([]T, error) {
	results := make([]T, 0, count)
	for range count {
		value, err := c.Receive(ctx)
		if err != nil {
			return results, err
		}
		results = append(results, value)
	}
	return results, nil
}

// WithTimeout оборачивает операцию с таймаутом.
func WithTimeout[T any](timeout time.Duration, operation func() (T, error)) (T, error) {
	resultChan := make(chan struct {
		value T
		err   error
	}, 1)

	go func() {
		defer close(resultChan)
		value, err := operation()
		resultChan <- struct {
			value T
			err   error
		}{value, err}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case result := <-resultChan:
		return result.value, result.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Map применяет функцию к каждому значению из канала и отправляет результат в новый канал.
func Map[T, R any](input *Channel[T], fn func(T) R) *Channel[R] {
	resultChan := New[R]()

	go func() {
		defer resultChan.Close()
		for {
			value, err := input.Receive(context.Background())
			if err != nil {
				return // канал закрыт
			}
			result := fn(value)
			if err := resultChan.Send(context.Background(), result); err != nil {
				return // выходной канал закрыт
			}
		}
	}()

	return resultChan
}

// Filter фильтрует значения из канала, отправляя только те, которые удовлетворяют условию, в новый канал.
func Filter[T any](input *Channel[T], predicate func(T) bool) *Channel[T] {
	output := New[T]()

	go func() {
		defer output.Close()
		for {
			value, err := input.Receive(context.Background())
			if err != nil {
				return // канал закрыт
			}
			if predicate(value) {
				if err := output.Send(context.Background(), value); err != nil {
					return // выходной канал закрыт
				}
			}
		}
	}()

	return output
}

// Take возвращает новый канал с первыми n значениями из исходного канала.
func Take[T any](input *Channel[T], n int) *Channel[T] {
	output := New[T]()

	go func() {
		defer output.Close()
		for range n {
			value, err := input.Receive(context.Background())
			if err != nil {
				return // канал закрыт
			}
			if err := output.Send(context.Background(), value); err != nil {
				return // выходной канал закрыт
			}
		}
	}()

	return output
}

// Skip пропускает первые n значений из канала и возвращает новый канал с оставшимися значениями.
func Skip[T any](input *Channel[T], n int) *Channel[T] {
	output := New[T]()

	go func() {
		defer output.Close()

		// Пропускаем первые n значений
		for range n {
			_, err := input.Receive(context.Background())
			if err != nil {
				return // канал закрыт
			}
		}

		// Передаем оставшиеся значения в выходной канал
		for {
			value, err := input.Receive(context.Background())
			if err != nil {
				return // канал закрыт
			}
			if err := output.Send(context.Background(), value); err != nil {
				return // выходной канал закрыт
			}
		}
	}()

	return output
}

// ConditionalSend sends a value to the channel only if the condition is true.
// Returns error if sending fails, nil otherwise.
func (c *Channel[T]) ConditionalSend(ctx context.Context, value T, condition func(T) bool) error {
	if condition(value) {
		return c.Send(ctx, value)
	}
	return nil // Condition not met, but not an error
}

// ConditionalReceive receives a value from the channel only if it satisfies the condition.
// Continues waiting until a value satisfies the condition or an error occurs.
func (c *Channel[T]) ConditionalReceive(ctx context.Context, condition func(T) bool) (T, error) {
	for {
		value, err := c.Receive(ctx)
		if err != nil {
			var zero T
			return zero, err
		}
		if condition(value) {
			return value, nil
		}
		// Continue waiting if condition is not satisfied
	}
}
