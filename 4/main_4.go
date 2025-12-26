package main

import (
	"fmt"
)

func RemoveIntersection(s1, s2 []string) []string {
	slice2Map := make(map[string]bool)
	var resultSlice []string
	for _, str := range s2 {
		slice2Map[str] = true
	}
	for _, str := range s1 {
		if !slice2Map[str] {
			resultSlice = append(resultSlice, str)
		}
	}
	return resultSlice
}

func main() {
	slice1 := []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}
	slice2 := []string{"banana", "date", "fig"}
	fmt.Println("Созданы два слайса:")
	fmt.Println("Первый слайс:", slice1)
	fmt.Println("Второй слайс:", slice2)

	resSlice := RemoveIntersection(slice1, slice2)
	fmt.Println("Создан новый слайс с элементами из первого слайса, которых нет во втором")
	fmt.Println("Результирующий слайс:", resSlice)
}
