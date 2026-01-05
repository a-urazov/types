package cast

import (
	"testing"
)

func TestConvert_StringToInt(t *testing.T) {
	_, err := To[int]("42")
	if err == nil {
		t.Errorf("ожидалась ошибка при преобразовании строки '42' в int, но ошибки не было")
	}
}

func TestConvert_IntToInt(t *testing.T) {
	result, err := To[int](42)
	if err != nil {
		t.Errorf("ошибка при преобразовании int в int: %v", err)
	}
	expected := 42
	if result != expected {
		t.Errorf("ожидалось '%d', получено '%d'", expected, result)
	}
}

func TestConvert_IntToFloat(t *testing.T) {
	result, err := To[float64](42)
	if err != nil {
		t.Errorf("ошибка при преобразовании int в float64: %v", err)
	}
	expected := float64(42)
	if result != expected {
		t.Errorf("ожидалось '%f', получено '%f'", expected, result)
	}
}

func TestConvert_StringToString(t *testing.T) {
	result, err := To[string]("hello")
	if err != nil {
		t.Errorf("ошибка при преобразовании string в string: %v", err)
	}
	expected := "hello"
	if result != expected {
		t.Errorf("ожидалось '%s', получено '%s'", expected, result)
	}
}

func TestStringConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{
			name:     "string to string",
			input:    "hello",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "int to string",
			input:    42,
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "float64 to string",
			input:    3.14,
			expected: "3.14",
			wantErr:  false,
		},
		{
			name:     "bool to string",
			input:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "nil to string",
			input:    nil,
			expected: "",
			wantErr:  false,
		},
	}

	converter := StringConverter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ожидалось '%s', получено '%s'", tt.expected, result)
			}
		})
	}
}

func TestNumberConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		wantErr  bool
	}{
		{
			name:     "float64 to float64",
			input:    3.14,
			expected: 3.14,
			wantErr:  false,
		},
		{
			name:     "int to float64",
			input:    42,
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "string to float64",
			input:    "123.45",
			expected: 123.45,
			wantErr:  false,
		},
		{
			name:     "bool true to float64",
			input:    true,
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "bool false to float64",
			input:    false,
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "nil to float64",
			input:    nil,
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "invalid string",
			input:    "not a number",
			expected: 0,
			wantErr:  true,
		},
	}

	converter := NumberConverter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ожидалось %f, получено %f", tt.expected, result)
			}
		})
	}
}

func TestBoolConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
		wantErr  bool
	}{
		{
			name:     "bool true",
			input:    true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string 'true'",
			input:    "true",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string '1'",
			input:    "1",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string 'yes'",
			input:    "yes",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string 'false'",
			input:    "false",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string '0'",
			input:    "0",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int 1",
			input:    1,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "int 0",
			input:    0,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "invalid string",
			input:    "maybe",
			expected: false,
			wantErr:  true,
		},
	}

	converter := BoolConverter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ожидалось %v, получено %v", tt.expected, result)
			}
		})
	}
}

func TestMapConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "map conversion",
			input: map[string]any{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "nil to map",
			input:   nil,
			wantErr: false,
		},
	}

	converter := MapConverter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("результат не должен быть nil")
			}
		})
	}
}

func TestSliceConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "slice conversion",
			input:   []any{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "nil to slice",
			input:   nil,
			wantErr: false,
		},
	}

	converter := SliceConverter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("результат не должен быть nil")
			}
		})
	}
}

func TestGetConverter(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		wantType any
	}{
		{
			name:     "string converter",
			typeName: "string",
			wantType: StringConverter{},
		},
		{
			name:     "number converter",
			typeName: "number",
			wantType: NumberConverter{},
		},
		{
			name:     "bool converter",
			typeName: "bool",
			wantType: BoolConverter{},
		},
		{
			name:     "map converter",
			typeName: "map",
			wantType: MapConverter{},
		},
		{
			name:     "slice converter",
			typeName: "slice",
			wantType: SliceConverter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := GetConverter(tt.typeName)
			if converter == nil {
				t.Errorf("конвертер не должен быть nil")
			}
		})
	}
}
