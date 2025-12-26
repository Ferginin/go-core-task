package main

import (
	"reflect"
	"testing"
)

// TestNewStringIntMap проверяет создание новой карты
func TestNewStringIntMap(t *testing.T) {
	sim := NewStringIntMap()

	if sim == nil {
		t.Fatal("NewStringIntMap returned nil")
	}

	if sim.data == nil {
		t.Error("data map was not initialized")
	}

	if len(sim.data) != 0 {
		t.Errorf("Expected empty map, got size %d", len(sim.data))
	}
}

// TestAdd проверяет добавление элементов
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int
		expected int
	}{
		{
			name:     "add positive integer",
			key:      "positive",
			value:    42,
			expected: 42,
		},
		{
			name:     "add negative integer",
			key:      "negative",
			value:    -15,
			expected: -15,
		},
		{
			name:     "add zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
		{
			name:     "add large number",
			key:      "large",
			value:    1000000,
			expected: 1000000,
		},
		{
			name:     "add with empty string key",
			key:      "",
			value:    99,
			expected: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewStringIntMap()
			sim.Add(tt.key, tt.value)

			value, exists := sim.Get(tt.key)
			if !exists {
				t.Errorf("Key %q not found after Add", tt.key)
				return
			}
			if value != tt.expected {
				t.Errorf("Expected value %d, got %d", tt.expected, value)
			}
		})
	}
}

// TestAddOverwrite проверяет перезапись существующего значения
func TestAddOverwrite(t *testing.T) {
	sim := NewStringIntMap()

	sim.Add("key", 10)
	value, _ := sim.Get("key")
	if value != 10 {
		t.Errorf("Expected 10, got %d", value)
	}

	// Перезаписываем значение
	sim.Add("key", 20)
	value, _ = sim.Get("key")
	if value != 20 {
		t.Errorf("Expected 20 after overwrite, got %d", value)
	}

	// Проверяем, что размер не увеличился
	if len(sim.data) != 1 {
		t.Errorf("Expected map size 1, got %d", len(sim.data))
	}
}

// TestRemove проверяет удаление элементов
func TestRemove(t *testing.T) {
	tests := []struct {
		name        string
		initialData map[string]int
		keyToRemove string
		shouldExist bool
	}{
		{
			name:        "remove existing key",
			initialData: map[string]int{"key": 100},
			keyToRemove: "key",
			shouldExist: false,
		},
		{
			name:        "remove non-existing key",
			initialData: map[string]int{"key1": 100},
			keyToRemove: "key2",
			shouldExist: false,
		},
		{
			name:        "remove from empty map",
			initialData: map[string]int{},
			keyToRemove: "key",
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewStringIntMap()
			for k, v := range tt.initialData {
				sim.Add(k, v)
			}

			sim.Remove(tt.keyToRemove)

			_, exists := sim.Get(tt.keyToRemove)
			if exists != tt.shouldExist {
				t.Errorf("Expected exists=%v after Remove, got exists=%v", tt.shouldExist, exists)
			}
		})
	}
}

// TestGet проверяет получение значений
func TestGet(t *testing.T) {
	tests := []struct {
		name          string
		initialData   map[string]int
		keyToGet      string
		expectedVal   int
		expectedExist bool
	}{
		{
			name:          "get existing key",
			initialData:   map[string]int{"key": 50},
			keyToGet:      "key",
			expectedVal:   50,
			expectedExist: true,
		},
		{
			name:          "get non-existing key",
			initialData:   map[string]int{"key1": 50},
			keyToGet:      "key2",
			expectedVal:   0,
			expectedExist: false,
		},
		{
			name:          "get from empty map",
			initialData:   map[string]int{},
			keyToGet:      "key",
			expectedVal:   0,
			expectedExist: false,
		},
		{
			name:          "get zero value",
			initialData:   map[string]int{"zero": 0},
			keyToGet:      "zero",
			expectedVal:   0,
			expectedExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewStringIntMap()
			for k, v := range tt.initialData {
				sim.Add(k, v)
			}

			value, exists := sim.Get(tt.keyToGet)

			if exists != tt.expectedExist {
				t.Errorf("Expected exists=%v, got exists=%v", tt.expectedExist, exists)
			}

			if exists && value != tt.expectedVal {
				t.Errorf("Expected value %d, got %d", tt.expectedVal, value)
			}
		})
	}
}

