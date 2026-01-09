package result

// Result представляет тип, который может содержать либо значение типа T, либо ошибку
type Result[T any] struct {
	value T
	err   error
}

// Ok создает Result с успешным значением
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err создает Result с ошибкой
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Ok проверяет, содержит ли Result успешное значение
func (r Result[T]) Ok() bool {
	return r.err == nil
}

// Err проверяет, содержит ли Result ошибку
func (r Result[T]) Err() bool {
	return r.err != nil
}

// Unwrap возвращает значение и ошибку
func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// UnwrapOr returns the value if successful, otherwise returns the default value
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value if successful, otherwise returns the result of the provided function
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.err != nil {
		return fn(r.err)
	}
	return r.value
}

// Expect returns the value if successful, otherwise panics with the provided message
func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(msg + ": " + r.err.Error())
	}
	return r.value
}

// UnwrapErr returns the error if present, otherwise panics
func (r Result[T]) UnwrapErr() error {
	if r.err == nil {
		panic("вызван UnwrapErr для значения Ok")
	}
	return r.err
}

// And returns the other Result if this one is Ok, otherwise returns the error
func (r Result[T]) And(other Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return other
}

// AndThen calls the provided function with the value if this Result is Ok
func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return fn(r.value)
}

// Or returns the other Result if this one is Err, otherwise returns this Result
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.err != nil {
		return other
	}
	return r
}

// OrElse returns the result of the provided function if this Result is Err
func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if r.err != nil {
		return fn(r.err)
	}
	return r
}

// Map applies a function to the value if the Result is Ok
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(fn(r.value))
}

// MapOr applies a function to the value if the Result is Ok, otherwise returns the default value
func (r Result[T]) MapOr(defaultValue T, fn func(T) T) T {
	if r.err != nil {
		return defaultValue
	}
	return fn(r.value)
}

// MapOrElse applies a function to the value if the Result is Ok, otherwise returns the result of the fallback function
func (r Result[T]) MapOrElse(fallbackFn func(error) T, fn func(T) T) T {
	if r.err != nil {
		return fallbackFn(r.err)
	}
	return fn(r.value)
}

// Match provides pattern matching for Result, similar to Rust implementation
func (r Result[T]) Match(okFn func(T) any, errFn func(error) any) any {
	if r.err != nil {
		return errFn(r.err)
	} else {
		return okFn(r.value)
	}
}

// Error converts an error to a Result
func Error[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Func executes a function that returns a value and error, and wraps them in a Result
func Func[T any](fn func() (T, error)) Result[T] {
	value, err := fn()
	return Error(value, err)
}
