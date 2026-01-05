package nullable

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	value := New("test")
	if value.V != "test" {
		t.Errorf("Expected 'test', got %v", value.V)
	}
	if !value.Valid {
		t.Error("Expected Valid to be true")
	}
}

func TestIsNull(t *testing.T) {
	value := New("test")
	if value.IsNull() {
		t.Error("Expected IsNull to return false for valid value")
	}

	nullValue := New[string]()
	if !nullValue.IsNull() {
		t.Error("Expected IsNull to return true for null value")
	}
}

func TestIsNotNull(t *testing.T) {
	value := New("test")
	if !value.IsNotNull() {
		t.Error("Expected IsNotNull to return true for valid value")
	}

	nullValue := New[string]()
	if nullValue.IsNotNull() {
		t.Error("Expected IsNotNull to return false for null value")
	}
}

func TestSet(t *testing.T) {
	var nullableValue Type[string] = New[string]()
	nullableValue.Set("new value")

	if nullableValue.V != "new value" {
		t.Errorf("Expected 'new value', got %v", nullableValue.V)
	}
	if !nullableValue.Valid {
		t.Error("Expected Valid to be true after Set")
	}
}

func TestSetNull(t *testing.T) {
	value := New("test")
	value.SetNull()

	var expected string
	if value.V != expected {
		t.Errorf("Expected zero value %v, got %v", expected, value.V)
	}
	if value.Valid {
		t.Error("Expected Valid to be false after SetNull")
	}
}

func TestGet(t *testing.T) {
	value := New("test")
	got, ok := value.Get()

	if got != "test" {
		t.Errorf("Expected 'test', got %v", got)
	}
	if !ok {
		t.Error("Expected ok to be true")
	}

	nullValue := New[string]()
	got2, ok2 := nullValue.Get()

	var expected string
	if got2 != expected {
		t.Errorf("Expected zero value %v, got %v", expected, got2)
	}
	if ok2 {
		t.Error("Expected ok to be false")
	}
}

func TestOr(t *testing.T) {
	value := New("test")
	result := value.Or("default")

	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}

	nullValue := New[string]()
	result2 := nullValue.Or("default")

	if result2 != "default" {
		t.Errorf("Expected 'default', got %v", result2)
	}
}

func TestString(t *testing.T) {
	value := New("test")
	if value.String() != "test" {
		t.Errorf("Expected 'test', got %v", value.String())
	}

	nullValue := New[string]()
	if nullValue.String() != "null" {
		t.Errorf("Expected 'null', got %v", nullValue.String())
	}
}

func TestMarshalJSON(t *testing.T) {
	value := New("test")
	jsonData, err := json.Marshal(value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(jsonData) != "\"test\"" {
		t.Errorf("Expected \"test\", got %s", string(jsonData))
	}

	nullValue := New[string]()
	jsonData2, err := json.Marshal(nullValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(jsonData2) != "null" {
		t.Errorf("Expected null, got %s", string(jsonData2))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var value Type[string]
	err := json.Unmarshal([]byte(`"test"`), &value)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value.V != "test" {
		t.Errorf("Expected 'test', got %v", value.V)
	}
	if !value.Valid {
		t.Error("Expected Valid to be true after unmarshaling valid JSON")
	}

	var nullValue Type[string]
	err = json.Unmarshal([]byte("null"), &nullValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var expected string
	if nullValue.V != expected {
		t.Errorf("Expected zero value %v, got %v", expected, nullValue.V)
	}
	if nullValue.Valid {
		t.Error("Expected Valid to be false after unmarshaling null JSON")
	}
}

func TestGenericWithInt(t *testing.T) {
	intValue := New(42)
	if intValue.V != 42 {
		t.Errorf("Expected 42, got %v", intValue.V)
	}
	if !intValue.Valid {
		t.Error("Expected Valid to be true")
	}

	nullInt := New[int]()
	if nullInt.V != 0 { // zero value for int
		t.Errorf("Expected 0, got %v", nullInt.V)
	}
	if nullInt.Valid {
		t.Error("Expected Valid to be false")
	}
}

func TestGenericWithBool(t *testing.T) {
	boolValue := New(true)
	if !boolValue.V {
		t.Error("Expected true, got false")
	}
	if !boolValue.Valid {
		t.Error("Expected Valid to be true")
	}

	nullBool := New[bool]()
	if nullBool.V {
		t.Error("Expected false, got true")
	}
	if nullBool.Valid {
		t.Error("Expected Valid to be false")
	}
}
