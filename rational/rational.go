package rational

import (
	"fmt"
	"math/big"
)

// Rational represents a rational number as a fraction of two integers (numerator/denominator)
type Rational struct {
	num *big.Int // numerator
	den *big.Int // denominator
}

// New creates a new rational number with the given numerator and denominator
func New(num, den int64) *Rational {
	if den == 0 {
		panic("знаменатель не может быть равен нулю")
	}

	n := big.NewInt(num)
	d := big.NewInt(den)

	// Handle negative denominators by moving the sign to the numerator
	if d.Sign() < 0 {
		n.Neg(n)
		d.Neg(d)
	}

	r := &Rational{num: n, den: d}
	return r.Reduce()
}

// NewFromBigInt creates a new rational number from big.Int values
func NewFromBigInt(num, den *big.Int) *Rational {
	if den.Sign() == 0 {
		panic("знаменатель не может быть равен нулю")
	}

	// Handle negative denominators by moving the sign to the numerator
	newNum := new(big.Int).Set(num)
	newDen := new(big.Int).Set(den)

	if newDen.Sign() < 0 {
		newNum.Neg(newNum)
		newDen.Neg(newDen)
	}

	r := &Rational{num: newNum, den: newDen}
	return r.Reduce()
}

// NewFromInt creates a rational number from an integer
func NewFromInt(value int64) *Rational {
	return &Rational{
		num: big.NewInt(value),
		den: big.NewInt(1),
	}
}

// Add adds two rational numbers
func (r *Rational) Add(other *Rational) *Rational {
	// a/b + c/d = (ad + bc) / bd
	num := new(big.Int).Mul(r.num, other.den)
	temp := new(big.Int).Mul(r.den, other.num)
	num.Add(num, temp)

	den := new(big.Int).Mul(r.den, other.den)

	return NewFromBigInt(num, den)
}

// Subtract subtracts the other rational number from this one
func (r *Rational) Subtract(other *Rational) *Rational {
	// a/b - c/d = (ad - bc) / bd
	num := new(big.Int).Mul(r.num, other.den)
	temp := new(big.Int).Mul(r.den, other.num)
	num.Sub(num, temp)

	den := new(big.Int).Mul(r.den, other.den)

	return NewFromBigInt(num, den)
}

// Multiply multiplies two rational numbers
func (r *Rational) Multiply(other *Rational) *Rational {
	// a/b * c/d = ac / bd
	num := new(big.Int).Mul(r.num, other.num)
	den := new(big.Int).Mul(r.den, other.den)

	return NewFromBigInt(num, den)
}

// Divide divides this rational number by the other
func (r *Rational) Divide(other *Rational) *Rational {
	if other.IsZero() {
		panic("деление на ноль")
	}

	// a/b ÷ c/d = ad / bc
	num := new(big.Int).Mul(r.num, other.den)
	den := new(big.Int).Mul(r.den, other.num)

	return NewFromBigInt(num, den)
}

// Invert returns the multiplicative inverse (reciprocal) of the rational number
func (r *Rational) Invert() *Rational {
	if r.IsZero() {
		panic("нельзя инвертировать ноль")
	}

	return NewFromBigInt(r.den, r.num)
}

// Negate returns the additive inverse of the rational number
func (r *Rational) Negate() *Rational {
	num := new(big.Int).Neg(r.num)
	return NewFromBigInt(num, r.den)
}

// Abs returns the absolute value of the rational number
func (r *Rational) Abs() *Rational {
	num := new(big.Int).Abs(r.num)
	// Denominator is always positive due to normalization
	return NewFromBigInt(num, r.den)
}

// Compare compares two rational numbers (-1 if less, 0 if equal, 1 if greater)
func (r *Rational) Compare(other *Rational) int {
	// a/b cmp c/d is equivalent to (ad - bc) cmp 0
	left := new(big.Int).Mul(r.num, other.den)
	right := new(big.Int).Mul(r.den, other.num)

	return left.Cmp(right)
}

// Sign returns -1 if the rational number is negative, 0 if zero, 1 if positive
func (r *Rational) Sign() int {
	return r.num.Sign()
}

// IsZero checks if the rational number is zero
func (r *Rational) IsZero() bool {
	return r.num.Sign() == 0
}

// IsPositive checks if the rational number is positive
func (r *Rational) IsPositive() bool {
	return r.Sign() > 0
}

// IsNegative checks if the rational number is negative
func (r *Rational) IsNegative() bool {
	return r.Sign() < 0
}

// ToFloat64 converts the rational number to a float64
func (r *Rational) ToFloat64() float64 {
	return float64(r.num.Int64()) / float64(r.den.Int64())
}

// ToInt64 converts the rational number to an int64, truncating fractional part
func (r *Rational) ToInt64() int64 {
	result := new(big.Int)
	result.Div(r.num, r.den)
	return result.Int64()
}

// String returns a string representation of the rational number
func (r *Rational) String() string {
	if r.den.Cmp(big.NewInt(1)) == 0 {
		return r.num.String()
	}
	return fmt.Sprintf("%s/%s", r.num.String(), r.den.String())
}

// Reduce reduces the rational number to its lowest terms
func (r *Rational) Reduce() *Rational {
	gcd := new(big.Int).GCD(nil, nil, r.num, r.den)

	if gcd.Sign() > 0 {
		num := new(big.Int).Div(r.num, gcd)
		den := new(big.Int).Div(r.den, gcd)

		return &Rational{num: num, den: den}
	}

	return r
}

// Numerator returns the numerator of the rational number
func (r *Rational) Numerator() *big.Int {
	return new(big.Int).Set(r.num)
}

// Denominator returns the denominator of the rational number
func (r *Rational) Denominator() *big.Int {
	return new(big.Int).Set(r.den)
}

// Equals checks if two rational numbers are equal
func (r *Rational) Equals(other *Rational) bool {
	return r.Compare(other) == 0
}

// Clone creates a copy of the rational number
func (r *Rational) Clone() *Rational {
	return &Rational{
		num: new(big.Int).Set(r.num),
		den: new(big.Int).Set(r.den),
	}
}

// Power raises the rational number to the given integer power
func (r *Rational) Power(exp int64) *Rational {
	if exp == 0 {
		return NewFromInt(1)
	}

	if exp > 0 {
		num := new(big.Int).Exp(r.num, big.NewInt(exp), nil)
		den := new(big.Int).Exp(r.den, big.NewInt(exp), nil)
		return NewFromBigInt(num, den)
	} else { // exp < 0
		// r^(-n) = (1/r)^n
		inverted := r.Invert()
		positiveExp := -exp
		num := new(big.Int).Exp(inverted.num, big.NewInt(positiveExp), nil)
		den := new(big.Int).Exp(inverted.den, big.NewInt(positiveExp), nil)
		return NewFromBigInt(num, den)
	}
}
