package main

import (
	"testing"
)

func TestSliceExample(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "only even numbers",
			input:    []int{2, 4, 6, 8},
			expected: []int{2, 4, 6, 8},
		},
		{
			name:     "only odd numbers",
			input:    []int{1, 3, 5, 7},
			expected: []int{},
		},
		{
			name:     "mixed numbers",
			input:    []int{1, 2, 3, 4, 5, 6},
			expected: []int{2, 4, 6},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "with zero",
			input:    []int{0, 1, 2},
			expected: []int{0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sliceExample(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestAddElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		addNum   int
		expected []int
	}{
		{
			name:     "add to non-empty slice",
			input:    []int{1, 2, 3},
			addNum:   4,
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "add to empty slice",
			input:    []int{},
			addNum:   1,
			expected: []int{1},
		},
		{
			name:     "add zero",
			input:    []int{5, 10},
			addNum:   0,
			expected: []int{5, 10, 0},
		},
		{
			name:     "add negative number",
			input:    []int{1, 2},
			addNum:   -5,
			expected: []int{1, 2, -5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]int, len(tt.input))
			copy(original, tt.input)

			result := addElements(tt.input, tt.addNum)

			// Проверяем, что оригинальный слайс не изменился
			for i := range tt.input {
				if tt.input[i] != original[i] {
					t.Errorf("Original slice was modified at index %d", i)
				}
			}

			// Проверяем результат
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestCopySlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
	}{
		{
			name:  "copy non-empty slice",
			input: []int{1, 2, 3, 4, 5},
		},
		{
			name:  "copy empty slice",
			input: []int{},
		},
		{
			name:  "copy slice with negatives",
			input: []int{-1, -2, -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := copySlice(tt.input)

			// Проверяем длину
			if len(result) != len(tt.input) {
				t.Errorf("Expected length %d, got %d", len(tt.input), len(result))
				return
			}

			// Проверяем содержимое
			for i := range tt.input {
				if result[i] != tt.input[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.input[i], result[i])
				}
			}

			// Проверяем, что это действительно копия (разные адреса)
			if len(tt.input) > 0 {
				result[0] = 999
				if tt.input[0] == 999 {
					t.Error("Slice was not copied, original was modified")
				}
			}
		})
	}
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		index    int
		expected []int
	}{
		{
			name:     "remove first element",
			input:    []int{1, 2, 3, 4, 5},
			index:    0,
			expected: []int{2, 3, 4, 5},
		},
		{
			name:     "remove last element",
			input:    []int{1, 2, 3, 4, 5},
			index:    4,
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "remove middle element",
			input:    []int{1, 2, 3, 4, 5},
			index:    2,
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "remove from single element slice",
			input:    []int{42},
			index:    0,
			expected: []int{},
		},
		{
			name:     "remove from two element slice",
			input:    []int{10, 20},
			index:    1,
			expected: []int{10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]int, len(tt.input))
			copy(original, tt.input)

			result := removeElement(tt.input, tt.index)

			// Проверяем, что оригинальный слайс не изменился
			if len(tt.input) != len(original) {
				t.Error("Original slice length was modified")
			}
			for i := range tt.input {
				if tt.input[i] != original[i] {
					t.Errorf("Original slice was modified at index %d", i)
				}
			}

			// Проверяем результат
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
