package rational

import (
	"math/big"
	"testing"
)

func TestNew(t *testing.T) {
	r := New(3, 4)
	if r.num.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("Expected numerator to be 3, got %s", r.num.String())
	}
	if r.den.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Expected denominator to be 4, got %s", r.den.String())
	}
}

func TestNewWithNegativeDenominator(t *testing.T) {
	r := New(3, -4)
	if r.num.Cmp(big.NewInt(-3)) != 0 {
		t.Errorf("Expected numerator to be -3, got %s", r.num.String())
	}
	if r.den.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Expected denominator to be 4, got %s", r.den.String())
	}
}

func TestNewFromInt(t *testing.T) {
	r := NewFromInt(5)
	if r.num.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("Expected numerator to be 5, got %s", r.num.String())
	}
	if r.den.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected denominator to be 1, got %s", r.den.String())
	}
}

func TestNewFromBigInt(t *testing.T) {
	num := big.NewInt(7)
	den := big.NewInt(8)
	r := NewFromBigInt(num, den)

	if r.num.Cmp(num) != 0 {
		t.Errorf("Expected numerator to be 7, got %s", r.num.String())
	}
	if r.den.Cmp(den) != 0 {
		t.Errorf("Expected denominator to be 8, got %s", r.den.String())
	}
}

func TestNewPanicsOnZeroDenominator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected New to panic with zero denominator")
		}
	}()

	New(1, 0)
}

func TestAdd(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(1, 3)
	result := r1.Add(r2)

	// Expected: 1/2 + 1/3 = 3/6 + 2/6 = 5/6
	expected := New(5, 6)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSubtract(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(1, 3)
	result := r1.Subtract(r2)

	// Expected: 1/2 - 1/3 = 3/6 - 2/6 = 1/6
	expected := New(1, 6)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMultiply(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(2, 3)
	result := r1.Multiply(r2)

	// Expected: 1/2 * 2/3 = 2/6 = 1/3
	expected := New(1, 3)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDivide(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(2, 3)
	result := r1.Divide(r2)

	// Expected: 1/2 ÷ 2/3 = 1/2 * 3/2 = 3/4
	expected := New(3, 4)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDivideByZero(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(0, 1)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected division by zero to panic")
		}
	}()

	r1.Divide(r2)
}

func TestInvert(t *testing.T) {
	r := New(3, 4)
	inverse := r.Invert()

	// Expected: 1/(3/4) = 4/3
	expected := New(4, 3)
	if !inverse.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, inverse)
	}
}

func TestInvertZero(t *testing.T) {
	r := New(0, 1)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected inversion of zero to panic")
		}
	}()

	r.Invert()
}

func TestNegate(t *testing.T) {
	r := New(3, 4)
	negated := r.Negate()

	// Expected: -(3/4) = -3/4
	expected := New(-3, 4)
	if !negated.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, negated)
	}
}

func TestAbs(t *testing.T) {
	r := New(-3, 4)
	abs := r.Abs()

	// Expected: |-3/4| = 3/4
	expected := New(3, 4)
	if !abs.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, abs)
	}

	r2 := New(3, -4)
	abs2 := r2.Abs()

	// Expected: |3/-4| = 3/4
	expected2 := New(3, 4)
	if !abs2.Equals(expected2) {
		t.Errorf("Expected %v, got %v", expected2, abs2)
	}
}

func TestCompare(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(2, 4) // Same as 1/2
	r3 := New(2, 3)

	if r1.Compare(r2) != 0 {
		t.Errorf("Expected %v to equal %v", r1, r2)
	}

	if r1.Compare(r3) >= 0 {
		t.Errorf("Expected %v to be less than %v", r1, r3)
	}

	if r3.Compare(r1) <= 0 {
		t.Errorf("Expected %v to be greater than %v", r3, r1)
	}
}

