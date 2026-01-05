package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Тестовые интерфейсы и структуры
type Logger interface {
	Log(msg string)
}

type SimpleLogger struct {
	name string
}

func (l *SimpleLogger) Log(msg string) {
	// Простая реализация для тестирования
}

type Database interface {
	Connect() error
}

type MockDatabase struct {
	connected bool
}

func (db *MockDatabase) Connect() error {
	db.connected = true
	return nil
}

type Service struct {
	logger Logger
	db     Database
}

// TestRegisterAndResolveSingleton проверяет регистрацию и разрешение Singleton зависимостей
func TestRegisterAndResolveSingleton(t *testing.T) {
	ctx := New()

	// Регистрируем Logger как Singleton
	err := RegisterSingleton[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{name: "test"}, nil
	})
	if err != nil {
		t.Fatalf("ошибка при регистрации: %v", err)
	}

	// Разрешаем Logger дважды
	logger1, err := Resolve[Logger](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении: %v", err)
	}

	logger2, err := Resolve[Logger](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении: %v", err)
	}

	// Проверяем, что это один и тот же экземпляр
	if logger1 != logger2 {
		t.Error("Singleton должен возвращать один и тот же экземпляр")
	}
}

// TestRegisterAndResolveTransient проверяет регистрацию и разрешение Transient зависимостей
func TestRegisterAndResolveTransient(t *testing.T) {
	ctx := New()

	// Регистрируем Logger как Transient
	err := RegisterTransient[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{name: "test"}, nil
	})
	if err != nil {
		t.Fatalf("ошибка при регистрации: %v", err)
	}

	// Разрешаем Logger дважды
	logger1, err := Resolve[Logger](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении: %v", err)
	}

	logger2, err := Resolve[Logger](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении: %v", err)
	}

	// Проверяем, что это разные экземпляры
	if logger1 == logger2 {
		t.Error("Transient должен возвращать разные экземпляры")
	}
}

// TestRegisterInstance проверяет регистрацию конкретного экземпляра
func TestRegisterInstance(t *testing.T) {
	ctx := New()

	logger := &SimpleLogger{name: "singleton"}

	err := RegisterInstance[Logger](ctx, logger)
	if err != nil {
		t.Fatalf("ошибка при регистрации: %v", err)
	}

	resolved, err := Resolve[Logger](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении: %v", err)
	}

	// Проверяем, что это точно тот же экземпляр
	if resolved != logger {
		t.Error("RegisterInstance должен возвращать точно тот же экземпляр")
	}
}

// TestResolveMissingDependency проверяет обработку отсутствующей зависимости
func TestResolveMissingDependency(t *testing.T) {
	ctx := New()

	_, err := Resolve[Logger](ctx)
	if err == nil {
		t.Error("должна быть ошибка при разрешении неизвестной зависимости")
	}
}

// TestDuplicateRegistration проверяет, что нельзя зарегистрировать зависимость дважды
func TestDuplicateRegistration(t *testing.T) {
	ctx := New()

	constructor := func(c *Context) (interface{}, error) {
		return &SimpleLogger{}, nil
	}

	err1 := RegisterSingleton[Logger](ctx, constructor)
	if err1 != nil {
		t.Fatalf("первая регистрация должна пройти: %v", err1)
	}

	err2 := RegisterSingleton[Logger](ctx, constructor)
	if err2 == nil {
		t.Error("вторая регистрация должна вызвать ошибку")
	}
}

// TestContains проверяет проверку наличия зависимости
func TestContains(t *testing.T) {
	ctx := New()

	if Contains[Logger](ctx) {
		t.Error("Logger не должен быть зарегистрирован")
	}

	RegisterSingleton[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{}, nil
	})

	if !Contains[Logger](ctx) {
		t.Error("Logger должен быть зарегистрирован")
	}
}

// TestGetServices проверяет получение количества зарегистрированных зависимостей
func TestGetServices(t *testing.T) {
	ctx := New()

	if count := GetServices(ctx); count != 0 {
		t.Errorf("изначально должно быть 0 зависимостей, получено %d", count)
	}

	RegisterSingleton[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{}, nil
	})

	if count := GetServices(ctx); count != 1 {
		t.Errorf("должна быть 1 зависимость, получено %d", count)
	}

	RegisterSingleton[Database](ctx, func(c *Context) (interface{}, error) {
		return &MockDatabase{}, nil
	})

	if count := GetServices(ctx); count != 2 {
		t.Errorf("должны быть 2 зависимости, получено %d", count)
	}
}

