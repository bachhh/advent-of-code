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

	var matrix [][]byte

	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes) == 0 {
			break
		}
		matrix = append(matrix, util.CloneSlice(bytes))
	}

	var moves []byte
	for scanner.Scan() {
		bytes := scanner.Bytes()
		moves = append(moves, bytes...)
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
			fmt.Println("boundary")
		} else if matrix[next.Row][next.Col] == '#' {
			fmt.Println("wall")
		} else if matrix[next.Row][next.Col] == '.' {
			matrix[next.Row][next.Col] = '@'
			matrix[cur.Row][cur.Col] = '.'
			cur = next
			fmt.Println("_____1______")
		} else if push(matrix, next, dir) {
			matrix[next.Row][next.Col] = '@'
			matrix[cur.Row][cur.Col] = '.'
			cur = next
			fmt.Println("_____2______")
		} else {
			fmt.Println("cannot push box")
			// unpushable box
		}

		util.PrintMatrix(matrix)
	}

	util.PrintMatrix(matrix)
	fmt.Println(string(moves))
}

func push(matrix [][]byte, cur Pair, dir util.Direction) bool {
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
	} else if push(matrix, next, dir) {
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
