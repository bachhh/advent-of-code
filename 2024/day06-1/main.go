package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	matrix := [][]byte{}

	printMatrix(matrix)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		cp := make([]byte, len(line))
		copy(cp, line)
		matrix = append(matrix, cp)
		// printMatrix(matrix)
		// fmt.Println(string(line))
	}
	fmt.Println(len(matrix), len(matrix[0]))
	printMatrix(matrix)

	counter := 1
	x, y := getPos(matrix, '^')
	fmt.Println(x, y)
	// printMatrix(matrix)
	var dir int // 0,1,2,3 = north, east south, west
	dir = 0
	for {

		var nextX, nextY int
		switch dir {
		case 0:
			nextX = x - 1
			nextY = y
		case 1:
			nextX = x
			nextY = y + 1
		case 2:
			nextX = x + 1
			nextY = y
		case 3:
			nextX = x
			nextY = y - 1
		}
		fmt.Println(x, y, nextX, nextY, dir)
		if nextX < 0 || nextY < 0 || nextX >= len(matrix) || nextY >= len(matrix[0]) {
			break // out of bound
		}
		if matrix[nextX][nextY] == '#' {
			dir = (dir + 1) % 4
			continue
		}

		x, y = nextX, nextY
		if matrix[x][y] == '.' {
			counter++
		}
		matrix[x][y] = 'X'
		// fmt.Println(counter)
		// printMatrix(matrix)
	}
	// printMatrix(matrix)
	fmt.Println(counter)
}

func printMatrix(matrix [][]byte) {
	fmt.Println()
	for i := range matrix {
		fmt.Println(string(matrix[i]))
	}
	fmt.Println()
}

func getPos(matrix [][]byte, target byte) (int, int) {
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == target {
				return i, j
			}
		}
	}
	return -1, -1
}
