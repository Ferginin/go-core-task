package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

func MakeNumsConveyor(ch1 <-chan uint8) <-chan float64 {
	chOut := make(chan float64)
	go func() {
		defer close(chOut)
		for num := range ch1 {
			outNum := float64(num)
			chOut <- math.Pow(outNum, 3)
		}
	}()

	return chOut
}

func main() {
	chIn := make(chan uint8)
	chOut := MakeNumsConveyor(chIn)

	go func() {
		defer close(chIn)
		for i := range 10 {
			time.Sleep(time.Millisecond * 200)
			n := uint8(rand.UintN(20))
			chIn <- n
			fmt.Printf("Значение uint8 №%d на первом конвеере: %d\n", i, n)
		}
	}()

	for num := range chOut {
		fmt.Printf("Куб от значения: %f\n\n", num)
	}
}
