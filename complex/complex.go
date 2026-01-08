package complex

import (
	"fmt"
	"math"
)

// Complex представляет комплексное число с вещественной и мнимой частями
type Complex struct {
	Real float64
	Imag float64
}

// New создает новое комплексное число с заданными вещественной и мнимой частями
func New(realPart, imagPart float64) Complex {
	return Complex{Real: realPart, Imag: imagPart}
}

// Add складывает два комплексных числа
func (c Complex) Add(other Complex) Complex {
	return Complex{
		Real: c.Real + other.Real,
		Imag: c.Imag + other.Imag,
	}
}

// Subtract вычитает одно комплексное число из другого
func (c Complex) Subtract(other Complex) Complex {
	return Complex{
		Real: c.Real - other.Real,
		Imag: c.Imag - other.Imag,
	}
}

// Multiply перемножает два комплексных числа
func (c Complex) Multiply(other Complex) Complex {
	return Complex{
		Real: c.Real*other.Real - c.Imag*other.Imag,
		Imag: c.Real*other.Imag + c.Imag*other.Real,
	}
}

// Divide делит одно комплексное число на другое
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

// Magnitude возвращает модуль комплексного числа
func (c Complex) Magnitude() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

// Conjugate возвращает комплексно-сопряженное число
func (c Complex) Conjugate() Complex {
	return Complex{
		Real: c.Real,
		Imag: -c.Imag,
	}
}

// String возвращает строковое представление комплексного числа
func (c Complex) String() string {
	if c.Imag >= 0 {
		return fmt.Sprintf("%.6f+%.6fi", c.Real, c.Imag)
	}
	return fmt.Sprintf("%.6f%.6fi", c.Real, c.Imag)
}

// Equals проверяет равенство двух комплексных чисел с учетом погрешности
func (c Complex) Equals(other Complex, tolerance float64) bool {
	return math.Abs(c.Real-other.Real) < tolerance && math.Abs(c.Imag-other.Imag) < tolerance
}

// Polar возвращает полярную форму (модуль и угол) комплексного числа
func (c Complex) Polar() (r, theta float64) {
	r = c.Magnitude()
	theta = math.Atan2(c.Imag, c.Real)
	return r, theta
}

// FromPolar создает комплексное число из полярных координат
func FromPolar(r, theta float64) Complex {
	return Complex{
		Real: r * math.Cos(theta),
		Imag: r * math.Sin(theta),
	}
}

// Power возводит комплексное число в заданную степень
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

// Sqrt возвращает квадратный корень комплексного числа
func (c Complex) Sqrt() Complex {
	r, theta := c.Polar()
	return FromPolar(math.Sqrt(r), theta/2)
}
