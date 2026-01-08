package option

import (
	"encoding/json"
	"testing"
)

func TestSome(t *testing.T) {
	value := 42
	opt := Some(value)

	if !opt.IsSome() {
		t.Error("Expected IsSome to return true for Some value")
	}

	if opt.IsNone() {
		t.Error("Expected IsNone to return false for Some value")
	}

	retrieved, ok := opt.Get()
	if !ok {
		t.Error("Expected Get to return true for Some value")
	}

	if retrieved != value {
		t.Errorf("Expected value to be %d, got %d", value, retrieved)
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()

	if opt.IsSome() {
		t.Error("Expected IsSome to return false for None value")
	}

	if !opt.IsNone() {
		t.Error("Expected IsNone to return true for None value")
	}

	_, ok := opt.Get()
	if ok {
		t.Error("Expected Get to return false for None value")
	}
}

func TestGetOrElse(t *testing.T) {
	// Test with Some value
	someOpt := Some(42)
	result1 := someOpt.GetOrElse(0)
	if result1 != 42 {
		t.Errorf("Expected GetOrElse to return Some value 42, got %d", result1)
	}

	// Test with None value
	noneOpt := None[int]()
	result2 := noneOpt.GetOrElse(100)
	if result2 != 100 {
		t.Errorf("Expected GetOrElse to return default value 100, got %d", result2)
	}
}

func TestGetOrCall(t *testing.T) {
	// Test with Some value
	someOpt := Some(42)
	result1 := someOpt.GetOrCall(func() int { return 99 })
	if result1 != 42 {
		t.Errorf("Expected GetOrCall to return Some value 42, got %d", result1)
	}

	// Test with None value
	noneOpt := None[int]()
	callCount := 0
	result2 := noneOpt.GetOrCall(func() int {
		callCount++
		return 100
	})
	if result2 != 100 {
		t.Errorf("Expected GetOrCall to return default value 100, got %d", result2)
	}
	if callCount != 1 {
		t.Error("Expected the default function to be called once for None value")
	}
}

func TestMap(t *testing.T) {
	// Test Map with Some value
	someOpt := Some(5)
	mapped := Map(someOpt, func(x int) string { return "value: " + string(rune(x+48)) })

	if !mapped.IsSome() {
		t.Error("Expected mapped Option to be Some")
	}

	mappedValue, ok := mapped.Get()
	if !ok || mappedValue != "value: 5" {
		t.Errorf("Expected mapped value to be 'value: 5', got '%s'", mappedValue)
	}

	// Test Map with None value - should remain None
	noneOpt := None[int]()
	mappedNone := Map(noneOpt, func(x int) string { return "should not reach here" })

	if !mappedNone.IsNone() {
		t.Error("Expected mapped Option to be None")
	}
}

func TestFlatMap(t *testing.T) {
	// Test FlatMap with Some value
	someOpt := Some(5)
	flatMapped := FlatMap(someOpt, func(x int) Option[string] {
		if x > 0 {
			return Some("positive: " + string(rune(x+48)))
		}
		return None[string]()
	})

	if !flatMapped.IsSome() {
		t.Error("Expected flatMapped Option to be Some")
	}

	flatMappedValue, ok := flatMapped.Get()
	if !ok || flatMappedValue != "positive: 5" {
		t.Errorf("Expected flatMapped value to be 'positive: 5', got '%s'", flatMappedValue)
	}

	// Test FlatMap with None value - should remain None
	noneOpt := None[int]()
	flatMappedNone := FlatMap(noneOpt, func(x int) Option[string] {
		return Some("should not reach here")
	})

	if !flatMappedNone.IsNone() {
		t.Error("Expected flatMapped Option to be None")
	}
}

func TestFilter(t *testing.T) {
	// Test Filter with Some value that passes predicate
	someOpt := Some(5)
	filtered := Filter(someOpt, func(x int) bool { return x > 0 })

	if !filtered.IsSome() {
		t.Error("Expected filtered Option to be Some")
	}

	filteredValue, ok := filtered.Get()
	if !ok || filteredValue != 5 {
		t.Errorf("Expected filtered value to be 5, got %d", filteredValue)
	}

	// Test Filter with Some value that fails predicate
	someOpt2 := Some(-5)
	filtered2 := Filter(someOpt2, func(x int) bool { return x > 0 })

	if !filtered2.IsNone() {
		t.Error("Expected filtered Option to be None when predicate fails")
	}

	// Test Filter with None value - should remain None
	noneOpt := None[int]()
	filtered3 := Filter(noneOpt, func(x int) bool { return true })

	if !filtered3.IsNone() {
		t.Error("Expected filtered Option to be None when starting with None")
	}
}

func TestString(t *testing.T) {
	someOpt := Some("hello")
	expected1 := "Some(hello)"

	if someOpt.String() != expected1 {
		t.Errorf("Expected String() to return '%s', got '%s'", expected1, someOpt.String())
	}

	noneOpt := None[string]()
	expected2 := "None"

	if noneOpt.String() != expected2 {
		t.Errorf("Expected String() to return '%s', got '%s'", expected2, noneOpt.String())
	}
}

func TestJSONMarshal(t *testing.T) {
	someOpt := Some("hello")
	data1, err := json.Marshal(someOpt)
	if err != nil {
		t.Errorf("Failed to marshal Some Option: %v", err)
	}

	expected1 := `"hello"`
	if string(data1) != expected1 {
		t.Errorf("Expected marshaled Some to be '%s', got '%s'", expected1, string(data1))
	}

	noneOpt := None[string]()
	data2, err := json.Marshal(noneOpt)
	if err != nil {
		t.Errorf("Failed to marshal None Option: %v", err)
	}

	expected2 := `null`
	if string(data2) != expected2 {
		t.Errorf("Expected marshaled None to be '%s', got '%s'", expected2, string(data2))
	}
}

func TestJSONUnmarshal(t *testing.T) {
	// Test unmarshaling a value (Some)
	jsonStr := `"hello"`
	var someOpt Option[string]
	err := json.Unmarshal([]byte(jsonStr), &someOpt)
	if err != nil {
		t.Errorf("Failed to unmarshal Some Option: %v", err)
	}

	if !someOpt.IsSome() {
		t.Error("Expected unmarshaled Option to be Some")
	}

	value, ok := someOpt.Get()
	if !ok || value != "hello" {
		t.Errorf("Expected unmarshaled value to be 'hello', got '%s'", value)
	}

	// Test unmarshaling null (None)
	jsonStr2 := `null`
	var noneOpt Option[string]
	err = json.Unmarshal([]byte(jsonStr2), &noneOpt)
	if err != nil {
		t.Errorf("Failed to unmarshal None Option: %v", err)
	}

	if !noneOpt.IsNone() {
		t.Error("Expected unmarshaled Option to be None")
	}
}
