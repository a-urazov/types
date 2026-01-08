package matrix

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		wantErr  bool
	}{
		{"Valid dimensions", 3, 3, false},
		{"Zero rows", 0, 3, true},
		{"Zero cols", 3, 0, true},
		{"Negative rows", -1, 3, true},
		{"Negative cols", 3, -1, true},
		{"Both negative", -1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New[int](tt.rows, tt.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewWithValues(t *testing.T) {
	tests := []struct {
		name     string
		values   [][]int
		wantErr  bool
	}{
		{
			"Valid 2x2 matrix",
			[][]int{{1, 2}, {3, 4}},
			false,
		},
		{
			"Empty values",
			[][]int{},
			true,
		},
		{
			"Empty row",
			[][]int{{}},
			true,
		},
		{
			"Different row lengths",
			[][]int{{1, 2}, {3}},
			true,
		},
		{
			"Valid 3x3 matrix",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWithValues(tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMatrix_GetSet(t *testing.T) {
	matrix, _ := New[int](3, 3)

	// Test setting and getting values
	err := matrix.Set(1, 1, 42)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	value, err := matrix.Get(1, 1)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if value != 42 {
		t.Errorf("Get() = %v, want %v", value, 42)
	}

	// Test bounds checking
	_, err = matrix.Get(-1, 0)
	if err == nil {
		t.Error("Get() with negative row should return error")
	}

	_, err = matrix.Get(0, -1)
	if err == nil {
		t.Error("Get() with negative col should return error")
	}

	_, err = matrix.Get(5, 0)
	if err == nil {
		t.Error("Get() with out-of-bounds row should return error")
	}

	_, err = matrix.Get(0, 5)
	if err == nil {
		t.Error("Get() with out-of-bounds col should return error")
	}

	err = matrix.Set(-1, 0, 1)
	if err == nil {
		t.Error("Set() with negative row should return error")
	}

	err = matrix.Set(0, -1, 1)
	if err == nil {
		t.Error("Set() with negative col should return error")
	}

	err = matrix.Set(5, 0, 1)
	if err == nil {
		t.Error("Set() with out-of-bounds row should return error")
	}

	err = matrix.Set(0, 5, 1)
	if err == nil {
		t.Error("Set() with out-of-bounds col should return error")
	}
}

func TestMatrix_Size(t *testing.T) {
	matrix, _ := New[int](4, 5)

	rows, cols := matrix.Size()
	if rows != 4 || cols != 5 {
		t.Errorf("Size() = (%d, %d), want (4, 5)", rows, cols)
	}

	if matrix.Rows() != 4 {
		t.Errorf("Rows() = %d, want %d", matrix.Rows(), 4)
	}

	if matrix.Cols() != 5 {
		t.Errorf("Cols() = %d, want %d", matrix.Cols(), 5)
	}
}

func TestMatrix_Clone(t *testing.T) {
	original, _ := NewWithValues([][]int{{1, 2, 3}, {4, 5, 6}})
	cloned := original.Clone()

	// Check that sizes match
	if original.Rows() != cloned.Rows() || original.Cols() != cloned.Cols() {
		t.Fatal("Cloned matrix has different size than original")
	}

	// Check that values match
	for i := 0; i < original.Rows(); i++ {
		for j := 0; j < original.Cols(); j++ {
			origVal, _ := original.Get(i, j)
			clonedVal, _ := cloned.Get(i, j)
			if origVal != clonedVal {
				t.Errorf("Clone differs at [%d][%d]: original=%d, cloned=%d", i, j, origVal, clonedVal)
			}
		}
	}

	// Modify original and ensure clone is unaffected
	original.Set(0, 0, 999)
	origVal, _ := original.Get(0, 0)
	clonedVal, _ := cloned.Get(0, 0)
	if origVal == clonedVal {
		t.Error("Clone should be independent of original")
	}
}

func TestMatrix_FillReset(t *testing.T) {
	matrix, _ := New[int](3, 3)

	// Fill with 5
	matrix.Fill(5)

	// Verify all values are 5
	for i := 0; i < matrix.Rows(); i++ {
		for j := 0; j < matrix.Cols(); j++ {
			val, _ := matrix.Get(i, j)
			if val != 5 {
				t.Errorf("Fill() failed at [%d][%d]: got %d, want 5", i, j, val)
			}
		}
	}

	// Reset to zero
	matrix.Reset()

	// Verify all values are zero
	for i := 0; i < matrix.Rows(); i++ {
		for j := 0; j < matrix.Cols(); j++ {
			val, _ := matrix.Get(i, j)
			var zero int
			if val != zero {
				t.Errorf("Reset() failed at [%d][%d]: got %d, want 0", i, j, val)
			}
		}
	}
}

func TestMatrix_IsEqual(t *testing.T) {
	matrix1, _ := NewWithValues([][]int{{1, 2}, {3, 4}})
	matrix2, _ := NewWithValues([][]int{{1, 2}, {3, 4}})
	matrix3, _ := NewWithValues([][]int{{1, 2}, {3, 5}})

	if !matrix1.IsEqual(matrix2) {
		t.Error("IsEqual() should return true for identical matrices")
	}

	if matrix1.IsEqual(matrix3) {
		t.Error("IsEqual() should return false for different matrices")
	}

	// Different sizes
	matrix4, _ := NewWithValues([][]int{{1, 2}})
	if matrix1.IsEqual(matrix4) {
		t.Error("IsEqual() should return false for matrices of different sizes")
	}
}

func TestMatrix_ForEach(t *testing.T) {
	matrix, _ := NewWithValues([][]int{{1, 2, 3}, {4, 5, 6}})

	count := 0
	sum := 0

	matrix.ForEach(func(val int, i, j int) {
		count++
		sum += val
	})

	expectedCount := matrix.Rows() * matrix.Cols()
	if count != expectedCount {
		t.Errorf("ForEach() called callback %d times, expected %d", count, expectedCount)
	}

	expectedSum := 1 + 2 + 3 + 4 + 5 + 6
	if sum != expectedSum {
		t.Errorf("ForEach() sum = %d, expected %d", sum, expectedSum)
	}
}