// TestNestedDependencies проверяет разрешение вложенных зависимостей
func TestNestedDependencies(t *testing.T) {
	ctx := New()

	// Регистрируем Logger
	RegisterSingleton[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{name: "logger"}, nil
	})

	// Регистрируем Database
	RegisterSingleton[Database](ctx, func(c *Context) (interface{}, error) {
		return &MockDatabase{}, nil
	})

	// Регистрируем Service, который зависит от Logger и Database
	RegisterSingleton[*Service](ctx, func(c *Context) (interface{}, error) {
		logger, err := Resolve[Logger](c)
		if err != nil {
			return nil, err
		}

		db, err := Resolve[Database](c)
		if err != nil {
			return nil, err
		}

		return &Service{
			logger: logger,
			db:     db,
		}, nil
	})

	// Разрешаем Service
	service, err := Resolve[*Service](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении Service: %v", err)
	}

	if service.logger == nil {
		t.Error("logger не должен быть nil")
	}

	if service.db == nil {
		t.Error("database не должна быть nil")
	}
}

// BenchmarkResolve измеряет производительность разрешения зависимостей
func BenchmarkResolve(b *testing.B) {
	ctx := New()
	RegisterSingleton[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{name: "benchmark"}, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve[Logger](ctx)
	}
}

