package main

import (
	"fmt"
	"math/rand"
)

func main() {
	nums := initSlice()
	fmt.Println("Исходный слайс:")
	fmt.Println(nums)

	fmt.Println("\nСлайс только с четными числами:")
	fmt.Println(sliceExample(nums))

	fmt.Println("\nДобавление элемента 52 в конец слайса:")
	fmt.Println(addElements(nums, 52))
	fmt.Println("Исходный слайс:")
	fmt.Println(nums)

	fmt.Println("\nКопирование слайса:")
	copiedNums := copySlice(nums)
	fmt.Println(copiedNums)
	fmt.Println("Исходный слайс:")
	fmt.Println(nums)
	fmt.Println("Изменим исходный слайс:")
	nums[0] = 111
	fmt.Println(nums)
	fmt.Println("Скопированный слайс:")
	fmt.Println(copiedNums)

	fmt.Println("\nУдаление элемента на 4 индексе:")
	fmt.Println(removeElement(nums, 4))
	fmt.Println("Исходный слайс:")
	fmt.Println(nums)
}

func initSlice() []int {
	nums := make([]int, 10)
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	return nums
}

func sliceExample(nums []int) []int {
	var oddNums []int
	for _, n := range nums {
		if n%2 == 0 {
			oddNums = append(oddNums, n)
		}
	}
	return oddNums
}

func addElements(nums []int, num int) []int {
	newNums := make([]int, len(nums))
	copy(newNums, nums)
	newNums = append(newNums, num)
	return newNums
}

func copySlice(nums []int) []int {
	newNums := make([]int, len(nums))
	copy(newNums, nums)
	return newNums
}

func removeElement(nums []int, indx int) []int {
	newNums := copySlice(nums)
	newNums = append(newNums[:indx], newNums[indx+1:]...)
	return newNums
}
