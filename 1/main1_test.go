package main

import (
	"testing"
)

func TestConcatVariables(t *testing.T) {
	vars := InitVariables()
	result := vars.ConcatVariables()

	if len(result) == 0 {
		t.Error("ConcatVariables returned empty string")
	}

	// Проверяем, что строка содержит ожидаемые подстроки
	expected := []string{"42", "75", "fa", "3.14", "Golang", "true"}
	for _, exp := range expected {
		if !contains(result, exp) {
			t.Errorf("Result doesn't contain expected substring: %s", exp)
		}
	}
}

func TestInsertSalt(t *testing.T) {
	runes := []rune("abcd")
	result := InsertSalt(runes, "XY")

	expected := []rune("abXYcd")
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	if string(result) != string(expected) {
		t.Errorf("Expected %s, got %s", string(expected), string(result))
	}
}

func TestHashRune(t *testing.T) {
	runes := []rune("test")
	hash := HashRune(runes)

	// SHA256 хэш всегда 64 символа в hex
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}
}

// Вспомогательная функция
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}
