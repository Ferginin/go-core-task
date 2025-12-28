package main

import (
	"reflect"
	"slices"
	"testing"
)

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []string
		slice2   []string
		expected []string
	}{
		{
			name:     "basic difference",
			slice1:   []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"},
			slice2:   []string{"banana", "date", "fig"},
			expected: []string{"apple", "cherry", "43", "lead", "gno1"},
		},
		{
			name:     "no common elements",
			slice1:   []string{"apple", "cherry"},
			slice2:   []string{"banana", "date"},
			expected: []string{"apple", "cherry"},
		},
		{
			name:     "all elements common",
			slice1:   []string{"apple", "banana"},
			slice2:   []string{"apple", "banana", "cherry"},
			expected: []string{},
		},
		{
			name:     "empty slice1",
			slice1:   []string{},
			slice2:   []string{"banana", "date"},
			expected: []string{},
		},
		{
			name:     "empty slice2",
			slice1:   []string{"apple", "banana"},
			slice2:   []string{},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "both empty",
			slice1:   []string{},
			slice2:   []string{},
			expected: []string{},
		},
		{
			name:     "duplicates in slice1",
			slice1:   []string{"apple", "apple", "banana", "cherry"},
			slice2:   []string{"banana"},
			expected: []string{"apple", "apple", "cherry"},
		},
		{
			name:     "duplicates in slice2",
			slice1:   []string{"apple", "banana", "cherry"},
			slice2:   []string{"banana", "banana"},
			expected: []string{"apple", "cherry"},
		},
		{
			name:     "single element in each",
			slice1:   []string{"apple"},
			slice2:   []string{"banana"},
			expected: []string{"apple"},
		},
		{
			name:     "empty strings",
			slice1:   []string{"", "apple", ""},
			slice2:   []string{""},
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveIntersection(tt.slice1, tt.slice2)

			if !slices.Equal(result, tt.expected) {
				t.Errorf("Difference(%v, %v) = %v; expected %v",
					tt.slice1, tt.slice2, result, tt.expected)
			}
		})
	}
}

func TestDifferenceUnique(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []string
		slice2   []string
		expected []string
	}{
		{
			name:     "basic unique difference",
			slice1:   []string{"apple", "banana", "cherry"},
			slice2:   []string{"banana"},
			expected: []string{"apple", "cherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveIntersection(tt.slice1, tt.slice2)

			// Сортируем для стабильного сравнения
			slices.Sort(result)
			slices.Sort(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DifferenceUnique(%v, %v) = %v; expected %v",
					tt.slice1, tt.slice2, result, tt.expected)
			}
		})
	}
}

func TestDifferenceEdgeCases(t *testing.T) {
	t.Run("nil slices", func(t *testing.T) {
		var slice1 []string = nil
		slice2 := []string{"banana"}

		result := RemoveIntersection(slice1, slice2)

		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})

	t.Run("slice2 is nil", func(t *testing.T) {
		slice1 := []string{"apple", "banana"}
		var slice2 []string = nil
		expected := []string{"apple", "banana"}

		result := RemoveIntersection(slice1, slice2)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
