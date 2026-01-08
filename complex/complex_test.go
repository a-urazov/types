package complex

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	c := New(3, 4)
	if c.Real != 3 {
		t.Errorf("Expected real part to be 3, got %f", c.Real)
	}
	if c.Imag != 4 {
		t.Errorf("Expected imaginary part to be 4, got %f", c.Imag)
	}
}

func TestAdd(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(1, 2)
	result := c1.Add(c2)

	expected := New(4, 6)
	tolerance := 1e-9
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSubtract(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(1, 2)
	result := c1.Subtract(c2)

	expected := New(2, 2)
	tolerance := 1e-9
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMultiply(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(1, 2)
	result := c1.Multiply(c2)

	// (3+4i) * (1+2i) = 3*1 - 4*2 + (3*2 + 4*1)i = -5 + 10i
	expected := New(-5, 10)
	tolerance := 1e-9
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDivide(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(1, 2)
	result := c1.Divide(c2)

	// (3+4i) / (1+2i) = ((3*1 + 4*2) + (4*1 - 3*2)i) / (1² + 2²) = (11 - 2i) / 5 = 2.2 - 0.4i
	expected := New(2.2, -0.4)
	tolerance := 1e-9
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDivideByZero(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(0, 0)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected division by zero to panic")
		}
	}()

	c1.Divide(c2)
}

func TestMagnitude(t *testing.T) {
	c := New(3, 4)
	mag := c.Magnitude()

	expected := 5.0 // sqrt(3² + 4²) = 5
	tolerance := 1e-9
	if math.Abs(mag-expected) > tolerance {
		t.Errorf("Expected magnitude to be %f, got %f", expected, mag)
	}
}

func TestConjugate(t *testing.T) {
	c := New(3, 4)
	conj := c.Conjugate()

	expected := New(3, -4)
	tolerance := 1e-9
	if !conj.Equals(expected, tolerance) {
		t.Errorf("Expected conjugate to be %v, got %v", expected, conj)
	}
}

func TestPolar(t *testing.T) {
	c := New(3, 4)
	r, theta := c.Polar()

	expectedR := 5.0 // sqrt(3² + 4²) = 5
	expectedTheta := math.Atan2(4, 3)

	tolerance := 1e-9
	if math.Abs(r-expectedR) > tolerance {
		t.Errorf("Expected magnitude to be %f, got %f", expectedR, r)
	}
	if math.Abs(theta-expectedTheta) > tolerance {
		t.Errorf("Expected angle to be %f, got %f", expectedTheta, theta)
	}
}

func TestFromPolar(t *testing.T) {
	r := 5.0
	theta := math.Pi / 4 // 45 degrees
	c := FromPolar(r, theta)

	expectedReal := r * math.Cos(theta)
	expectedImag := r * math.Sin(theta)

	tolerance := 1e-9
	if math.Abs(c.Real-expectedReal) > tolerance {
		t.Errorf("Expected real part to be %f, got %f", expectedReal, c.Real)
	}
	if math.Abs(c.Imag-expectedImag) > tolerance {
		t.Errorf("Expected imaginary part to be %f, got %f", expectedImag, c.Imag)
	}
}

func TestPower(t *testing.T) {
	c := New(1, 1) // 1+i

	// (1+i)^2 = 1 + 2i + i^2 = 1 + 2i - 1 = 2i
	result := c.Power(2)
	expected := New(0, 2)
	tolerance := 1e-9
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Any number to the power of 0 should be 1+0i
	result = c.Power(0)
	expected = New(1, 0)
	if !result.Equals(expected, tolerance) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Any number to the power of 1 should be itself
	result = c.Power(1)
	if !result.Equals(c, tolerance) {
		t.Errorf("Expected %v, got %v", c, result)
	}
}

func TestSqrt(t *testing.T) {
	c := New(0, 1) // i
	result := c.Sqrt()

	// sqrt(i) = (1+i)/sqrt(2)
	expectedReal := 1.0 / math.Sqrt(2)
	expectedImag := 1.0 / math.Sqrt(2)

	tolerance := 1e-9
	if math.Abs(result.Real-expectedReal) > tolerance {
		t.Errorf("Expected real part to be %f, got %f", expectedReal, result.Real)
	}
	if math.Abs(result.Imag-expectedImag) > tolerance {
		t.Errorf("Expected imaginary part to be %f, got %f", expectedImag, result.Imag)
	}
}

func TestEquals(t *testing.T) {
	c1 := New(3, 4)
	c2 := New(3.0000001, 4.0000001)

	// Should be equal within tolerance
	if !c1.Equals(c2, 1e-5) {
		t.Errorf("Expected %v and %v to be equal within tolerance", c1, c2)
	}

	// Should not be equal outside tolerance
	if c1.Equals(c2, 1e-8) {
		t.Errorf("Expected %v and %v to not be equal outside tolerance", c1, c2)
	}
}

func TestString(t *testing.T) {
	c := New(3, 4)
	expected := "3.000000+4.000000i"

	if c.String() != expected {
		t.Errorf("Expected string representation to be %s, got %s", expected, c.String())
	}

	c = New(3, -4)
	expected = "3.000000-4.000000i"

	if c.String() != expected {
		t.Errorf("Expected string representation to be %s, got %s", expected, c.String())
	}
}