package main

import (
	"bufio"
	"fmt"
	"os"
)

// a lazy way to think about this is we have basically 4 distinc matrices, accounting for inversion
var (
	east = [][]byte{
		{'M', '.', 'S'},
		{'.', 'A', '.'},
		{'M', '.', 'S'},
	}
	west = [][]byte{
		{'S', '.', 'M'},
		{'.', 'A', '.'},
		{'S', '.', 'M'},
	}
	south = [][]byte{
		{'M', '.', 'M'},
		{'.', 'A', '.'},
		{'S', '.', 'S'},
	}
	north = [][]byte{
		{'S', '.', 'S'},
		{'.', 'A', '.'},
		{'M', '.', 'M'},
	}
)

// just move through the original matrix and check if it contains any of these 4 matrices
func compareMatrix(matrix [][]byte, row, col int, target [][]byte) bool {
	// bound check
	if col+len(target) > len(matrix) || row+len(target[0]) > len(matrix[0]) {
		return false
	}

	for r := range target {
		for c := range target {
			if target[r][c] == '.' {
				// match any char is fine
				continue
			} else if matrix[row+r][col+c] != target[r][c] {
				return false
			}
		}
	}
	return true
}

func main() {
	var matrix [][]byte

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Stop reading on an empty line (optional)
			break
		}
		matrix = append(matrix, []byte(line))
	}

	counter := 0
	for row := range matrix {
		for col := range matrix[0] {
			if compareMatrix(matrix, row, col, east) {
				counter++
			}
			if compareMatrix(matrix, row, col, south) {
				counter++
			}
			if compareMatrix(matrix, row, col, north) {
				counter++
			}
			if compareMatrix(matrix, row, col, west) {
				counter++
			}
		}
	}
	fmt.Println(counter)
}
