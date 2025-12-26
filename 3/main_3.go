package main

import "fmt"

type StringIntMapInterface interface {
	Add(key string, value int)
	Remove(key string)
	Copy() map[string]int
	Exists(key string) bool
	Get(key string) (int, bool)
}

type StringIntMap struct {
	data map[string]int
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{
		data: make(map[string]int),
	}
}

func (m *StringIntMap) Add(key string, value int) {
	m.data[key] = value
}

func (m *StringIntMap) Remove(key string) {
	delete(m.data, key)
}

func (m *StringIntMap) Copy() map[string]int {
	var newMap = make(map[string]int)
	for k, v := range m.data {
		newMap[k] = v
	}
	return newMap
}

func (m *StringIntMap) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *StringIntMap) Get(key string) (int, bool) {
	val, ok := m.data[key]
	return val, ok
}

func main() {
	SIMap := NewStringIntMap()
	fmt.Println("Исходная мапа:", SIMap)

	SIMap.Add("first", 11)
	SIMap.Add("second", 22)
	SIMap.Add("third", 33)
	fmt.Println("\nДобавили данные:", SIMap.data)

	SIMap.Remove("second")
	fmt.Println("\nУбрали second из мапы:", SIMap.data)

	copiedMap := SIMap.Copy()
	SIMap.Add("fourth", 44)
	fmt.Println("\nСкопированная мапа:", copiedMap)
	fmt.Println("Исходная мапа:", SIMap.data)

	fmt.Println("Есть ли элемент forth в мапе:", SIMap.Exists("fourth"))
	fmt.Println("Есть ли элемент fifth в мапе:", SIMap.Exists("fifth"))

	val, ok := SIMap.Get("first")
	fmt.Println("Есть ли элемент first в мапе:", ok, "он равен:", val)
}
