package context

import (
	"maps"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// ConstructorFunc представляет функцию для создания экземпляра зависимости
type ConstructorFunc func(ctx *Context) (any, error)

// ServiceDescriptor описывает зарегистрированную зависимость
type ServiceDescriptor struct {
	typ         reflect.Type
	constructor ConstructorFunc
	lifetime    Lifetime
	instance    any
	mu          sync.RWMutex
}

// Lifetime определяет время жизни зависимости
type Lifetime int

const (
	// Transient - новый экземпляр создается каждый раз
	Transient Lifetime = iota
	// Singleton - один экземпляр на всё время жизни приложения
	Singleton
	// Scoped - один экземпляр в пределах scope (для будущих версий)
	Scoped
)

// Context является DI контейнером для управления зависимостями
// и реализует interface context.Context из стандартной библиотеки
type Context struct {
	services map[reflect.Type]*ServiceDescriptor
	mu       sync.RWMutex

	// Поля для реализации context.Context
	parent   context.Context
	deadline time.Time
	done     chan struct{}
	err      error
	values   map[any]any
	mu2      sync.RWMutex
}

// New создает новый пустой DI контейнер
func New() *Context {
	return &Context{
		services: make(map[reflect.Type]*ServiceDescriptor),
		parent:   context.Background(),
		done:     make(chan struct{}),
		values:   make(map[any]any),
	}
}

// NewWithContext создает новый DI контейнер с указанным parent context
func NewWithContext(parent context.Context) *Context {
	return &Context{
		services: make(map[reflect.Type]*ServiceDescriptor),
		parent:   parent,
		done:     make(chan struct{}),
		values:   make(map[any]any),
	}
}

// getType получает reflect.Type для типа T
func getType[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

// Register регистрирует зависимость с использованием конструктора
// T - тип интерфейса или структуры
// constructor - функция создания экземпляра
// lifetime - время жизни зависимости (Singleton, Transient и т.д.)
func Register[T any](ctx *Context, constructor ConstructorFunc, lifetime Lifetime) error {
	typ := getType[T]()

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if _, exists := ctx.services[typ]; exists {
		return fmt.Errorf("зависимость типа %v уже зарегистрирована", typ)
	}

	ctx.services[typ] = &ServiceDescriptor{
		typ:         typ,
		constructor: constructor,
		lifetime:    lifetime,
	}

	return nil
}

// RegisterSingleton регистрирует зависимость как Singleton (один экземпляр на приложение)
func RegisterSingleton[T any](ctx *Context, constructor ConstructorFunc) error {
	return Register[T](ctx, constructor, Singleton)
}

// RegisterTransient регистрирует зависимость как Transient (новый экземпляр каждый раз)
func RegisterTransient[T any](ctx *Context, constructor ConstructorFunc) error {
	return Register[T](ctx, constructor, Transient)
}

// RegisterInstance регистрирует конкретный экземпляр как Singleton
func RegisterInstance[T any](ctx *Context, instance T) error {
	typ := getType[T]()

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if _, exists := ctx.services[typ]; exists {
		return fmt.Errorf("зависимость типа %v уже зарегистрирована", typ)
	}

	ctx.services[typ] = &ServiceDescriptor{
		typ:      typ,
		lifetime: Singleton,
		instance: instance,
	}

	return nil
}

// Resolve получает зарегистрированную зависимость по типу
// Если зависимость не найдена, возвращает ошибку
func Resolve[T any](ctx *Context) (T, error) {
	var t T
	typ := getType[T]()

	ctx.mu.RLock()
	descriptor, exists := ctx.services[typ]
	ctx.mu.RUnlock()

	if !exists {
		return t, fmt.Errorf("зависимость типа %v не найдена", typ)
	}

	// Если это singleton и уже создан, возвращаем cached экземпляр
	if descriptor.lifetime == Singleton {
		descriptor.mu.RLock()
		if descriptor.instance != nil {
			defer descriptor.mu.RUnlock()
			return descriptor.instance.(T), nil
		}
		descriptor.mu.RUnlock()
	}

	// Создаем новый экземпляр
	instance, err := descriptor.constructor(ctx)
	if err != nil {
		return t, fmt.Errorf("ошибка при создании зависимости типа %v: %w", typ, err)
	}

	result, ok := instance.(T)
	if !ok {
		return t, fmt.Errorf("неверный тип возвращаемого значения для %v", typ)
	}

	// Для singleton сохраняем в кэше
	if descriptor.lifetime == Singleton {
		descriptor.mu.Lock()
		descriptor.instance = result
		descriptor.mu.Unlock()
	}

	return result, nil
}

// Contains проверяет, зарегистрирована ли зависимость для указанного типа
func Contains[T any](ctx *Context) bool {
	typ := getType[T]()

	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	_, exists := ctx.services[typ]
	return exists
}

// GetServices возвращает количество зарегистрированных зависимостей (для отладки)
func GetServices(ctx *Context) int {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	return len(ctx.services)
}

// === Реализация interface context.Context ===

// Deadline возвращает время deadline для этого контекста
// Если deadline не установлен, возвращает нулевое время и false
func (c *Context) Deadline() (time.Time, bool) {
	if c.deadline.IsZero() {
		return time.Time{}, false
	}
	return c.deadline, true
}

// Done возвращает канал, который закрывается при отмене контекста
func (c *Context) Done() <-chan struct{} {
	return c.done
}

// Err возвращает причину отмены контекста, или nil если контекст не отменен
func (c *Context) Err() error {
	c.mu2.RLock()
	defer c.mu2.RUnlock()
	return c.err
}

// Value возвращает значение, связанное с ключом
func (c *Context) Value(key any) any {
	c.mu2.RLock()
	defer c.mu2.RUnlock()
	return c.values[key]
}

// SetValue устанавливает значение для указанного ключа
func (c *Context) SetValue(key any, value any) {
	c.mu2.Lock()
	defer c.mu2.Unlock()
	c.values[key] = value
}

// Cancel отменяет контекст с ошибкой context.Canceled
func (c *Context) Cancel() {
	c.cancelWithErr(context.Canceled)
}

// CancelWithError отменяет контекст с указанной ошибкой
func (c *Context) CancelWithError(err error) {
	if err == nil {
		err = context.Canceled
	}
	c.cancelWithErr(err)
}

// cancelWithErr внутренний метод для отмены контекста
func (c *Context) cancelWithErr(err error) {
	c.mu2.Lock()
	if c.err != nil {
		c.mu2.Unlock()
		return // уже отменен
	}
	c.err = err
	c.mu2.Unlock()

	close(c.done)
}

// WithDeadline возвращает новый контекст с установленным deadline
func (c *Context) WithDeadline(deadline time.Time) (*Context, context.CancelFunc) {
	newCtx := &Context{
		services: c.services,
		parent:   c,
		deadline: deadline,
		done:     make(chan struct{}),
		values:   make(map[any]any),
	}

	// Копируем существующие значения
	c.mu2.RLock()
	maps.Copy(newCtx.values, c.values)
	c.mu2.RUnlock()

	// Запускаем горутину для отмены по deadline
	go func() {
		select {
		case <-newCtx.done:
		case <-c.Done():
			newCtx.CancelWithError(context.Canceled)
		case <-time.After(time.Until(deadline)):
			newCtx.CancelWithError(context.DeadlineExceeded)
		}
	}()

	return newCtx, func() { newCtx.Cancel() }
}

// WithTimeout возвращает новый контекст с установленным timeout
func (c *Context) WithTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	return c.WithDeadline(time.Now().Add(timeout))
}

// WithCancel возвращает новый контекст с функцией отмены
func (c *Context) WithCancel() (*Context, context.CancelFunc) {
	newCtx := &Context{
		services: c.services,
		parent:   c,
		done:     make(chan struct{}),
		values:   make(map[any]any),
	}

	// Копируем существующие значения
	c.mu2.RLock()
	maps.Copy(newCtx.values, c.values)
	c.mu2.RUnlock()

	// Отслеживаем отмену parent контекста
	go func() {
		select {
		case <-newCtx.done:
		case <-c.Done():
			newCtx.CancelWithError(c.Err())
		}
	}()

	return newCtx, func() { newCtx.Cancel() }
}

// WithValue возвращает новый контекст с установленным значением
func (c *Context) WithValue(key any, value any) *Context {
	newCtx := &Context{
		services: c.services,
		parent:   c,
		done:     make(chan struct{}),
		values:   make(map[any]any),
	}

	// Копируем существующие значения
	c.mu2.RLock()
	for k, v := range c.values {
		newCtx.values[k] = v
	}
	c.mu2.RUnlock()

	newCtx.SetValue(key, value)

	// Отслеживаем отмену parent контекста
	go func() {
		select {
		case <-newCtx.done:
		case <-c.Done():
			newCtx.CancelWithError(c.Err())
		}
	}()

	return newCtx
}
