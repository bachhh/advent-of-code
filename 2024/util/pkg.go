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

func PrintMatrixInt(matrix [][]byte) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%02d", matrix[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
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

type Direction int

const (
	North Direction = iota
	East
	South
	West
	NorthEast
	SouthEast
	SouthWest
	NorthWest
)

type Pair struct {
	Row, Col int
}

func Walk(point Pair, dir Direction) Pair {
	switch dir {
	case North:
		return Pair{Row: point.Row - 1, Col: point.Col}
	case East:
		return Pair{Row: point.Row, Col: point.Col + 1}
	case South:
		return Pair{Row: point.Row + 1, Col: point.Col}
	case West:
		return Pair{Row: point.Row, Col: point.Col - 1}
	case NorthEast:
		return Pair{Row: point.Row - 1, Col: point.Col + 1}
	case SouthEast:
		return Pair{Row: point.Row + 1, Col: point.Col + 1}
	case SouthWest:
		return Pair{Row: point.Row + 1, Col: point.Col - 1}
	case NorthWest:
		return Pair{Row: point.Row - 1, Col: point.Col - 1}
	}

	panic(fmt.Sprintf("unknown direction %d", dir))
}

func IsPairInbound[T any](point Pair, matrix [][]T) bool {
	return point.Row >= 0 && point.Row < len(matrix) && point.Col >= 0 && point.Col < len(matrix[0])
}
