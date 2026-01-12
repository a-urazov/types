package streams

import (
	"reflect"
	"strings"
	"testing"
)

func TestOf(t *testing.T) {
	data := []int{1, 2, 3}
	s := Of(data)
	if s == nil {
		t.Error("Of() returned nil")
	}
	collected := s.Collect()
	if !reflect.DeepEqual(collected, data) {
		t.Errorf("Collect() got = %v, want %v", collected, data)
	}
}

func TestFromReader(t *testing.T) {
	reader := strings.NewReader("a\nb\nc")
	s := FromReader(reader)
	result := s.Collect()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromReader() got = %v, want %v", result, expected)
	}
}

