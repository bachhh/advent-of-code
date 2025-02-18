package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"aoc2024/util"
)

type Pair = util.Pair

const (
	maxX = 101
	maxY = 103
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	file, err := os.Open("test_input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var ogMatrix [][]byte
	disableTransform := true

	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes) == 0 {
			break
		}
		ogMatrix = append(ogMatrix, util.CloneSlice(bytes))
	}
	var moves []byte
	for scanner.Scan() {
		bytes := scanner.Bytes()
		moves = append(moves, bytes...)
	}
	var matrix [][]byte
	if disableTransform {
		matrix = util.CloneMatrix(ogMatrix)
	} else {
		matrix = util.NewMatrix[byte](len(ogMatrix), len(ogMatrix[0])*2)
		for i := range ogMatrix {
			for j := range ogMatrix[i] {
				switch ogMatrix[i][j] {
				case '#':
					matrix[i][j*2] = '#'
					matrix[i][j*2+1] = '#'
				case 'O':
					matrix[i][j*2] = '['
					matrix[i][j*2+1] = ']'
				case '.':
					matrix[i][j*2] = '.'
					matrix[i][j*2+1] = '.'
				case '@':
					matrix[i][j*2] = '@'
					matrix[i][j*2+1] = '.'
				default:
					panic("????")
				}
			}
		}

	}

	var cur Pair
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == '@' {
				cur.Row, cur.Col = i, j
			}
		}
	}

	util.PrintMatrix(matrix)
	for _, move := range moves {
		dir := dirMap[move]
		next := util.Walk(cur, dir)
		fmt.Println(string(move))

		// robot walk into the wall
		if next.Row <= 0 || next.Row >= len(matrix)-1 || next.Col <= 0 || next.Col >= len(matrix[0])-1 {
		} else if matrix[next.Row][next.Col] == '#' {
			fmt.Println("wall")
		} else if matrix[next.Row][next.Col] == '.' {
			matrix[next.Row][next.Col] = '@'
			matrix[cur.Row][cur.Col] = '.'
			cur = next
			fmt.Println("free")
		} else if push(matrix, next, dir, false) {
			matrix[next.Row][next.Col] = '@'
			matrix[cur.Row][cur.Col] = '.'
			cur = next
			fmt.Println("box")
		} else {
			fmt.Println("cannot push box")
			// unpushable box
		}

		util.PrintMatrix(matrix)
	}
	total := 0
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == 'O' {
				total += 100*i + j
			}
		}
	}

	util.PrintMatrix(matrix)
	fmt.Println(string(moves))
	fmt.Println(total)
}

func push(matrix [][]byte, cur Pair, dir util.Direction, fromNextCell bool) bool {
	// since map is twice as wide, all boxes occupy 2 cells
	// if fromNextCell, we already came from the 2nd cell of the box, so no need to recurse here
	// if push to east / west, no need to check the 2nd cell, all boxes lie on the same line
	if !fromNextCell && (dir == util.North || dir == util.South) {
		if matrix[cur.Row][cur.Col] == '[' {
			if !push(matrix, Pair{Row: cur.Row, Col: cur.Col + 1}, dir, true) {
				return false
			}
		} else if matrix[cur.Row][cur.Col] == ']' {
			if !push(matrix, Pair{Row: cur.Row, Col: cur.Col - 1}, dir, true) {
				return false
			}
		} else {
			panic("pushing something not a box")
		}
	}
	next := util.Walk(cur, dir)
	// if the next object is a wall, cannot pushed
	if next.Row <= 0 || next.Row >= len(matrix)-1 || next.Col <= 0 || next.Col >= len(matrix[0])-1 {
		return false
	} else if matrix[next.Row][next.Col] == '#' {
		return false
	}

	// we have a blank space on next, move the current object to that next space
	// since we don't have previous object info, leave the moving to the parent
	if matrix[next.Row][next.Col] == '.' {
		matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		return true
	} else if push(matrix, next, dir, false) {
		// we have another box, recursively check if the next space can be moved
		// if yes, move the current object to the next space
		matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		return true
	}
	return false
}

var dirMap = map[byte]util.Direction{
	'^': util.North,
	'<': util.West,
	'>': util.East,
	'v': util.South,
}
