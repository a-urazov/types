package either

import (
	"encoding/json"
	"fmt"
)

// Either представляет значение одного из двух возможных типов, Left или Right.
// По соглашению, Right используется для успеха, а Left - для ошибки.
type Either[L any, R any] struct {
	left    L
	right   R
	isRight bool
}

// Left создает новый Either со значением Left.
func Left[L any, R any](value L) Either[L, R] {
	return Either[L, R]{left: value, right: *new(R), isRight: false}
}

// Right создает новый Either со значением Right.
func Right[L any, R any](value R) Either[L, R] {
	return Either[L, R]{left: *new(L), right: value, isRight: true}
}

// IsLeft проверяет, содержит ли Either значение Left.
func (e Either[L, R]) IsLeft() bool {
	return !e.isRight
}

// IsRight checks if the Either contains a Right value.
func (e Either[L, R]) IsRight() bool {
	return e.isRight
}

// GetLeft returns the Left value and a boolean indicating if it's valid.
// The boolean will be false if the Either contains a Right value.
func (e Either[L, R]) GetLeft() (L, bool) {
	if e.IsLeft() {
		return e.left, true
	}
	var zero L
	return zero, false
}

// GetRight returns the Right value and a boolean indicating if it's valid.
// The boolean will be false if the Either contains a Left value.
func (e Either[L, R]) GetRight() (R, bool) {
	if e.IsRight() {
		return e.right, true
	}
	var zero R
	return zero, false
}

// OrLeft returns the Left value if present, otherwise returns the provided default.
func (e Either[L, R]) OrLeft(defaultValue L) L {
	if e.IsLeft() {
		return e.left
	}
	return defaultValue
}

// OrRight returns the Right value if present, otherwise returns the provided default.
func (e Either[L, R]) OrRight(defaultValue R) R {
	if e.IsRight() {
		return e.right
	}
	return defaultValue
}

// MapRight applies a function to the Right value if present, returning a new Either.
// If the Either contains a Left value, it returns the same Left value with a transformed Right type.
func MapRight[L, R, T any](e Either[L, R], fn func(R) T) Either[L, T] {
	if e.IsRight() {
		return Right[L, T](fn(e.right))
	}
	// We only need the left value, the T type is ignored in Left case
	return Left[L, T](e.left)
}

// FlatMapRight applies a function that returns an Either to the Right value if present.
// If the Either contains a Left value, it returns a Left with the same left value type.
func FlatMapRight[L, R, T any](e Either[L, R], fn func(R) Either[L, T]) Either[L, T] {
	if e.IsRight() {
		return fn(e.right)
	}
	// Preserve the original left value
	return Left[L, T](e.left)
}

// MapLeft applies a function to the Left value if present, returning a new Either.
// If the Either contains a Right value, it returns the same Right value with a transformed Left type.
func MapLeft[L, R, T any](e Either[L, R], fn func(L) T) Either[T, R] {
	if e.IsLeft() {
		return Left[T, R](fn(e.left))
	}
	// We only need the right value, the T type is ignored in Right case
	return Right[T, R](e.right)
}

// FlatMapLeft applies a function that returns an Either to the Left value if present.
// If the Either contains a Right value, it returns a Right with the same right value type.
func FlatMapLeft[L, R, T any](e Either[L, R], fn func(L) Either[T, R]) Either[T, R] {
	if e.IsLeft() {
		return fn(e.left)
	}
	// Preserve the original right value
	return Right[T, R](e.right)
}

// String returns a string representation of the Either.
func (e Either[L, R]) String() string {
	if e.IsRight() {
		return fmt.Sprintf("Right(%v)", e.right)
	}
	return fmt.Sprintf("Left(%v)", e.left)
}

// MarshalJSON implements the json.Marshaler interface.
func (e Either[L, R]) MarshalJSON() ([]byte, error) {
	if e.IsRight() {
		// Create a wrapper object with a "right" field
		wrapper := map[string]any{"right": e.right}
		return json.Marshal(wrapper)
	}
	// Create a wrapper object with a "left" field
	wrapper := map[string]any{"left": e.left}
	return json.Marshal(wrapper)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Either[L, R]) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a Right value first
	var rightWrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &rightWrapper); err != nil {
		return err
	}

	if rightData, exists := rightWrapper["right"]; exists {
		var rightVal R
		if err := json.Unmarshal(rightData, &rightVal); err != nil {
			return err
		}
		e.right = rightVal
		e.isRight = true
		// Initialize left to zero value
		e.left = *new(L)
		return nil
	}

	if leftData, exists := rightWrapper["left"]; exists {
		var leftVal L
		if err := json.Unmarshal(leftData, &leftVal); err != nil {
			return err
		}
		e.left = leftVal
		e.isRight = false
		// Initialize right to zero value
		e.right = *new(R)
		return nil
	}

	return fmt.Errorf("invalid Either format")
}
