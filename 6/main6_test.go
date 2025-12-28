package main

import (
	"testing"
	"time"
)

// helper-функция для чтения с таймаутом
func recvWithTimeout(t *testing.T, ch <-chan int, timeout time.Duration) (int, bool) {
	t.Helper()

	select {
	case n, ok := <-ch:
		return n, ok
	case <-time.After(timeout):
		t.Fatal("timeout waiting for value from channel")
		return 0, false
	}
}

// TestNewRandNumGenerator_BasicCreation проверяет создание генератора
func TestNewRandNumGenerator_BasicCreation(t *testing.T) {
	gen := NewRandNumGenerator(0, 10)
	defer gen.Stop()

	if gen == nil {
		t.Fatal("NewRandNumGenerator returned nil")
	}

	if gen.numChan == nil {
		t.Error("numChan is nil")
	}

	if gen.stop == nil {
		t.Error("stop channel is nil")
	}
}

// TestNewRandNumGenerator_GeneratesInRange проверяет, что числа в заданном диапазоне
func TestNewRandNumGenerator_GeneratesInRange(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{
			name: "positive range",
			min:  0,
			max:  10,
		},
		{
			name: "negative range",
			min:  -10,
			max:  -1,
		},
		{
			name: "mixed range",
			min:  -5,
			max:  5,
		},
		{
			name: "large range",
			min:  0,
			max:  1000,
		},
		{
			name: "single value",
			min:  42,
			max:  42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewRandNumGenerator(tt.min, tt.max)
			defer gen.Stop()

			n, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
			if !ok {
				t.Fatal("channel closed unexpectedly")
			}

			if n < tt.min || n > tt.max {
				t.Errorf("expected number in [%d,%d], got %d", tt.min, tt.max, n)
			}
		})
	}
}

// TestRandNumGenerator_MultipleValues проверяет генерацию нескольких чисел подряд
func TestRandNumGenerator_MultipleValues(t *testing.T) {
	min, max := 1, 100
	gen := NewRandNumGenerator(min, max)
	defer gen.Stop()

	const count = 20
	for i := 0; i < count; i++ {
		n, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
		if !ok {
			t.Fatalf("channel closed unexpectedly on iteration %d", i)
		}

		if n < min || n > max {
			t.Errorf("iteration %d: expected number in [%d,%d], got %d", i, min, max, n)
		}
	}
}

// TestRandNumGenerator_Stop проверяет корректную остановку генератора
func TestRandNumGenerator_Stop(t *testing.T) {
	gen := NewRandNumGenerator(0, 10)

	// читаем одно значение, чтобы убедиться что генератор работает
	_, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
	if !ok {
		t.Fatal("channel closed before Stop was called")
	}

	// останавливаем генератор
	gen.Stop()

	// даём горутине время на завершение
	time.Sleep(50 * time.Millisecond)

	// канал должен быть закрыт
	n, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
	if ok {
		t.Errorf("expected closed channel after Stop, got value %d", n)
	}
}

// TestRandNumGenerator_StopIdempotent проверяет, что повторный вызов Stop безопасен
func TestRandNumGenerator_StopIdempotent(t *testing.T) {
	gen := NewRandNumGenerator(0, 10)

	// первый вызов Stop
	gen.Stop()

	// второй вызов не должен паниковать
	gen.Stop()

	// для уверенности
	gen.Stop()

	// канал должен быть закрыт
	_, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
	if ok {
		t.Error("expected closed channel after multiple Stop calls")
	}
}

// TestRandNumGenerator_ConcurrentReads проверяет конкурентное чтение из канала
func TestRandNumGenerator_ConcurrentReads(t *testing.T) {
	gen := NewRandNumGenerator(0, 100)
	defer gen.Stop()

	const goroutines = 5
	const readsPerGoroutine = 10

	done := make(chan bool, goroutines)
	errors := make(chan error, goroutines*readsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			for j := 0; j < readsPerGoroutine; j++ {
				select {
				case n, ok := <-gen.Next():
					if !ok {
						errors <- nil
						return
					}
					if n < 0 || n > 100 {
						t.Errorf("goroutine %d: value out of range: %d", id, n)
					}
				case <-time.After(200 * time.Millisecond):
					t.Errorf("goroutine %d: timeout", id)
					return
				}
			}
		}(i)
	}

	// ждём завершения всех горутин
	for i := 0; i < goroutines; i++ {
		<-done
	}

	close(errors)
}

// TestRandNumGenerator_StopDuringRead проверяет остановку во время чтения
func TestRandNumGenerator_StopDuringRead(t *testing.T) {
	gen := NewRandNumGenerator(0, 10)

	// читаем несколько значений
	for i := 0; i < 3; i++ {
		_, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
		if !ok {
			t.Fatal("channel closed prematurely")
		}
	}

	// останавливаем во время активной работы
	gen.Stop()

	// следующее чтение должно вернуть ok=false
	select {
	case _, ok := <-gen.Next():
		if ok {
			t.Error("expected channel to be closed after Stop")
		}
	case <-time.After(100 * time.Millisecond):
		// канал может быть уже закрыт, это нормально
	}
}

// TestRandNumGenerator_ZeroRange проверяет диапазон из одного числа
func TestRandNumGenerator_ZeroRange(t *testing.T) {
	gen := NewRandNumGenerator(5, 5)
	defer gen.Stop()

	for i := 0; i < 5; i++ {
		n, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
		if !ok {
			t.Fatalf("channel closed on iteration %d", i)
		}

		if n != 5 {
			t.Errorf("expected 5, got %d", n)
		}
	}
}

// TestRandNumGenerator_ImmediateStop проверяет остановку сразу после создания
func TestRandNumGenerator_ImmediateStop(t *testing.T) {
	gen := NewRandNumGenerator(0, 10)
	gen.Stop()

	// даём горутине время завершиться
	time.Sleep(50 * time.Millisecond)

	_, ok := recvWithTimeout(t, gen.Next(), 100*time.Millisecond)
	if ok {
		t.Error("expected closed channel after immediate Stop")
	}
}
