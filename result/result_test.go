package result

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	value := 42
	r := Ok(value)

	if !r.Ok() {
		t.Errorf("Ok(%d) should be Ok, got Err", value)
	}

	if r.Err() {
		t.Errorf("Ok(%d) should not be Err", value)
	}

	resultValue, err := r.Unwrap()
	if resultValue != value {
		t.Errorf("Expected value %d, got %d", value, resultValue)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestErr(t *testing.T) {
	testErr := errors.New("test error")
	r := Err[int](testErr)

	if r.Ok() {
		t.Errorf("Err() should not be Ok")
	}

	if !r.Err() {
		t.Errorf("Err() should be Err")
	}

	_, err := r.Unwrap()
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestUnwrapOr(t *testing.T) {
	value := 42
	defaultValue := 0

	// Test with Ok value
	r1 := Ok(value)
	result := r1.UnwrapOr(defaultValue)
	if result != value {
		t.Errorf("Expected %d, got %d", value, result)
	}

	// Test with Err value
	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	result = r2.UnwrapOr(defaultValue)
	if result != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, result)
	}
}

func TestUnwrapOrElse(t *testing.T) {
	value := 42

	// Test with Ok value
	r1 := Ok(value)
	result := r1.UnwrapOrElse(func(err error) int {
		return 0 // Should not be called
	})
	if result != value {
		t.Errorf("Expected %d, got %d", value, result)
	}

	// Test with Err value
	defaultValue := 100
	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	result = r2.UnwrapOrElse(func(err error) int {
		if err != testErr {
			t.Errorf("Expected error %v, got %v", testErr, err)
		}
		return defaultValue
	})
	if result != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, result)
	}
}

func TestExpect(t *testing.T) {
	value := 42
	r := Ok(value)

	result := r.Expect("Unexpected error")
	if result != value {
		t.Errorf("Expected %d, got %d", value, result)
	}

	// Test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expect on Err should panic")
		}
	}()

	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	r2.Expect("This should panic")
}

func TestUnwrapErr(t *testing.T) {
	testErr := errors.New("test error")
	r := Err[int](testErr)

	resultErr := r.UnwrapErr()
	if resultErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, resultErr)
	}

	// Test panic case for Ok value
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("UnwrapErr on Ok should panic")
		}
	}()

	r2 := Ok(42)
	r2.UnwrapErr()
}

func TestAnd(t *testing.T) {
	value1 := 10
	value2 := 20

	// Test Ok and Ok
	r1 := Ok(value1)
	r2 := Ok(value2)
	result := r1.And(r2)

	if !result.Ok() {
		t.Errorf("Ok and Ok should be Ok")
	}

	resultValue, _ := result.Unwrap()
	if resultValue != value2 {
		t.Errorf("Expected %d, got %d", value2, resultValue)
	}

	// Test Err and Ok
	err1 := errors.New("error 1")
	r3 := Err[int](err1)
	r4 := Ok(30)
	result = r3.And(r4)

	if !result.Err() {
		t.Errorf("Err and Ok should be Err")
	}

	_, resultErr := result.Unwrap()
	if resultErr != err1 {
		t.Errorf("Expected error %v, got %v", err1, resultErr)
	}
}

func TestAndThen(t *testing.T) {
	value := 10

	// Test with Ok value
	r1 := Ok(value)
	result := r1.AndThen(func(v int) Result[int] {
		return Ok(v * 2)
	})

	if !result.Ok() {
		t.Errorf("AndThen with Ok should return Ok")
	}

	resultValue, _ := result.Unwrap()
	expected := value * 2
	if resultValue != expected {
		t.Errorf("Expected %d, got %d", expected, resultValue)
	}

	// Test with Err value
	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	result = r2.AndThen(func(v int) Result[int] {
		return Ok(v * 2) // Should not be called
	})

	if !result.Err() {
		t.Errorf("AndThen with Err should return Err")
	}

	_, resultErr := result.Unwrap()
	if resultErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, resultErr)
	}
}

func TestOr(t *testing.T) {
	value1 := 10
	value2 := 20

	// Test Ok or anything (should return first)
	r1 := Ok(value1)
	r2 := Ok(value2)
	result := r1.Or(r2)

	if !result.Ok() {
		t.Errorf("Ok or anything should return first Ok")
	}

	resultValue, _ := result.Unwrap()
	if resultValue != value1 {
		t.Errorf("Expected %d, got %d", value1, resultValue)
	}

	// Test Err or Ok
	testErr := errors.New("test error")
	r3 := Err[int](testErr)
	r4 := Ok(value2)
	result = r3.Or(r4)

	if !result.Ok() {
		t.Errorf("Err or Ok should return second Ok")
	}

	resultValue, _ = result.Unwrap()
	if resultValue != value2 {
		t.Errorf("Expected %d, got %d", value2, resultValue)
	}
}

