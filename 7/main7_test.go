package main

import (
	"slices"
	"sort"
	"testing"
	"time"
)

// helper: читаем все значения из канала с таймаутом
func collectValues(t *testing.T, ch <-chan int, timeout time.Duration) []int {
	t.Helper()

	var result []int
	done := make(chan bool)

	go func() {
		for val := range ch {
			result = append(result, val)
		}
		done <- true
	}()

	select {
	case <-done:
		return result
	case <-time.After(timeout):
		t.Fatalf("timeout collecting values from channel")
		return nil
	}
}

// TestMerge_SingleChannel проверяет слияние одного канала
func TestMerge_SingleChannel(t *testing.T) {
	ch := make(chan int)

	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
	}()

	merged := MergeChannels(ch)
	result := collectValues(t, merged, 1*time.Second)

	expected := []int{1, 2, 3}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestMerge_MultipleChannels проверяет слияние нескольких каналов
func TestMerge_MultipleChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	go func() {
		ch2 <- 10
		ch2 <- 20
		close(ch2)
	}()

	go func() {
		ch3 <- 100
		ch3 <- 200
		close(ch3)
	}()

	merged := MergeChannels(ch1, ch2, ch3)
	result := collectValues(t, merged, 1*time.Second)

	sort.Ints(result)
	expected := []int{1, 2, 10, 20, 100, 200}

	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestMerge_EmptyChannels проверяет слияние пустых каналов
func TestMerge_EmptyChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		close(ch1)
	}()

	go func() {
		close(ch2)
	}()

	merged := MergeChannels(ch1, ch2)
	result := collectValues(t, merged, 1*time.Second)

	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

// TestMerge_ChannelsClosedImmediately проверяет закрытие выходного канала
func TestMerge_ChannelsClosedImmediately(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	close(ch1)
	close(ch2)

	merged := MergeChannels(ch1, ch2)

	// должен закрыться почти сразу
	select {
	case val, ok := <-merged:
		if ok {
			t.Errorf("expected channel to be closed, got value %d", val)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for channel to close")
	}
}

// TestMerge_DifferentNumberOfElements проверяет каналы с разным количеством элементов
func TestMerge_DifferentNumberOfElements(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		ch1 <- 1
		close(ch1)
	}()

	go func() {
		ch2 <- 10
		ch2 <- 20
		ch2 <- 30
		close(ch2)
	}()

	go func() {
		ch3 <- 100
		ch3 <- 200
		close(ch3)
	}()

	merged := MergeChannels(ch1, ch2, ch3)
	result := collectValues(t, merged, 1*time.Second)

	sort.Ints(result)
	expected := []int{1, 10, 20, 30, 100, 200}

	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestMerge_LargeNumberOfChannels проверяет слияние большого количества каналов
func TestMerge_LargeNumberOfChannels(t *testing.T) {
	const numChannels = 100
	const valuesPerChannel = 10

	channels := make([]<-chan int, numChannels)

	for i := 0; i < numChannels; i++ {
		ch := make(chan int)
		channels[i] = ch

		go func(chNum int, c chan int) {
			for j := 0; j < valuesPerChannel; j++ {
				c <- chNum*1000 + j
			}
			close(c)
		}(i, ch)
	}

	merged := MergeChannels(channels...)
	result := collectValues(t, merged, 2*time.Second)

	if len(result) != numChannels*valuesPerChannel {
		t.Errorf("expected %d values, got %d", numChannels*valuesPerChannel, len(result))
	}
}

// TestMerge_ConcurrentReading проверяет конкурентное чтение из слитого канала
func TestMerge_ConcurrentReading(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 50; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 100; i < 150; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	merged := MergeChannels(ch1, ch2)

	// читаем из разных горутин одновременно
	ch1Result := make(chan []int)
	ch2Result := make(chan []int)

	go func() {
		var vals []int
		for i := 0; i < 50; i++ {
			vals = append(vals, <-merged)
		}
		ch1Result <- vals
	}()

	go func() {
		var vals []int
		for i := 0; i < 50; i++ {
			vals = append(vals, <-merged)
		}
		ch2Result <- vals
	}()

	r1 := <-ch1Result
	r2 := <-ch2Result

	total := append(r1, r2...)
	if len(total) != 100 {
		t.Errorf("expected 100 values, got %d", len(total))
	}
}

// TestMerge_PartialClosing проверяет поведение при закрытии некоторых каналов
func TestMerge_PartialClosing(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	go func() {
		ch2 <- 10
		ch2 <- 20
		close(ch2)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch3 <- 100
		close(ch3)
	}()

	merged := MergeChannels(ch1, ch2, ch3)
	result := collectValues(t, merged, 1*time.Second)

	sort.Ints(result)
	expected := []int{1, 2, 10, 20, 100}

	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