// BenchmarkResolveTransient измеряет производительность разрешения Transient зависимостей
func BenchmarkResolveTransient(b *testing.B) {
	ctx := New()
	RegisterTransient[Logger](ctx, func(c *Context) (interface{}, error) {
		return &SimpleLogger{name: "benchmark"}, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Resolve[Logger](ctx)
	}
}

// === Тесты для функций context.Context ===

// TestContextDeadline проверяет функцию Deadline
func TestContextDeadline(t *testing.T) {
	ctx := New()

	// По умолчанию deadline не установлен
	if deadline, ok := ctx.Deadline(); ok {
		t.Errorf("deadline не должен быть установлен, получено: %v", deadline)
	}

	// Установим deadline
	expectedDeadline := time.Now().Add(time.Hour)
	newCtx, _ := ctx.WithDeadline(expectedDeadline)

	if deadline, ok := newCtx.Deadline(); !ok {
		t.Error("deadline должен быть установлен")
	} else if deadline != expectedDeadline {
		t.Errorf("неверный deadline, ожидалось %v, получено %v", expectedDeadline, deadline)
	}
}

// TestContextCancel проверяет функцию Cancel
func TestContextCancel(t *testing.T) {
	ctx := New()

	if err := ctx.Err(); err != nil {
		t.Errorf("ошибка не должна быть установлена, получено: %v", err)
	}

	// Отменяем контекст
	ctx.Cancel()

	if err := ctx.Err(); err != context.Canceled {
		t.Errorf("ошибка должна быть context.Canceled, получено: %v", err)
	}

	// Проверяем, что канал Done закрыт
	select {
	case <-ctx.Done():
		// Корректно
	default:
		t.Error("канал Done должен быть закрыт")
	}
}

// TestContextWithCancel проверяет WithCancel
func TestContextWithCancel(t *testing.T) {
	ctx := New()
	newCtx, cancel := ctx.WithCancel()

	// Отменяем дочерний контекст
	cancel()

	if err := newCtx.Err(); err != context.Canceled {
		t.Errorf("ошибка должна быть context.Canceled, получено: %v", err)
	}

	// Родительский контекст должен остаться активным
	if err := ctx.Err(); err != nil {
		t.Errorf("родительский контекст не должен быть отменен, получено: %v", err)
	}
}

// TestContextWithTimeout проверяет WithTimeout
func TestContextWithTimeout(t *testing.T) {
	ctx := New()
	newCtx, cancel := ctx.WithTimeout(100 * time.Millisecond)
	defer cancel()

	// Проверяем, что контекст еще не отменен
	if err := newCtx.Err(); err != nil {
		t.Errorf("контекст не должен быть отменен сразу, получено: %v", err)
	}

	// Ждем истечения timeout
	time.Sleep(150 * time.Millisecond)

	if err := newCtx.Err(); err != context.DeadlineExceeded {
		t.Errorf("ошибка должна быть context.DeadlineExceeded, получено: %v", err)
	}
}

// TestContextValue проверяет Value и SetValue
func TestContextValue(t *testing.T) {
	ctx := New()

	// По умолчанию значение nil
	if val := ctx.Value("key"); val != nil {
		t.Errorf("значение должно быть nil, получено: %v", val)
	}

	// Устанавливаем значение
	ctx.SetValue("key", "value")

	if val := ctx.Value("key"); val != "value" {
		t.Errorf("значение должно быть 'value', получено: %v", val)
	}
}

// TestContextWithValue проверяет WithValue
func TestContextWithValue(t *testing.T) {
	ctx := New()
	ctx.SetValue("parentKey", "parentValue")

	newCtx := ctx.WithValue("childKey", "childValue")

	// Дочерний контекст должен иметь оба значения
	if val := newCtx.Value("parentKey"); val != "parentValue" {
		t.Errorf("значение parentKey должно быть 'parentValue', получено: %v", val)
	}

	if val := newCtx.Value("childKey"); val != "childValue" {
		t.Errorf("значение childKey должно быть 'childValue', получено: %v", val)
	}

	// Родительский контекст не должен иметь childKey
	if val := ctx.Value("childKey"); val != nil {
		t.Errorf("родительский контекст не должен иметь childKey, получено: %v", val)
	}
}

// TestContextWithValueInheritance проверяет наследование значений
func TestContextWithValueInheritance(t *testing.T) {
	ctx := New()
	ctx.SetValue("key1", "value1")

	newCtx1 := ctx.WithValue("key2", "value2")
	newCtx2 := newCtx1.WithValue("key3", "value3")

	// newCtx2 должен иметь все три значения
	if val := newCtx2.Value("key1"); val != "value1" {
		t.Errorf("key1 должно быть 'value1', получено: %v", val)
	}

	if val := newCtx2.Value("key2"); val != "value2" {
		t.Errorf("key2 должно быть 'value2', получено: %v", val)
	}

	if val := newCtx2.Value("key3"); val != "value3" {
		t.Errorf("key3 должно быть 'value3', получено: %v", val)
	}
}

// TestContextDIWithContextMethods проверяет DI контейнер с методами context.Context
func TestContextDIWithContextMethods(t *testing.T) {
	ctx := New()
	ctx.SetValue("appName", "TestApp")

	// Регистрируем сервис, который использует значение из контекста
	RegisterSingleton[*Service](ctx, func(c *Context) (interface{}, error) {
		// Получаем значение из контекста
		appName := c.Value("appName")

		return &Service{
			logger: &SimpleLogger{name: appName.(string)},
		}, nil
	})

	service, err := Resolve[*Service](ctx)
	if err != nil {
		t.Fatalf("ошибка при разрешении Service: %v", err)
	}

	if service.logger == nil {
		t.Error("logger не должен быть nil")
	}
}

// TestParentContextPropagation проверяет распространение отмены от parent контекста
func TestParentContextPropagation(t *testing.T) {
	ctx := New()
	newCtx, _ := ctx.WithCancel()

	// Отменяем parent контекст
	ctx.Cancel()

	// Ждем распространения отмены
	time.Sleep(10 * time.Millisecond)

	if err := newCtx.Err(); err == nil {
		t.Error("дочерний контекст должен быть отменен")
	}
}

// TestCancelWithError проверяет CancelWithError
func TestCancelWithError(t *testing.T) {
	ctx := New()
	customErr := fmt.Errorf("custom error")

	ctx.CancelWithError(customErr)

	if err := ctx.Err(); err != customErr {
		t.Errorf("ошибка должна быть customErr, получено: %v", err)
	}
}

// TestNewWithContext проверяет создание контекста с parent контекстом
func TestNewWithContext(t *testing.T) {
	parentCtx := context.Background()
	ctx := NewWithContext(parentCtx)

	if ctx == nil {
		t.Error("контекст не должен быть nil")
	}

	if err := ctx.Err(); err != nil {
		t.Errorf("контекст не должен быть отменен, получено: %v", err)
	}
}
