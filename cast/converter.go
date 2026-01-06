package cast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Converter определяет интерфейс для преобразования типов
type Converter interface {
	// Convert преобразует исходное значение в целевой тип
	Convert(source any) (any, error)
}

// ConvertFunc - функциональный тип для преобразования
type ConvertFunc func(source any) (any, error)

// Convert реализует интерфейс Converter для функции
func (cf ConvertFunc) Convert(source any) (any, error) {
	return cf(source)
}

// To преобразует значение в целевой тип To
// Использует встроенные преобразования или JSON маршалинг/анмаршалинг
// Тип источника определяется динамически в runtime
func To[T any](value any) (T, error) {
	var zero T

	// Если значение nil, возвращаем нулевое значение
	if value == nil {
		return zero, nil
	}

	// Получаем тип источника и целевого типа
	fromType := reflect.TypeOf(value)
	toType := reflect.TypeOf(zero)

	// Если типы одинаковые, просто конвертируем
	if fromType == toType {
		return value.(T), nil
	}

	// Пытаемся преобразовать через reflect, если возможно
	if fromType != nil && fromType.ConvertibleTo(toType) {
		fromVal := reflect.ValueOf(value)
		return fromVal.Convert(toType).Interface().(T), nil
	}

	// Пытаемся преобразовать через JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return zero, fmt.Errorf("ошибка маршалинга в JSON: %w", err)
	}

	var result T
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return zero, fmt.Errorf("ошибка анмаршалинга из JSON: %w", err)
	}

	return result, nil
}

// ConvertTo преобразует значение к целевому типу с использованием reflect
func ConvertTo(source any, targetType reflect.Type) (any, error) {
	if source == nil {
		return reflect.Zero(targetType).Interface(), nil
	}

	sourceValue := reflect.ValueOf(source)
	sourceType := sourceValue.Type()

	// Если типы совпадают
	if sourceType == targetType {
		return source, nil
	}

	// Попытка прямого преобразования через reflect
	if sourceType.ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType).Interface(), nil
	}

	// Преобразование через JSON
	jsonData, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга в JSON: %w", err)
	}

	result := reflect.New(targetType)
	err = json.Unmarshal(jsonData, result.Interface())
	if err != nil {
		return nil, fmt.Errorf("ошибка анмаршалинга из JSON: %w", err)
	}

	return result.Elem().Interface(), nil
}

// StringConverter преобразует любое значение в строку
type StringConverter struct{}

func (sc StringConverter) Convert(source any) (any, error) {
	if source == nil {
		return "", nil
	}

	switch v := source.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v), nil
	case float32, float64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(jsonData), nil
	}
}

// NumberConverter преобразует значения в числовой формат (float64)
type NumberConverter struct{}

func (nc NumberConverter) Convert(source any) (any, error) {
	if source == nil {
		return float64(0), nil
	}

	switch v := source.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		num, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return nil, fmt.Errorf("невозможно преобразовать строку в число: %w", err)
		}
		return num, nil
	case bool:
		if v {
			return float64(1), nil
		}
		return float64(0), nil
	default:
		return nil, fmt.Errorf("невозможно преобразовать тип %T в число", source)
	}
}

// BoolConverter преобразует значения в булев формат
type BoolConverter struct{}

func (bc BoolConverter) Convert(source any) (any, error) {
	if source == nil {
		return false, nil
	}

	switch v := source.(type) {
	case bool:
		return v, nil
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		switch s {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off", "":
			return false, nil
		default:
			return nil, fmt.Errorf("невозможно преобразовать строку '%s' в bool", v)
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v) != "0", nil
	case float32, float64:
		return fmt.Sprintf("%v", v) != "0", nil
	default:
		return nil, fmt.Errorf("невозможно преобразовать тип %T в bool", source)
	}
}

// MapConverter преобразует значения в map
type MapConverter struct{}

func (mc MapConverter) Convert(source any) (any, error) {
	if source == nil {
		return map[string]any{}, nil
	}

	switch v := source.(type) {
	case map[string]any:
		return v, nil
	default:
		// Пытаемся преобразовать через JSON
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга в JSON: %w", err)
		}

		var result map[string]any
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования в map: %w", err)
		}
		return result, nil
	}
}

// SliceConverter преобразует значения в slice
type SliceConverter struct{}

func (sc SliceConverter) Convert(source any) (any, error) {
	if source == nil {
		return []any{}, nil
	}

	switch v := source.(type) {
	case []any:
		return v, nil
	default:
		sourceValue := reflect.ValueOf(v)
		if sourceValue.Kind() != reflect.Slice && sourceValue.Kind() != reflect.Array {
			// Пытаемся преобразовать через JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("ошибка маршалинга в JSON: %w", err)
			}

			var result []any
			err = json.Unmarshal(jsonData, &result)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования в slice: %w", err)
			}
			return result, nil
		}

		// Конвертируем slice/array
		result := make([]any, sourceValue.Len())
		for i := 0; i < sourceValue.Len(); i++ {
			result[i] = sourceValue.Index(i).Interface()
		}
		return result, nil
	}
}

// GetConverter возвращает подходящий конвертер для целевого типа
func GetConverter(targetType string) Converter {
	switch targetType {
	case "string":
		return StringConverter{}
	case "number", "float", "float64", "int":
		return NumberConverter{}
	case "bool", "boolean":
		return BoolConverter{}
	case "map":
		return MapConverter{}
	case "slice", "array":
		return SliceConverter{}
	default:
		return ConvertFunc(func(source any) (any, error) {
			return nil, fmt.Errorf("неизвестный целевой тип: %s", targetType)
		})
	}
}
