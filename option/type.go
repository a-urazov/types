package option

import (
	"encoding/json"
	"fmt"
)

// Option represents a value that may or may not be present.
// It's a functional alternative to nullable types, avoiding null pointer exceptions.
type Option[T any] struct {
	value *T
}

// Some creates an Option with a value present.
func Some[T any](value T) Option[T] {
	return Option[T]{value: &value}
}

// None creates an Option with no value present.
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// IsSome checks if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone checks if the Option does not contain a value.
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Get returns the value and a boolean indicating if it's present.
func (o Option[T]) Get() (T, bool) {
	if o.IsSome() {
		return *o.value, true
	}
	var zero T
	return zero, false
}

// GetOrElse returns the value if present, otherwise returns the provided default.
func (o Option[T]) GetOrElse(defaultValue T) T {
	if o.IsSome() {
		return *o.value
	}
	return defaultValue
}

// GetOrCall returns the value if present, otherwise returns the result of calling the provided function.
func (o Option[T]) GetOrCall(f func() T) T {
	if o.IsSome() {
		return *o.value
	}
	return f()
}

// Map applies a function to the value if present, returning a new Option.
// If the Option is None, it returns None.
func Map[T, U any](opt Option[T], fn func(T) U) Option[U] {
	if opt.IsSome() {
		value, _ := opt.Get()
		return Some(fn(value))
	}
	return None[U]()
}

// FlatMap applies a function that returns an Option to the value if present.
// If the Option is None, it returns None.
func FlatMap[T, U any](opt Option[T], fn func(T) Option[U]) Option[U] {
	if opt.IsSome() {
		value, _ := opt.Get()
		return fn(value)
	}
	return None[U]()
}

// Filter returns the Option if the value satisfies the predicate, otherwise returns None.
func Filter[T any](opt Option[T], predicate func(T) bool) Option[T] {
	if opt.IsSome() {
		value, _ := opt.Get()
		if predicate(value) {
			return opt
		}
	}
	return None[T]()
}

// String returns a string representation of the Option.
func (o Option[T]) String() string {
	if o.IsSome() {
		value, _ := o.Get()
		return fmt.Sprintf("Some(%v)", value)
	}
	return "None"
}

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsSome() {
		value, _ := o.Get()
		return json.Marshal(value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		o.value = nil
		return err
	}

	o.value = &value
	return nil
}