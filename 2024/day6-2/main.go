package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	matrix := [][]byte{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		cp := make([]byte, len(line))
		copy(cp, line)
		matrix = append(matrix, cp)
	}

	counter := 1
	x, y := getPos(matrix, '^')
	matrix[x][y] = 'X'
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
		// fmt.Println(x, y, nextX, nextY, dir)
		if nextX < 0 || nextY < 0 || nextX >= len(matrix) || nextY >= len(matrix[0]) {
			break // out of bound
		}
		if matrix[nextX][nextY] == '#' {
			dir = (dir + 1) % 4
			continue
		}

		if canPlaceObs(matrix, x, y, dir) {
			fmt.Println("placing", x, y)
			matrix[x][y] = '+'
			counter++
		}
		x, y = nextX, nextY
		// if matrix[x][y] == '.' { counter++ }
		matrix[x][y] = 'X'
	}
	printMatrix(matrix)
	fmt.Println(counter)
}

// place that we can place obs is 1. a visited position and b. there is a path
// to the right of facing direction
func canPlaceObs(matrix [][]byte, x, y int, dir int) bool {
	if matrix[x][y] != 'X' {
		return false
	}
	// assume we place an obstacle, now turn right
	dir = (dir + 1) % 4

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
	return matrix[nextX][nextY] == 'X'
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
