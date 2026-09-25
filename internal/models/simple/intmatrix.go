/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package simple

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// IntMatrix represents a square matrix of integers
type IntMatrix [][]int

// NewIntMatrix creates a new square matrix of size n x n
func NewIntMatrix(n int) IntMatrix {
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}
	return matrix
}

// Size returns the dimension of the square matrix
func (m IntMatrix) Size() int {
	return len(m)
}

// Scan converts a database value (string) to IntMatrix
func (m *IntMatrix) Scan(value interface{}) error {
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case nil:
		*m = make(IntMatrix, 0)
		return nil
	default:
		return fmt.Errorf("failed to convert database value to bytes")
	}

	if str == "" {
		*m = make([][]int, 0)
		return nil
	}

	// Split the string into elements
	elements := strings.Split(str, "+")

	// Calculate matrix size (N for NxN matrix)
	n := int(math.Sqrt(float64(len(elements))))
	if n*n != len(elements) {
		return fmt.Errorf("invalid matrix size: element count must be a perfect square")
	}

	// Create the matrix
	matrix := NewIntMatrix(n)

	// Fill the matrix
	for i := 0; i < len(elements); i++ {
		num, err := strconv.Atoi(elements[i])
		if err != nil {
			return fmt.Errorf("failed to convert element to integer: %v", err)
		}
		row := i / n
		col := i % n
		matrix[row][col] = num
	}

	*m = matrix
	return nil
}

// Value converts IntMatrix to a database-storable string
func (m IntMatrix) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "", nil
	}

	n := len(m)
	elements := make([]string, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			elements[i*n+j] = strconv.Itoa(m[i][j])
		}
	}

	return strings.Join(elements, "+"), nil
}

// String returns a string representation of the matrix
func (m IntMatrix) String() string {
	if len(m) == 0 {
		return ""
	}

	n := len(m)
	elements := make([]string, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			elements[i*n+j] = strconv.Itoa(m[i][j])
		}
	}

	return strings.Join(elements, "+")
}

// Get returns the value at the specified position
func (m IntMatrix) Get(row, col int) (int, error) {
	if row < 0 || row >= len(m) || col < 0 || col >= len(m) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return m[row][col], nil
}

// Set sets the value at the specified position
func (m IntMatrix) Set(row, col, value int) error {
	if row < 0 || row >= len(m) || col < 0 || col >= len(m) {
		return fmt.Errorf("index out of bounds")
	}
	m[row][col] = value
	return nil
}
