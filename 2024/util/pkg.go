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

// SwapSlice swap all element of slice[a:b] with slice[c:d]
func SwapSlice[T ~[]E, E any](s T, a, b, c, d int) error {
	if a >= b || c >= d || b > len(s) || d > len(s) {
		return fmt.Errorf("invalid indices")
	}
	if c-d > b-a {
		return fmt.Errorf("2 sub-slice have different sizes")
	}
	tmp := make(T, b-a)
	copy(tmp, s[a:b])
	copy(s[a:b], s[c:d])
	copy(s[c:d], tmp)
	return nil
}
