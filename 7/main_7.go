package main

import (
	"fmt"
	"sync"
	"time"
)

func MergeChannels(channels ...<-chan int) <-chan int {
	if len(channels) == 0 {
		return nil
	}
	resChan := make(chan int)
	var wg sync.WaitGroup

	for _, channel := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				resChan <- v
			}
		}(channel)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()
	return resChan
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			ch1 <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch1)
	}()
	go func() {
		for i := 10; i <= 12; i++ {
			ch2 <- i
			time.Sleep(150 * time.Millisecond)
		}
		close(ch2)
	}()
	go func() {
		for i := 100; i <= 102; i++ {
			ch3 <- i
			time.Sleep(80 * time.Millisecond)
		}
		close(ch3)
	}()

	merged := MergeChannels(ch1, ch2, ch3)
	fmt.Println("Значения из слитого канала:")
	for val := range merged {
		fmt.Println(val)
	}
}