func TestOrElse(t *testing.T) {
	value := 10

	// Test with Ok value (function should not be called)
	r1 := Ok(value)
	result := r1.OrElse(func(err error) Result[int] {
		return Ok(999) // Should not be called
	})

	if !result.Ok() {
		t.Errorf("OrElse with Ok should return original Ok")
	}

	resultValue, _ := result.Unwrap()
	if resultValue != value {
		t.Errorf("Expected %d, got %d", value, resultValue)
	}

	// Test with Err value
	testErr := errors.New("test error")
	defaultValue := 50
	r2 := Err[int](testErr)
	result = r2.OrElse(func(err error) Result[int] {
		if err != testErr {
			t.Errorf("Expected error %v, got %v", testErr, err)
		}
		return Ok(defaultValue)
	})

	if !result.Ok() {
		t.Errorf("OrElse with Err should return result of function")
	}

	resultValue, _ = result.Unwrap()
	if resultValue != defaultValue {
		t.Errorf("Expected %d, got %d", defaultValue, resultValue)
	}
}

func TestMap(t *testing.T) {
	initialValue := 10
	expectedValue := initialValue * 2

	// Test with Ok value
	r := Ok(initialValue)
	result := r.Map(func(v int) int {
		return v * 2
	})

	if !result.Ok() {
		t.Errorf("Map with Ok should return Ok")
	}

	resultValue, _ := result.Unwrap()
	if resultValue != expectedValue {
		t.Errorf("Expected %d, got %d", expectedValue, resultValue)
	}

	// Test with Err value
	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	result = r2.Map(func(v int) int {
		return v * 2 // Should not be called
	})

	if !result.Err() {
		t.Errorf("Map with Err should return Err")
	}

	_, resultErr := result.Unwrap()
	if resultErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, resultErr)
	}
}

func TestMapOr(t *testing.T) {
	initialValue := 10
	defaultValue := 0
	expectedValue := initialValue + 5

	// Test with Ok value
	r := Ok(initialValue)
	result := r.MapOr(defaultValue, func(v int) int {
		return v + 5
	})

	if result != expectedValue {
		t.Errorf("Expected %d, got %d", expectedValue, result)
	}

	// Test with Err value
	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	result = r2.MapOr(defaultValue, func(v int) int {
		return v + 5 // Should not be called
	})

	if result != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, result)
	}
}

func TestMapOrElse(t *testing.T) {
	initialValue := 10
	expectedValue := initialValue + 5

	// Test with Ok value
	r := Ok(initialValue)
	result := r.MapOrElse(
		func(err error) int { return 999 }, // Should not be called
		func(v int) int { return v + 5 },
	)

	if result != expectedValue {
		t.Errorf("Expected %d, got %d", expectedValue, result)
	}

	// Test with Err value
	testErr := errors.New("test error")
	defaultValue := 20
	r2 := Err[int](testErr)
	result = r2.MapOrElse(
		func(err error) int {
			if err != testErr {
				t.Errorf("Expected error %v, got %v", testErr, err)
			}
			return defaultValue
		},
		func(v int) int { return v + 5 }, // Should not be called
	)

	if result != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, result)
	}
}

func TestMatch(t *testing.T) {
	// Test with Ok value
	var okValue int
	var errReceived error

	r1 := Ok(42)
	r1.Match(
		func(v int) any { okValue = v; return nil },
		func(err error) any { errReceived = err; return nil },
	)

	if okValue != 42 {
		t.Errorf("Expected okValue to be 42, got %d", okValue)
	}

	if errReceived != nil {
		t.Errorf("Expected no error, got %v", errReceived)
	}

	// Test with Err value
	okValue = 0
	errReceived = nil

	testErr := errors.New("test error")
	r2 := Err[int](testErr)
	r2.Match(
		func(v int) any { okValue = v; return nil },
		func(err error) any { errReceived = err; return nil },
	)

	if okValue != 0 {
		t.Errorf("Expected okValue to be 0, got %d", okValue)
	}

	if errReceived != testErr {
		t.Errorf("Expected error %v, got %v", testErr, errReceived)
	}
}

func TestFromError(t *testing.T) {
	value := 42

	// Test with nil error
	r1 := Error(value, nil)

	if !r1.Ok() {
		t.Errorf("FromError with nil error should be Ok")
	}

	resultValue, err := r1.Unwrap()
	if resultValue != value {
		t.Errorf("Expected value %d, got %d", value, resultValue)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test with non-nil error
	testErr := errors.New("test error")
	r2 := Error(value, testErr)

	if !r2.Err() {
		t.Errorf("FromError with error should be Err")
	}

	_, resultErr := r2.Unwrap()
	if resultErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, resultErr)
	}
}

func TestFromFunc(t *testing.T) {
	value := 42

	// Test function that returns no error
	fn1 := func() (int, error) {
		return value, nil
	}
	r1 := Func(fn1)

	if !r1.Ok() {
		t.Errorf("FromFunc with function returning nil error should be Ok")
	}

	resultValue, err := r1.Unwrap()
	if resultValue != value {
		t.Errorf("Expected value %d, got %d", value, resultValue)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test function that returns an error
	testErr := errors.New("test error")
	fn2 := func() (int, error) {
		return value, testErr
	}
	r2 := Func(fn2)

	if !r2.Err() {
		t.Errorf("FromFunc with function returning error should be Err")
	}

	_, resultErr := r2.Unwrap()
	if resultErr != testErr {
		t.Errorf("Expected error %v, got %v", testErr, resultErr)
	}
}
