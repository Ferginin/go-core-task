package main

import "fmt"

func FindIntersection(slice1, slice2 []int) (bool, []int) {
	var resSlice []int
	s1Map := make(map[int]bool)
	for _, v := range slice1 {
		s1Map[v] = true
	}
	s2Map := make(map[int]bool)
	for _, v := range slice2 {
		s2Map[v] = true
	}

	for k := range s2Map {
		if s1Map[k] {
			resSlice = append(resSlice, k)
		}
	}

	if len(resSlice) == 0 {
		return false, resSlice
	}

	return true, resSlice
}

func main() {
	slice1 := []int{65, 3, 58, 678, 64}
	slice2 := []int{64, 2, 3, 43}
	fmt.Println("Два слайса:\nПервый -", slice1, "второй -", slice2)

	ok, result := FindIntersection(slice1, slice2)
	fmt.Println("есть пересечения:", ok, "\nИтоговый слайс:", result)
}
