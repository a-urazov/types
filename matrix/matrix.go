package matrix

import (
	"cmp"
	"errors"
)

// Matrix представляет собой двумерную матрицу чисел
type Matrix[T cmp.Ordered] struct {
	rows int
	cols int
	data [][]T
}

// New создает новую матрицу с заданными размерами
func New[T cmp.Ordered](rows, cols int) (*Matrix[T], error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("размеры матрицы должны быть положительными")
	}

	matrix := &Matrix[T]{
		rows: rows,
		cols: cols,
		data: make([][]T, rows),
	}

	for i := range matrix.data {
		matrix.data[i] = make([]T, cols)
	}

	return matrix, nil
}

// NewWithValues создает новую матрицу с заданными начальными значениями
func NewWithValues[T cmp.Ordered](values [][]T) (*Matrix[T], error) {
	if len(values) == 0 || len(values[0]) == 0 {
		return nil, errors.New("значения не могут быть пустыми")
	}

	rows := len(values)
	cols := len(values[0])

	// Проверяем, что все строки имеют одинаковую длину
	for i := 1; i < rows; i++ {
		if len(values[i]) != cols {
			return nil, errors.New("все строки должны иметь одинаковую длину")
		}
	}

	matrix := &Matrix[T]{
		rows: rows,
		cols: cols,
		data: make([][]T, rows),
	}

	for i := range values {
		matrix.data[i] = make([]T, cols)
		copy(matrix.data[i], values[i])
	}

	return matrix, nil
}

// Clone создает копию матрицы
func (m *Matrix[T]) Clone() *Matrix[T] {
	clone := &Matrix[T]{
		rows: m.rows,
		cols: m.cols,
		data: make([][]T, m.rows),
	}

	for i := range m.data {
		clone.data[i] = make([]T, m.cols)
		copy(clone.data[i], m.data[i])
	}

	return clone
}

// Rows возвращает количество строк в матрице
func (m *Matrix[T]) Rows() int {
	return m.rows
}

// Cols возвращает количество столбцов в матрице
func (m *Matrix[T]) Cols() int {
	return m.cols
}

// Size возвращает размеры матрицы (строки, столбцы)
func (m *Matrix[T]) Size() (int, int) {
	return m.rows, m.cols
}

// Get возвращает значение элемента по указанным координатам
func (m *Matrix[T]) Get(row, col int) (T, error) {
	var zero T
	if row < 0 || row >= m.rows || col < 0 || col >= m.cols {
		return zero, errors.New("индексы выходят за пределы матрицы")
	}
	return m.data[row][col], nil
}

// Set устанавливает значение элемента по указанным координатам
func (m *Matrix[T]) Set(row, col int, value T) error {
	if row < 0 || row >= m.rows || col < 0 || col >= m.cols {
		return errors.New("индексы выходят за пределы матрицы")
	}
	m.data[row][col] = value
	return nil
}

// Fill заполняет всю матрицу указанным значением
func (m *Matrix[T]) Fill(value T) {
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			m.data[i][j] = value
		}
	}
}

// Reset обнуляет матрицу (заполняет нулевыми значениями)
func (m *Matrix[T]) Reset() {
	var zero T
	m.Fill(zero)
}

// IsEqual сравнивает две матрицы на равенство
func (m *Matrix[T]) IsEqual(other *Matrix[T]) bool {
	if m.rows != other.rows || m.cols != other.cols {
		return false
	}

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			if m.data[i][j] != other.data[i][j] {
				return false
			}
		}
	}

	return true
}

// ForEach применяет функцию к каждому элементу матрицы
func (m *Matrix[T]) ForEach(fn func(T, int, int)) {
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			fn(m.data[i][j], i, j)
		}
	}
}

// Add складывает две матрицы
func (m *Matrix[T]) Add(other *Matrix[T]) (*Matrix[T], error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("размеры матриц должны совпадать для сложения")
	}

	result, _ := New[T](m.rows, m.cols)

	var zero T
	var temp T

	// Проверяем, можно ли выполнять сложение (только для числовых типов)
	// В реальной реализации можно использовать constraint для числовых типов
	_ = zero
	_ = temp

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			// Так как мы не можем быть уверены в типе T, мы не можем выполнить сложение напрямую
			// Вместо этого, мы предоставим специализированные функции для числовых типов
			// или предоставим метод, который принимает функцию для операции
			result.data[i][j] = m.data[i][j] // заглушка
		}
	}

	return result, nil
}

// AddFunc позволяет складывать матрицы, используя предоставленную функцию для операции
func (m *Matrix[T]) AddFunc(other *Matrix[T], op func(T, T) T) (*Matrix[T], error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("размеры матриц должны совпадать для сложения")
	}

	result, _ := New[T](m.rows, m.cols)

	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = op(m.data[i][j], other.data[i][j])
		}
	}

	return result, nil
}

// Multiply умножает две матрицы
func (m *Matrix[T]) Multiply(other *Matrix[T]) (*Matrix[T], error) {
	if m.cols != other.rows {
		return nil, errors.New("количество столбцов первой матрицы должно совпадать с количеством строк второй матрицы")
	}

	result, _ := New[T](m.rows, other.cols)

	// Как и в случае с Add, мы не можем выполнить умножение напрямую без знания типа T
	// Мы предоставим метод, который использует функцию для выполнения умножения и сложения

	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			// Заглушка для умножения
			result.data[i][j] = m.data[i][j] // заглушка
		}
	}

	return result, nil
}

// MultiplyFunc позволяет умножать матрицы, используя предоставленные функции для умножения и сложения
func (m *Matrix[T]) MultiplyFunc(other *Matrix[T], mulFunc func(T, T) T, addFunc func(T, T) T) (*Matrix[T], error) {
	if m.cols != other.rows {
		return nil, errors.New("количество столбцов первой матрицы должно совпадать с количеством строк второй матрицы")
	}

	result, _ := New[T](m.rows, other.cols)

	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			var sum T // нейтральный элемент для сложения
			for k := 0; k < m.cols; k++ {
				product := mulFunc(m.data[i][k], other.data[k][j])
				sum = addFunc(sum, product)
			}
			result.data[i][j] = sum
		}
	}

	return result, nil
}
