package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type Variables struct {
	numDecimal     int64 // Десятичная система
	numOctal       int64 // Восьмеричная система
	numHexadecimal int64 // Шестнадцатиричная система
	floatVar       float64
	stringVar      string
	boolVar        bool
	complexNum     complex64
}

func main() {
	Vars := InitVariables()

	Vars.PrintType()

	concatVars := Vars.ConcatVariables()
	fmt.Println("\nString from variables:", concatVars)

	runes := []rune(concatVars)
	fmt.Println("Runes from string:", runes)

	saltedRunes := InsertSalt(runes, "go-2024")
	fmt.Println("Salted runes:", saltedRunes)

	hashedRunes := HashRune(saltedRunes)
	fmt.Println("Hashed runes:", hashedRunes)
}

func InitVariables() Variables {
	return Variables{
		numDecimal:     42,
		numOctal:       075,
		numHexadecimal: 0xFA,
		floatVar:       3.14,
		stringVar:      "Golang",
		boolVar:        true,
		complexNum:     1 + 2i,
	}
}

func (v Variables) PrintType() {
	fmt.Printf("Тип переменной numDecimal: %T\n", v.numDecimal)
	fmt.Printf("Тип переменной numOctal: %T\n", v.numOctal)
	fmt.Printf("Тип переменной numHexadecimal: %T\n", v.numHexadecimal)
	fmt.Printf("Тип переменной floatVar: %T\n", v.floatVar)
	fmt.Printf("Тип переменной stringVar: %T\n", v.stringVar)
	fmt.Printf("Тип переменной boolVar: %T\n", v.boolVar)
	fmt.Printf("Тип переменной complexNum: %T\n", v.complexNum)
}

func (v Variables) ConcatVariables() string {
	var builder strings.Builder

	builder.WriteString(strconv.FormatInt(v.numDecimal, 10))
	builder.WriteString(strconv.FormatInt(v.numOctal, 8))
	builder.WriteString(strconv.FormatInt(v.numHexadecimal, 16))
	builder.WriteString(strconv.FormatFloat(v.floatVar, 'f', -1, 64))
	builder.WriteString(v.stringVar)
	builder.WriteString(strconv.FormatBool(v.boolVar))
	builder.WriteString(fmt.Sprintf("%v", v.complexNum))

	return builder.String()
}

func InsertSalt(runes []rune, salt string) []rune {
	saltRunes := []rune(salt)
	middle := len(runes) / 2

	result := make([]rune, 0, len(runes)+len(saltRunes))
	result = append(result, runes[:middle]...)
	result = append(result, saltRunes...)
	result = append(result, runes[middle:]...)

	return result
}

func HashRune(v []rune) string {
	hash := sha256.Sum256([]byte(string(v)))
	return hex.EncodeToString(hash[:])
}
