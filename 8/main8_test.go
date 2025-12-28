package main

import (
	"testing"
	"time"
)

// GetCounter возвращает текущее значение счётчика
func (cwg *CustomWG) GetCounter() int {
	cwg.mtx.Lock()
	defer cwg.mtx.Unlock()
	return cwg.count
}

// TestNewCustomWaitGroup проверяет создание новой CustomWaitGroup
func TestNewCustomWaitGroup(t *testing.T) {
	cwg := NewCustomWG()

	if cwg == nil {
		t.Fatal("NewCustomWaitGroup returned nil")
	}

	if cwg.GetCounter() != 0 {
		t.Errorf("expected counter to be 0, got %d", cwg.GetCounter())
	}

	if cwg.sem == nil {
		t.Error("semaphore is nil")
	}
}

// TestCustomWaitGroup_Add проверяет добавление счётчика
func TestCustomWaitGroup_Add(t *testing.T) {
	tests := []struct {
		name     string
		adds     []int
		expected int
	}{
		{
			name:     "single add",
			adds:     []int{1},
			expected: 1,
		},
		{
			name:     "multiple adds",
			adds:     []int{1, 1, 1},
			expected: 3,
		},
		{
			name:     "add and subtract",
			adds:     []int{5, -2},
			expected: 3,
		},
		{
			name:     "add zero",
			adds:     []int{0},
			expected: 0,
		},
		{
			name:     "large number",
			adds:     []int{1000},
			expected: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwg := NewCustomWG()

			for _, delta := range tt.adds {
				cwg.Add(delta)
			}

			if cwg.GetCounter() != tt.expected {
				t.Errorf("expected counter %d, got %d", tt.expected, cwg.GetCounter())
			}
		})
	}
}

// TestCustomWaitGroup_Done проверяет метод Done
func TestCustomWaitGroup_Done(t *testing.T) {
	cwg := NewCustomWG()
	cwg.Add(3)

	if cwg.GetCounter() != 3 {
		t.Fatalf("expected counter 3, got %d", cwg.GetCounter())
	}

	cwg.Done()
	if cwg.GetCounter() != 2 {
		t.Errorf("after Done(), expected counter 2, got %d", cwg.GetCounter())
	}

	cwg.Done()
	cwg.Done()
	if cwg.GetCounter() != 0 {
		t.Errorf("after 3x Done(), expected counter 0, got %d", cwg.GetCounter())
	}
}

// TestCustomWaitGroup_Wait проверяет блокирование Wait()
func TestCustomWaitGroup_Wait(t *testing.T) {
	cwg := NewCustomWG()
	cwg.Add(1)

	done := make(chan bool)

	go func() {
		cwg.Wait()
		done <- true
	}()

	// даём горутине время начать ждать
	time.Sleep(100 * time.Millisecond)

	// проверяем, что горутина всё ещё ждет (done канал не сигнализировал)
	select {
	case <-done:
		t.Fatal("Wait() returned too early")
	case <-time.After(100 * time.Millisecond):
		// нормально, горутина ещё ждит
	}

	cwg.Done()

	// теперь горутина должна завершиться
	select {
	case <-done:
		// горутина завершилась
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait() did not return after Done()")
	}
}

// TestCustomWaitGroup_WaitWithZeroCounter проверяет Wait() когда счётчик уже 0
func TestCustomWaitGroup_WaitWithZeroCounter(t *testing.T) {
	cwg := NewCustomWG()

	done := make(chan bool)

	go func() {
		cwg.Wait() // счётчик уже 0, не должен блокироваться
		done <- true
	}()

	select {
	case <-done:
		// нормально, Wait() не блокировалась
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait() blocked when counter was already 0")
	}
}

// TestCustomWaitGroup_NegativeCounterPanic проверяет паники при отрицательном счётчике
func TestCustomWaitGroup_NegativeCounterPanic(t *testing.T) {
	cwg := NewCustomWG()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative counter, but got none")
		}
	}()

	cwg.Add(-1)
}

// TestCustomWaitGroup_DoneWithoutAdd проверяет Done() без предварительного Add()
func TestCustomWaitGroup_DoneWithoutAdd(t *testing.T) {
	cwg := NewCustomWG()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Done without Add, but got none")
		}
	}()

	cwg.Done() // должно вызвать панику: счётчик -1
}
