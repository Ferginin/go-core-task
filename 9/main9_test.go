package main

import (
	"math"
	"slices"
	"testing"
	"time"
)

// вспомогательная функция: чтение всех значений из канала с таймаутом
func collectWithTimeout[T any](t *testing.T, ch <-chan T, timeout time.Duration) []T {
	t.Helper()

	var res []T
	done := make(chan struct{})

	go func() {
		for v := range ch {
			res = append(res, v)
		}
		close(done)
	}()

	select {
	case <-done:
		return res
	case <-time.After(timeout):
		t.Fatalf("timeout while collecting from channel")
		return nil
	}
}

func TestMakeNumsConveyor_Basic(t *testing.T) {
	in := make(chan uint8)
	out := MakeNumsConveyor(in)

	go func() {
		defer close(in)
		nums := []uint8{1, 2, 3, 4}
		for _, n := range nums {
			in <- n
		}
	}()

	got := collectWithTimeout(t, out, time.Second)

	expected := []float64{
		math.Pow(1, 3),
		math.Pow(2, 3),
		math.Pow(3, 3),
		math.Pow(4, 3),
	}

	if !slices.Equal(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestMakeNumsConveyor_EmptyInput(t *testing.T) {
	in := make(chan uint8)
	out := MakeNumsConveyor(in)

	// сразу закрываем вход, ничего не отправляя
	close(in)

	got := collectWithTimeout(t, out, time.Second)
	if len(got) != 0 {
		t.Errorf("expected empty output slice, got %v", got)
	}
}

func TestMakeNumsConveyor_SingleValue(t *testing.T) {
	in := make(chan uint8)
	out := MakeNumsConveyor(in)

	go func() {
		defer close(in)
		in <- 5
	}()

	got := collectWithTimeout(t, out, time.Second)
	if len(got) != 1 {
		t.Fatalf("expected 1 value, got %d: %v", len(got), got)
	}

	expected := math.Pow(5, 3)
	if got[0] != expected {
		t.Errorf("expected %v, got %v", expected, got[0])
	}
}

func TestMakeNumsConveyor_MaxUint8(t *testing.T) {
	in := make(chan uint8)
	out := MakeNumsConveyor(in)

	go func() {
		defer close(in)
		in <- math.MaxUint8 // 255
	}()

	got := collectWithTimeout(t, out, time.Second)
	if len(got) != 1 {
		t.Fatalf("expected 1 value, got %d: %v", len(got), got)
	}

	expected := math.Pow(float64(math.MaxUint8), 3)
	if got[0] != expected {
		t.Errorf("expected %v, got %v", expected, got[0])
	}
}

func TestMakeNumsConveyor_OrderPreserved(t *testing.T) {
	in := make(chan uint8)
	out := MakeNumsConveyor(in)

	input := []uint8{0, 1, 2, 10, 3}
	go func() {
		defer close(in)
		for _, n := range input {
			in <- n
		}
	}()

	got := collectWithTimeout(t, out, time.Second)

	if len(got) != len(input) {
		t.Fatalf("expected %d values, got %d (%v)", len(input), len(got), got)
	}

	for i, n := range input {
		expected := math.Pow(float64(n), 3)
		if got[i] != expected {
			t.Errorf("at index %d: expected %v, got %v", i, expected, got[i])
		}
	}
}
