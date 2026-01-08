package either

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLeft(t *testing.T) {
	leftValue := "error occurred"
	e := Left[string, int](leftValue)

	if e.IsRight() {
		t.Error("Expected IsRight to return false for Left value")
	}

	if !e.IsLeft() {
		t.Error("Expected IsLeft to return true for Left value")
	}

	left, ok := e.GetLeft()
	if !ok {
		t.Error("Expected GetLeft to return true for Left value")
	}

	if left != leftValue {
		t.Errorf("Expected left value to be '%s', got '%s'", leftValue, left)
	}

	_, ok = e.GetRight()
	if ok {
		t.Error("Expected GetRight to return false for Left value")
	}
}

func TestRight(t *testing.T) {
	rightValue := 42
	e := Right[string, int](rightValue)

	if !e.IsRight() {
		t.Error("Expected IsRight to return true for Right value")
	}

	if e.IsLeft() {
		t.Error("Expected IsLeft to return false for Right value")
	}

	right, ok := e.GetRight()
	if !ok {
		t.Error("Expected GetRight to return true for Right value")
	}

	if right != rightValue {
		t.Errorf("Expected right value to be %d, got %d", rightValue, right)
	}

	_, ok = e.GetLeft()
	if ok {
		t.Error("Expected GetLeft to return false for Right value")
	}
}

func TestOrLeft(t *testing.T) {
	leftValue := "error"
	defaultValue := "default"

	e := Left[string, int](leftValue)
	result := e.OrLeft(defaultValue)

	if result != leftValue {
		t.Errorf("Expected OrLeft to return left value '%s', got '%s'", leftValue, result)
	}

	e2 := Right[string, int](42)
	result2 := e2.OrLeft(defaultValue)

	if result2 != defaultValue {
		t.Errorf("Expected OrLeft to return default value '%s', got '%s'", defaultValue, result2)
	}
}

func TestOrRight(t *testing.T) {
	rightValue := 42
	defaultValue := 0

	e := Right[string, int](rightValue)
	result := e.OrRight(defaultValue)

	if result != rightValue {
		t.Errorf("Expected OrRight to return right value %d, got %d", rightValue, result)
	}

	e2 := Left[string, int]("error")
	result2 := e2.OrRight(defaultValue)

	if result2 != defaultValue {
		t.Errorf("Expected OrRight to return default value %d, got %d", defaultValue, result2)
	}
}

func TestMapRight(t *testing.T) {
	// Test MapRight with Right value
	e1 := Right[string, int](5)
	mapped := MapRight(e1, func(x int) string { return fmt.Sprintf("value: %d", x) })

	if !mapped.IsRight() {
		t.Error("Expected mapped Either to be Right")
	}

	right, ok := mapped.GetRight()
	if !ok || right != "value: 5" {
		t.Errorf("Expected mapped Right value to be 'value: 5', got '%s'", right)
	}

	// Test MapRight with Left value - should remain Left
	e2 := Left[string, int]("error")
	mapped2 := MapRight(e2, func(x int) string { return fmt.Sprintf("value: %d", x) })

	if !mapped2.IsLeft() {
		t.Error("Expected mapped Either to be Left")
	}

	left, ok := mapped2.GetLeft()
	if !ok || left != "error" {
		t.Errorf("Expected mapped Left value to be 'error', got '%s'", left)
	}
}

func TestString(t *testing.T) {
	e1 := Left[string, int]("error")
	expected1 := "Left(error)"

	if e1.String() != expected1 {
		t.Errorf("Expected String() to return '%s', got '%s'", expected1, e1.String())
	}

	e2 := Right[string, int](42)
	expected2 := "Right(42)"

	if e2.String() != expected2 {
		t.Errorf("Expected String() to return '%s', got '%s'", expected2, e2.String())
	}
}

func TestJSONMarshal(t *testing.T) {
	e1 := Left[string, int]("error occurred")
	data1, err := json.Marshal(e1)
	if err != nil {
		t.Errorf("Failed to marshal Left Either: %v", err)
	}

	expected1 := `{"left":"error occurred"}`
	if string(data1) != expected1 {
		t.Errorf("Expected marshaled Left to be '%s', got '%s'", expected1, string(data1))
	}

	e2 := Right[string, int](123)
	data2, err := json.Marshal(e2)
	if err != nil {
		t.Errorf("Failed to marshal Right Either: %v", err)
	}

	expected2 := `{"right":123}`
	if string(data2) != expected2 {
		t.Errorf("Expected marshaled Right to be '%s', got '%s'", expected2, string(data2))
	}
}

