package main

import (
	"fmt"
	"sync"
	"time"
)

type CustomWG struct {
	sem   chan struct{}
	count int
	mtx   sync.Mutex
}

func NewCustomWG() *CustomWG {
	return &CustomWG{
		count: 0,
		sem:   make(chan struct{}),
	}
}

func (cwg *CustomWG) Add(delta int) {
	cwg.mtx.Lock()
	defer cwg.mtx.Unlock()

	cwg.count += delta
	if cwg.count < 0 {
		panic("Number of goroutines is negative")
	}
	if cwg.count == 0 && len(cwg.sem) == 0 {
		select {
		case cwg.sem <- struct{}{}:
		default:
		}
	}
}

func (cwg *CustomWG) Done() {
	cwg.Add(-1)
}

func (cwg *CustomWG) Wait() {
	if cwg.count == 0 {
		return
	}
	<-cwg.sem
}

func main() {
	fmt.Println("Создание кастомной Waitgroup")
	cwg := NewCustomWG()

	fmt.Println("Запуск 5 горутин...")
	for i := 1; i <= 5; i++ {
		cwg.Add(1)

		go func(id int) {
			defer cwg.Done()

			fmt.Printf("Горутина %d: начало работы\n", id)
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("Горутина %d: конец работы\n", id)
		}(i)
	}

	fmt.Println("Основной поток: ждём завершения всех горутин...")
	cwg.Wait()
	fmt.Println("Основной поток: все горутины завершились!")
}
