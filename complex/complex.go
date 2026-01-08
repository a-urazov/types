package complex

import (
	"fmt"
	"math"
)

// Complex represents a complex number with real and imaginary parts
type Complex struct {
	Real float64
	Imag float64
}

// New creates a new complex number with the given real and imaginary parts
func New(real, imag float64) Complex {
	return Complex{Real: real, Imag: imag}
}

// Add adds two complex numbers
func (c Complex) Add(other Complex) Complex {
	return Complex{
		Real: c.Real + other.Real,
		Imag: c.Imag + other.Imag,
	}
}

// Subtract subtracts the other complex number from this one
func (c Complex) Subtract(other Complex) Complex {
	return Complex{
		Real: c.Real - other.Real,
		Imag: c.Imag - other.Imag,
	}
}

// Multiply multiplies two complex numbers
func (c Complex) Multiply(other Complex) Complex {
	return Complex{
		Real: c.Real*other.Real - c.Imag*other.Imag,
		Imag: c.Real*other.Imag + c.Imag*other.Real,
	}
}

// Divide divides this complex number by the other
func (c Complex) Divide(other Complex) Complex {
	denominator := other.Real*other.Real + other.Imag*other.Imag
	if denominator == 0 {
		panic("division by zero")
	}

	return Complex{
		Real: (c.Real*other.Real + c.Imag*other.Imag) / denominator,
		Imag: (c.Imag*other.Real - c.Real*other.Imag) / denominator,
	}
}

// Magnitude returns the magnitude (or modulus) of the complex number
func (c Complex) Magnitude() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

// Conjugate returns the complex conjugate
func (c Complex) Conjugate() Complex {
	return Complex{
		Real: c.Real,
		Imag: -c.Imag,
	}
}

// String returns a string representation of the complex number
func (c Complex) String() string {
	if c.Imag >= 0 {
		return fmt.Sprintf("%.6f+%.6fi", c.Real, c.Imag)
	}
	return fmt.Sprintf("%.6f%.6fi", c.Real, c.Imag)
}

// Equals checks if two complex numbers are equal within a tolerance
func (c Complex) Equals(other Complex, tolerance float64) bool {
	return math.Abs(c.Real-other.Real) < tolerance && math.Abs(c.Imag-other.Imag) < tolerance
}

// Polar returns the polar form (magnitude and angle) of the complex number
func (c Complex) Polar() (r, theta float64) {
	r = c.Magnitude()
	theta = math.Atan2(c.Imag, c.Real)
	return r, theta
}

// FromPolar creates a complex number from polar coordinates
func FromPolar(r, theta float64) Complex {
	return Complex{
		Real: r * math.Cos(theta),
		Imag: r * math.Sin(theta),
	}
}

// Power raises the complex number to the given power
func (c Complex) Power(n int) Complex {
	if n == 0 {
		return New(1, 0)
	}
	if n == 1 {
		return c
	}

	result := c
	for i := 2; i <= n; i++ {
		result = result.Multiply(c)
	}
	return result
}

// Sqrt returns the square root of the complex number
func (c Complex) Sqrt() Complex {
	r, theta := c.Polar()
	return FromPolar(math.Sqrt(r), theta/2)
}