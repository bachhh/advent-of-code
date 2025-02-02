package util

import (
	"fmt"
	"strconv"
)

func CloneSlice[T ~[]E, E any](slice T) T {
	cp := make(T, len(slice))
	copy(cp, slice)
	return cp
}

func PrintMatrix(matrix [][]byte) {
	for i := range matrix {
		fmt.Println(string(matrix[i]))
	}
	fmt.Println()
}

func PrintMatrixRefresh(matrix [][]byte) {
	fmt.Print("\033[H\033[2J")
	PrintMatrix(matrix)
}

func IsAlphaNumeric(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

func CharToInt(char byte) (int, error) {
	return strconv.Atoi(string(char))
}

func IntToChar(num int) (byte, error) {
	if num == 0 {
		return byte('0'), nil
	}
	if num > 9 {
		panic("can only convert single digit")
	}
	return byte('0' + num), nil
}
