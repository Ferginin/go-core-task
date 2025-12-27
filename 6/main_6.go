package main

import (
	"fmt"
	"math/rand"
	"time"
)

type RandNumGenerator struct {
	numChan chan int
	stop    chan struct{}
}

func NewRandNumGenerator(min, max int) *RandNumGenerator {
	if min > max {
		min, max = max, min
	}
	gen := &RandNumGenerator{
		numChan: make(chan int),
		stop:    make(chan struct{}),
	}
	go func() {
		defer close(gen.numChan)
		for {
			select {
			case <-gen.stop:
				return
			case gen.numChan <- rand.Intn(max-min+1) + min:
			}
		}
	}()
	return gen
}

func (gen *RandNumGenerator) Next() <-chan int {
	return gen.numChan
}

func (gen *RandNumGenerator) Stop() {
	select {
	case <-gen.stop:
		return
	default:
		close(gen.stop)
	}
}

func main() {
	fmt.Println("Создание генератора случайных чисел")
	gen := NewRandNumGenerator(20, 100)
	defer gen.Stop()
	for i := 0; i < 10; i++ {
		n := <-gen.Next()
		fmt.Println("Число - ", n)
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("Остановка генератора...")
}