// TestExists проверяет наличие ключей
func TestExists(t *testing.T) {
	tests := []struct {
		name        string
		initialData map[string]int
		keyToCheck  string
		expected    bool
	}{
		{
			name:        "key exists",
			initialData: map[string]int{"key": 100},
			keyToCheck:  "key",
			expected:    true,
		},
		{
			name:        "key does not exist",
			initialData: map[string]int{"key1": 100},
			keyToCheck:  "key2",
			expected:    false,
		},
		{
			name:        "check in empty map",
			initialData: map[string]int{},
			keyToCheck:  "key",
			expected:    false,
		},
		{
			name:        "key with zero value exists",
			initialData: map[string]int{"zero": 0},
			keyToCheck:  "zero",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewStringIntMap()
			for k, v := range tt.initialData {
				sim.Add(k, v)
			}

			result := sim.Exists(tt.keyToCheck)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestCopy проверяет копирование карты
func TestCopy(t *testing.T) {
	sim := NewStringIntMap()
	sim.Add("first", 11)
	sim.Add("second", 22)
	sim.Add("third", 33)

	copied := sim.Copy()

	// Проверяем размер
	if len(copied) != 3 {
		t.Errorf("Expected copied map size 3, got %d", len(copied))
	}

	// Проверяем содержимое
	expectedData := map[string]int{
		"first":  11,
		"second": 22,
		"third":  33,
	}

	if !reflect.DeepEqual(copied, expectedData) {
		t.Errorf("Copied map content mismatch.\nExpected: %v\nGot: %v", expectedData, copied)
	}
}

// TestCopyIsIndependent проверяет независимость копии
func TestCopyIsIndependent(t *testing.T) {
	sim := NewStringIntMap()
	sim.Add("first", 11)
	sim.Add("second", 22)

	copied := sim.Copy()

	// Изменяем оригинальную карту
	sim.Add("third", 33)
	sim.Add("first", 100)
	sim.Remove("second")

	// Проверяем, что копия не изменилась
	if len(copied) != 2 {
		t.Errorf("Copy was affected: expected size 2, got %d", len(copied))
	}

	if value, exists := copied["first"]; !exists || value != 11 {
		t.Errorf("Copy was affected: 'first' expected 11, got %d (exists: %v)", value, exists)
	}

	if value, exists := copied["second"]; !exists || value != 22 {
		t.Errorf("Copy was affected: 'second' expected 22, got %d (exists: %v)", value, exists)
	}

	if _, exists := copied["third"]; exists {
		t.Error("Copy was affected: 'third' should not exist in copy")
	}
}

// TestCopyEmpty проверяет копирование пустой карты
func TestCopyEmpty(t *testing.T) {
	sim := NewStringIntMap()
	copied := sim.Copy()

	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got size %d", len(copied))
	}

	if copied == nil {
		t.Error("Copy returned nil instead of empty map")
	}
}

// TestCopyModifyOriginal проверяет, что изменение копии не влияет на оригинал
func TestCopyModifyOriginal(t *testing.T) {
	sim := NewStringIntMap()
	sim.Add("key", 100)

	copied := sim.Copy()

	// Изменяем копию
	copied["key"] = 200
	copied["newkey"] = 300

	// Проверяем, что оригинал не изменился
	if value, _ := sim.Get("key"); value != 100 {
		t.Errorf("Original was affected: expected 100, got %d", value)
	}

	if sim.Exists("newkey") {
		t.Error("Original was affected: 'newkey' should not exist")
	}
}
