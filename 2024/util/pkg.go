package util

import (
	"bufio"
	"cmp"
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

// same as print matrix, but take any function and and transform to character
func PrintMatrixTransform[T any](isRefresh bool, matrix [][]T, transform func(c T) string) {
	if isRefresh {
		fmt.Print("\033[H\033[2J")
	}
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%s", transform(matrix[i][j]))
		}
		fmt.Println()
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

var ManhattanDirs = []Direction{
	North,
	East,
	South,
	West,
}

var ChebysevDirs = []Direction{
	North,
	East,
	South,
	West,
	NorthEast,
	SouthEast,
	SouthWest,
	NorthWest,
}

// String method to implement fmt.Stringer interface
func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	case NorthEast:
		return "NorthEast"
	case SouthEast:
		return "SouthEast"
	case SouthWest:
		return "SouthWest"
	case NorthWest:
		return "NorthWest"
	default:
		return "Unknown Direction"
	}
}

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

func ManhattanDistance(from, to Pair) int {
	return Abs(from.Row-to.Row) + Abs(from.Col-to.Col)
}

func IsPairInbound[T any](point Pair, matrix [][]T) bool {
	return point.Row >= 0 && point.Row < len(matrix) && point.Col >= 0 && point.Col < len(matrix[0])
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Abs[T Number](val T) T {
	if val < 0 {
		return -val
	}
	return val
}

func NewMatrix[T any](row, col int) [][]T {
	matrix := make([][]T, row)
	for i := range matrix {
		matrix[i] = make([]T, col)
	}
	return matrix
}

func CloneMatrix[T any](matrix [][]T) [][]T {
	cloneMatrix := make([][]T, len(matrix))
	for i := range matrix {
		cloneMatrix[i] = make([]T, len(matrix[i]))
		copy(cloneMatrix[i], matrix[i])
	}
	return cloneMatrix
}

func ScanMatrix(scanner *bufio.Scanner) [][]byte {
	var matrix [][]byte

	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes) == 0 {
			break
		}
		matrix = append(matrix, CloneSlice(bytes))
	}
	return matrix
}

func FindMatrix[T any](matrix [][]T, foo func(T) bool) (int, int) {
	for i := range matrix {
		for j := range matrix[i] {
			if foo(matrix[i][j]) {
				return i, j
			}
		}
	}
	return -1, -1
}

func CloneMap[K comparable, V any](m map[K]V) map[K]V {
	cl := map[K]V{}
	for k, v := range m {
		cl[k] = v
	}
	return cl
}

func IsCollinear(a, b, c Pair) bool {
	// Check vertical alignment
	if a.Row == b.Row && b.Row == c.Row {
		return IsInIntervalIncl(b.Col, a.Col, c.Col)
	}

	// Check horizontal alignment
	if a.Col == b.Col && b.Col == c.Col {
		return IsInIntervalIncl(b.Row, a.Row, c.Row)
	}

	// Not collinear in horizontal or vertical direction
	return false
}

// check if A is in the interval [x, y] inclusive
func IsInIntervalIncl[T cmp.Ordered](a, x, y T) bool {
	if x > y {
		x, y = y, x
	}
	return a >= x && a <= y
}

// same as IsInIntervalIncl but exclusive, .e.g (x, y)
func IsInIntervalExcl[T cmp.Ordered](a, x, y T) bool {
	if x > y {
		x, y = y, x
	}
	return a > x && a < y
}
