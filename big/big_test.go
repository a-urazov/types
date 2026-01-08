package big

import (
	"testing"
)

func TestNewGenericConstructor(t *testing.T) {
	// Test integer types
	b1 := New(42)
	if b1.typ != Integer {
		t.Errorf("Expected Integer type, got %v", b1.typ)
	}
	val, ok := b1.ToInt64()
	if val != 42 || !ok {
		t.Errorf("Expected value 42, got %v", val)
	}

	b2 := New[int64](123)
	if b2.typ != Integer {
		t.Errorf("Expected Integer type, got %v", b2.typ)
	}

	b3 := New[uint32](99)
	if b3.typ != Integer {
		t.Errorf("Expected Integer type, got %v", b3.typ)
	}

	// Test decimal types
	b4 := New(3.14)
	if b4.typ != Decimal {
		t.Errorf("Expected Decimal type, got %v", b4.typ)
	}

	b5 := New[float32](2.71)
	if b5.typ != Decimal {
		t.Errorf("Expected Decimal type, got %v", b5.typ)
	}

	// Test operations work with new generic constructor
	result := b1.Add(b2)
	if result.typ != Decimal { // Addition of integers results in decimal
		t.Errorf("Expected Decimal type after addition, got %v", result.typ)
	}

	// Test that comparison works
	if b1.Compare(New(42)) != 0 {
		t.Errorf("Expected equal comparison")
	}
	if b1.Compare(New(43)) >= 0 {
		t.Errorf("Expected b1 < 43")
	}
}

func TestToInt64(t *testing.T) {
	b := New[int64](123)
	val, ok := b.ToInt64()
	if val != 123 || !ok {
		t.Errorf("Expected (123, true), got (%v, %v)", val, ok)
	}
}

func TestToFloat64(t *testing.T) {
	b := New(3.14)
	val, ok := b.ToFloat64()
	if val != 3.14 || !ok {
		t.Errorf("Expected (3.14, true), got (%v, %v)", val, ok)
	}
}

func TestStringRepresentation(t *testing.T) {
	b1 := New(42)
	if b1.String() != "42" {
		t.Errorf("Expected '42', got '%v'", b1.String())
	}

	b2 := New(3.14)
	expected := b2.String()
	if expected != "3.14" {
		t.Errorf("Expected '3.14', got '%v'", expected)
	}
}
