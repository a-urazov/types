package big

import (
	"fmt"
	"math/big"
)

// Number constraint определяет типы, которые могут быть использованы в универсальном конструкторе
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | big.Int | big.Float | string
}

// Type представляет тип большого числа
type Type uint8

const (
	Integer Type = iota
	Decimal
)

// Big представляет большое число, которое может быть целым или десятичным
type Big struct {
	typ Type
	int *big.Int   // Используется, когда тип - Integer
	dec *big.Float // Используется, когда тип - Decimal
}

// New[T Number] создает новый объект Big на основе переданного значения
func New[T Number](value ...T) *Big {
	if len(value) == 0 {
		return newInteger(0)
	}
	switch any(value[0]).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// Преобразуем в int64 для создания целого числа
		return newInteger(toInt64(value[0]))
	case float32, float64:
		// Преобразуем в float64 для создания десятичного числа
		return newDecimal(toFloat64(value[0]))
	case big.Int:
		v := any(value[0]).(big.Int)
		return newIntegerFromBigInt(&v)
	case big.Float:
		v := any(value[0]).(big.Float)
		return newDecimalFromBigFloat(&v)
	case string:
		return parse(any(value[0]).(string))  // Вызов Parse, который нужно определить
	default:
		// Возвращаем нулевое значение для неожиданного типа
		return newInteger(0)
	}
}

// parse парсит строку и возвращает объект Big
func parse(s string) *Big {
    // Пробуем сначала распознать как целое число
    if i, ok := new(big.Int).SetString(s, 10); ok {
        return newIntegerFromBigInt(i)
    }

    // Если не удалось как целое, пробуем как десятичное
    if f, ok := new(big.Float).SetString(s); ok {
        return newDecimalFromBigFloat(f)
    }

    // Если оба способа не удались, возвращаем 0
    return newInteger(0)
}

// Вспомогательная функция для преобразования чисел в int64
func toInt64[T Number](value T) int64 {
	switch v := any(value).(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	default:
		return int64(0) // fallback
	}
}

// Вспомогательная функция для преобразования чисел в float64
func toFloat64[T Number](value T) float64 {
	switch v := any(value).(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		// Если это целое число, преобразуем его в float64
		return float64(toInt64(value))
	}
}

// newInteger создает новое большое целое число
func newInteger(value int64) *Big {
	return &Big{
		typ: Integer,
		int: big.NewInt(value),
	}
}

// newIntegerFromBigInt создает новое большое целое число из *big.Int
func newIntegerFromBigInt(value *big.Int) *Big {
	return &Big{
		typ: Integer,
		int: new(big.Int).Set(value),
	}
}

// newDecimal создает новое большое десятичное число из float64
func newDecimal(value float64) *Big {
	return &Big{
		typ: Decimal,
		dec: big.NewFloat(value),
	}
}

// newDecimalFromBigFloat создает новое большое десятичное число из *big.Float
func newDecimalFromBigFloat(value *big.Float) *Big {
	return &Big{
		typ: Decimal,
		dec: new(big.Float).Set(value),
	}
}

// newDecimalFromString создает новое большое десятичное число из строки
func newDecimalFromString(s string) (*Big, error) {
	f := new(big.Float)
	_, ok := f.SetString(s)
	if !ok {
		return nil, fmt.Errorf("invalid decimal string: %s", s)
	}
	return &Big{
		typ: Decimal,
		dec: f,
	}, nil
}

// Add складывает два больших числа
func (b *Big) Add(other *Big) *Big {
	// Convert to decimal for calculation to maintain consistent return type
	bDec := b.ToDecimal()
	otherDec := other.ToDecimal()
	result := new(big.Float).Add(bDec.dec, otherDec.dec)
	return newDecimalFromBigFloat(result)
}

// Subtract вычитает одно большое число из другого
func (b *Big) Subtract(other *Big) *Big {
	// Convert to decimal for calculation to maintain consistent return type
	bDec := b.ToDecimal()
	otherDec := other.ToDecimal()
	result := new(big.Float).Sub(bDec.dec, otherDec.dec)
	return newDecimalFromBigFloat(result)
}

// Multiply перемножает два больших числа
func (b *Big) Multiply(other *Big) *Big {
	// Convert to decimal for calculation to maintain consistent return type
	bDec := b.ToDecimal()
	otherDec := other.ToDecimal()
	result := new(big.Float).Mul(bDec.dec, otherDec.dec)
	return newDecimalFromBigFloat(result)
}

// Divide делит одно большое число на другое
func (b *Big) Divide(other *Big) *Big {
	if other.IsZero() {
		panic("division by zero")
	}

	// Convert to decimal for calculation to preserve precision
	bDec := b.ToDecimal()
	otherDec := other.ToDecimal()
	result := new(big.Float).Quo(bDec.dec, otherDec.dec)
	return newDecimalFromBigFloat(result)
}