func TestSign(t *testing.T) {
	pos := New(1, 2)
	if pos.Sign() <= 0 {
		t.Errorf("Expected positive number to have positive sign")
	}

	neg := New(-1, 2)
	if neg.Sign() >= 0 {
		t.Errorf("Expected negative number to have negative sign")
	}

	zero := New(0, 1)
	if zero.Sign() != 0 {
		t.Errorf("Expected zero to have zero sign")
	}
}

func TestIsZero(t *testing.T) {
	zero := New(0, 1)
	if !zero.IsZero() {
		t.Errorf("Expected zero to be zero")
	}

	nonZero := New(1, 2)
	if nonZero.IsZero() {
		t.Errorf("Expected non-zero to not be zero")
	}
}

func TestIsPositive(t *testing.T) {
	pos := New(1, 2)
	if !pos.IsPositive() {
		t.Errorf("Expected positive number to be positive")
	}

	neg := New(-1, 2)
	if neg.IsPositive() {
		t.Errorf("Expected negative number to not be positive")
	}

	zero := New(0, 1)
	if zero.IsPositive() {
		t.Errorf("Expected zero to not be positive")
	}
}

func TestIsNegative(t *testing.T) {
	neg := New(-1, 2)
	if !neg.IsNegative() {
		t.Errorf("Expected negative number to be negative")
	}

	pos := New(1, 2)
	if pos.IsNegative() {
		t.Errorf("Expected positive number to not be negative")
	}

	zero := New(0, 1)
	if zero.IsNegative() {
		t.Errorf("Expected zero to not be negative")
	}
}

func TestToFloat64(t *testing.T) {
	r := New(1, 4)
	expected := 0.25
	if r.ToFloat64() != expected {
		t.Errorf("Expected %f, got %f", expected, r.ToFloat64())
	}
}

func TestToInt64(t *testing.T) {
	r := New(7, 2) // 3.5
	expected := int64(3)
	if r.ToInt64() != expected {
		t.Errorf("Expected %d, got %d", expected, r.ToInt64())
	}
}

func TestString(t *testing.T) {
	r1 := New(3, 4)
	expected1 := "3/4"
	if r1.String() != expected1 {
		t.Errorf("Expected %s, got %s", expected1, r1.String())
	}

	r2 := New(5, 1)
	expected2 := "5"
	if r2.String() != expected2 {
		t.Errorf("Expected %s, got %s", expected2, r2.String())
	}
}

func TestReduce(t *testing.T) {
	// Create 4/8 which should reduce to 1/2
	r := New(4, 8)
	reduced := r.Reduce()
	expected := New(1, 2)

	if !reduced.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, reduced)
	}
}

func TestNumeratorAndDenominator(t *testing.T) {
	r := New(3, 4)
	num := r.Numerator()
	den := r.Denominator()

	if num.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("Expected numerator to be 3, got %s", num.String())
	}
	if den.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Expected denominator to be 4, got %s", den.String())
	}
}

func TestEquals(t *testing.T) {
	r1 := New(1, 2)
	r2 := New(2, 4) // Same as 1/2 when reduced
	r3 := New(1, 3)

	if !r1.Equals(r2) {
		t.Errorf("Expected %v to equal %v", r1, r2)
	}

	if r1.Equals(r3) {
		t.Errorf("Expected %v to not equal %v", r1, r3)
	}
}

func TestClone(t *testing.T) {
	r1 := New(3, 4)
	r2 := r1.Clone()

	if !r1.Equals(r2) {
		t.Errorf("Expected clone to equal original")
	}

	// Modify clone and ensure original is unchanged
	r2 = New(1, 2)
	if r1.Equals(r2) {
		t.Errorf("Original should be unchanged after modifying clone")
	}
}

func TestPower(t *testing.T) {
	r := New(2, 3)

	// Test power of 0 -> should be 1
	result := r.Power(0)
	expected := NewFromInt(1)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test power of 1 -> should be same
	result = r.Power(1)
	expected = r
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test power of 2 -> (2/3)^2 = 4/9
	result = r.Power(2)
	expected = New(4, 9)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test negative power -> (2/3)^(-1) = 3/2
	result = r.Power(-1)
	expected = New(3, 2)
	if !result.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
