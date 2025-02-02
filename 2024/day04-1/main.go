package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	forward  = []byte{'X', 'M', 'A', 'S'}
	backward = []byte{'S', 'A', 'M', 'X'}
)

// direction:
// 6 . 5 . 4
// . 6 5 4 .
// 7 7 * 1 1
// . 8 3 2 .
// 8 . 3 . 2
// . . . . .
func check(matrix [][]byte, row, col int, text []byte, direction int) bool {
	switch direction {
	case 1:
		// bound check
		if col+len(text) > len(matrix) {
			return false
		}
		for i := range text {
			if matrix[row][col+i] != text[i] {
				return false
			}
		}
		return true
	case 2:
		if row+len(text) > len(matrix) || col+len(text) > len(matrix[0]) {
			return false
		}
		for i := range text {
			if matrix[row+i][col+i] != text[i] {
				return false
			}
		}
		return true
	case 3:
		if row+len(text) > len(matrix[0]) {
			return false
		}
		for i := range text {
			if matrix[row+i][col] != text[i] {
				return false
			}
		}
		return true
	case 4:
		if row-len(text) < -1 || col+len(text) > len(matrix) {
			return false
		}
		for i := range text {
			if matrix[row-i][col+i] != text[i] {
				return false
			}
		}
		return true
	default:
		return false
	}
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
	allDir := []int{1, 2, 3, 4}
	for row := range matrix {
		for col := range matrix[0] {
			for _, dir := range allDir {
				if check(matrix, row, col, forward, dir) {
					fmt.Println(row, col, dir)
					counter++
				}
				if check(matrix, row, col, backward, dir) {
					fmt.Println(row, col, dir)
					counter++
				}
			}
		}
	}
	fmt.Println(counter)
}