func TestJSONUnmarshal(t *testing.T) {
	// Test unmarshaling Left
	jsonStr := `{"left":"error occurred"}`
	var e1 Either[string, int]
	err := json.Unmarshal([]byte(jsonStr), &e1)
	if err != nil {
		t.Errorf("Failed to unmarshal Left Either: %v", err)
	}

	if !e1.IsLeft() {
		t.Error("Expected unmarshaled Either to be Left")
	}

	left, ok := e1.GetLeft()
	if !ok || left != "error occurred" {
		t.Errorf("Expected unmarshaled Left value to be 'error occurred', got '%s'", left)
	}

	// Test unmarshaling Right
	jsonStr2 := `{"right":456}`
	var e2 Either[string, int]
	err = json.Unmarshal([]byte(jsonStr2), &e2)
	if err != nil {
		t.Errorf("Failed to unmarshal Right Either: %v", err)
	}

	if !e2.IsRight() {
		t.Error("Expected unmarshaled Either to be Right")
	}

	right, ok := e2.GetRight()
	if !ok || right != 456 {
		t.Errorf("Expected unmarshaled Right value to be 456, got %d", right)
	}
}
func TestFlatMapRight(t *testing.T) {
	// Test FlatMapRight with Right value
	e1 := Right[string, int](5)
	flatMapped := FlatMapRight(e1, func(x int) Either[string, string] {
		if x > 0 {
			return Right[string, string](fmt.Sprintf("positive: %d", x))
		}
		return Left[string, string]("negative or zero")
	})

	if !flatMapped.IsRight() {
		t.Error("Expected flatMapped Either to be Right")
	}

	right, ok := flatMapped.GetRight()
	if !ok || right != "positive: 5" {
		t.Errorf("Expected flatMapped Right value to be 'positive: 5', got '%s'", right)
	}

	// Test FlatMapRight with Left value - should remain Left with same value
	e2 := Left[string, int]("initial error")
	flatMapped2 := FlatMapRight(e2, func(x int) Either[string, string] {
		return Right[string, string]("should not reach here")
	})

	if !flatMapped2.IsLeft() {
		t.Error("Expected flatMapped Either to be Left")
	}

	left, ok := flatMapped2.GetLeft()
	if !ok || left != "initial error" {
		t.Errorf("Expected flatMapped Left value to be 'initial error', got '%s'", left)
	}
}
func TestMapLeft(t *testing.T) {
	// Test MapLeft with Left value
	e1 := Left[string, int]("error occurred")
	converted := MapLeft(e1, func(s string) error { return fmt.Errorf("%s", s) })

	if !converted.IsLeft() {
		t.Error("Expected converted Either to be Left")
	}

	left, ok := converted.GetLeft()
	if !ok || left.Error() != "error occurred" {
		t.Errorf("Expected converted Left value to be error 'error occurred', got '%v'", left)
	}

	// Test MapLeft with Right value - should remain Right
	e2 := Right[string, int](42)
	converted2 := MapLeft(e2, func(s string) error { return fmt.Errorf("%s", s) })

	if !converted2.IsRight() {
		t.Error("Expected converted Either to be Right")
	}

	right, ok := converted2.GetRight()
	if !ok || right != 42 {
		t.Errorf("Expected converted Right value to be 42, got %d", right)
	}
}
func TestFlatMapLeft(t *testing.T) {
        // Test FlatMapLeft with Left value
        e1 := Left[string, int]("error occurred")
        flatMapped := FlatMapLeft(e1, func(s string) Either[error, int] {
                if s == "error occurred" {
                        return Left[error, int](fmt.Errorf("wrapped: %s", s))
                }
                return Right[error, int](0)
        })

        if !flatMapped.IsLeft() {
                t.Error("Expected flatMapped Either to be Left")
        }

        left, ok := flatMapped.GetLeft()
        if !ok || left.Error() != "wrapped: error occurred" {
                t.Errorf("Expected flatMapped Left value to be 'wrapped: error occurred', got '%v'", left)
        }

        // Test FlatMapLeft with Right value - should remain Right with same value
        e2 := Right[string, int](42)
        flatMapped2 := FlatMapLeft(e2, func(s string) Either[error, int] {
                return Left[error, int](fmt.Errorf("should not reach here"))
        })

        if !flatMapped2.IsRight() {
                t.Error("Expected flatMapped Either to be Right")
        }

        right, ok := flatMapped2.GetRight()
        if !ok || right != 42 {
                t.Errorf("Expected flatMapped Right value to be 42, got %d", right)
        }
}