// Mod calculates the modulo of two big numbers
func (b *Big) Mod(other *Big) *Big {
	// Convert to decimal for calculation to maintain consistent return type
	bDec := b.ToDecimal()
	otherDec := other.ToDecimal()

	// Calculate quotient
	quo := new(big.Float).Quo(bDec.dec, otherDec.dec)

	// Truncate to integer part
	quoInt, _ := quo.Int(nil)

	// Calculate remainder: remainder = dividend - divisor * quotient_int
	divisor := otherDec.dec
	product := new(big.Float).Mul(divisor, new(big.Float).SetInt(quoInt))
	result := new(big.Float).Sub(bDec.dec, product)

	return newDecimalFromBigFloat(result)
}

// Pow raises this big number to the power of the exponent
func (b *Big) Pow(exponent int64) *Big {
	// Convert to decimal for calculation to maintain consistent return type
	bDec := b.ToDecimal()
	result := new(big.Float).Set(bDec.dec)

	for i := int64(1); i < exponent; i++ {
		result.Mul(result, bDec.dec)
	}

	return newDecimalFromBigFloat(result)
}

// Compare compares two big numbers (-1 if less, 0 if equal, 1 if greater)
func (b *Big) Compare(other *Big) int {
	if b.typ == Integer && other.typ == Integer {
		return b.int.Cmp(other.int)
	} else {
		// Convert both to decimals for comparison
		bDec := b.ToDecimal()
		otherDec := other.ToDecimal()

		return bDec.dec.Cmp(otherDec.dec)
	}
}

// Sign returns -1 if the number is negative, 0 if zero, 1 if positive
func (b *Big) Sign() int {
	if b.typ == Integer {
		return b.int.Sign()
	} else {
		return b.dec.Sign()
	}
}

// IsZero checks if the number is zero
func (b *Big) IsZero() bool {
	return b.Sign() == 0
}

// IsPositive checks if the number is positive
func (b *Big) IsPositive() bool {
	return b.Sign() > 0
}

// IsNegative checks if the number is negative
func (b *Big) IsNegative() bool {
	return b.Sign() < 0
}

// ToInteger converts the number to a big integer (truncates decimal part if applicable)
func (b *Big) ToInteger() *big.Int {
	if b.typ == Integer {
		return new(big.Int).Set(b.int)
	} else {
		// Convert decimal to integer by truncating
		intVal, _ := b.dec.Int(nil)
		return intVal
	}
}

// ToDecimal converts the number to a big decimal
func (b *Big) ToDecimal() *Big {
	if b.typ == Decimal {
		return &Big{
			typ: Decimal,
			dec: new(big.Float).Set(b.dec),
		}
	} else {
		// Convert integer to decimal
		floatVal := new(big.Float).SetInt(b.int)
		return &Big{
			typ: Decimal,
			dec: floatVal,
		}
	}
}

// ToFloat64 converts the number to a float64
func (b *Big) ToFloat64() (float64, bool) {
	if b.typ == Integer {
		f, acc := b.int.Float64()
		return f, acc == big.Exact
	} else {
		f, acc := b.dec.Float64()
		return f, acc == big.Exact
	}
}

// ToInt64 converts the number to an int64
func (b *Big) ToInt64() (int64, bool) {
	if b.typ == Integer {
		return b.int.Int64(), true
	} else {
		// Convert decimal to int64 by truncating
		intVal, _ := b.dec.Int64()
		return intVal, true
	}
}

// String returns a string representation of the number
func (b *Big) String() string {
	if b.typ == Integer {
		return b.int.String()
	} else {
		return b.dec.Text('g', 10)
	}
}

// Equals checks if two big numbers are equal
func (b *Big) Equals(other *Big) bool {
	return b.Compare(other) == 0
}

// Clone creates a copy of the number
func (b *Big) Clone() *Big {
	if b.typ == Integer {
		return &Big{
			typ: Integer,
			int: new(big.Int).Set(b.int),
		}
	} else {
		return &Big{
			typ: Decimal,
			dec: new(big.Float).Set(b.dec),
		}
	}
}

// Sqrt computes the square root of the number
func (b *Big) Sqrt() *Big {
	if b.typ == Integer {
		// Convert to decimal for square root calculation
		dec := b.ToDecimal()
		return dec.Sqrt()
	} else {
		if b.dec.Sign() < 0 {
			panic("square root of negative number")
		}

		result := new(big.Float)
		result.Sqrt(b.dec)
		return newDecimalFromBigFloat(result)
	}
}

// Abs returns the absolute value of the number
func (b *Big) Abs() *Big {
	if b.typ == Integer {
		result := new(big.Int).Abs(b.int)
		return newIntegerFromBigInt(result)
	} else {
		result := new(big.Float).Abs(b.dec)
		return newDecimalFromBigFloat(result)
	}
}
