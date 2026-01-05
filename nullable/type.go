package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Type представляет универсальный обнуляемый тип, который может содержать значение типа T или быть нулевым
type Type[T any] struct {
	V     T
	Valid bool
}

// New создает новый Type[T] со значением
func New[T any](value ...T) Type[T] {
	if len(value) > 1 {
		panic("nullable.New() принимает не более 1 значения")
	}
	if len(value) == 0 {
		var zero T
		return Type[T]{V: zero, Valid: false}
	}
	return Type[T]{V: value[0], Valid: true}
}

// IsNull проверяет, является ли Type[T] нулевым
func (t *Type[T]) IsNull() bool {
	return !t.Valid
}

// IsNotNull проверяет, имеет ли Type[T] значение
func (t *Type[T]) IsNotNull() bool {
	return t.Valid
}

// Set устанавливает значение и помечает его как действительное
func (t *Type[T]) Set(value T) {
	t.V = value
	t.Valid = true
}

// SetNull устанавливает значение в null
func (t *Type[T]) SetNull() {
	var zero T
	t.V = zero
	t.Valid = false
}

// Get возвращает значение и логическое значение, указывающее, является ли оно действительным
func (t Type[T]) Get() (T, bool) {
	return t.V, t.Valid
}

// Or возвращает значение, если оно действительное, в противном случае возвращает предоставленное по умолчанию
func (t Type[T]) Or(defaultValue T) T {
	if t.Valid {
		return t.V
	}
	return defaultValue
}

// String возвращает строковое представление Type[T]
func (t Type[T]) String() string {
	if t.Valid {
		return fmt.Sprintf("%v", t.V)
	}
	return "null"
}

// MarshalJSON реализует интерфейс json.Marshaler
func (t Type[T]) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.V)
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (t *Type[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var zero T
		t.V = zero
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.V); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// Value реализует интерфейс driver.Valuer
func (t Type[T]) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.V, nil
}

// Scan реализует интерфейс sql.Scanner
func (t *Type[T]) Scan(value any) error {
	if value == nil {
		var zero T
		t.V = zero
		t.Valid = false
		return nil
	}

	// Преобразование считанного значения в целевой тип
	convertedValue, ok := value.(T)
	if !ok {
		// Если прямое преобразование не удалось, возможно, нам нужно обработать конкретные преобразования типов
		// Например, преобразование из типов базы данных в целевой тип
		return fmt.Errorf("не удается сканировать %T в Type[%T]", value, t.V)
	}

	t.V = convertedValue
	t.Valid = true
	return nil
}
