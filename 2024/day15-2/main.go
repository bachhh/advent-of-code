package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"aoc2024/util"
)

type Pair = util.Pair

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

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var ogMatrix [][]byte
	disableTransform := false

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
	fmt.Println(len(matrix))
	for _, move := range moves {
		dir := dirMap[move]
		next := util.Walk(cur, dir)

		// robot walk into the wall
		if next.Row <= 0 || next.Row >= len(matrix)-1 || next.Col <= 0 || next.Col >= len(matrix[0])-1 {
		} else if matrix[next.Row][next.Col] == '#' {
			fmt.Println("wall")
		} else if matrix[next.Row][next.Col] == '.' {

			matrix[next.Row][next.Col] = '@'
			matrix[cur.Row][cur.Col] = '.'
			cur = next
			fmt.Println("free")
		} else if checkPush(matrix, next, dir) {
			push(matrix, next, dir)

			fmt.Println("box")
			util.PrintMatrix(matrix)

			matrix[cur.Row][cur.Col] = '.'
			matrix[next.Row][next.Col] = '@'
			cur = next
			fmt.Println("box2")
			util.PrintMatrix(matrix)
		} else {
			fmt.Println("cannot push box")
			// unpushable box
		}

		util.PrintMatrix(matrix)
	}
	total := 0

	util.PrintMatrix(matrix)
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == '[' {
				leftDis := j
				total += (i)*100 + leftDis
			}
		}
	}

	util.PrintMatrix(matrix)
	// time.Sleep(time.Second)
	fmt.Println(string(moves))
	fmt.Println(total)
}

func checkPush(matrix [][]byte, cur Pair, dir util.Direction) bool {
	// special case for pushing box north / south
	if dir == util.North || dir == util.South {
		// left normalized
		if matrix[cur.Row][cur.Col] == ']' {
			cur = Pair{Row: cur.Row, Col: cur.Col - 1}
		}
		next1 := util.Walk(cur, dir)
		// if the next object is a wall, cannot pushed
		if next1.Row <= 0 || next1.Row >= len(matrix)-1 || next1.Col <= 0 || next1.Col >= len(matrix[0])-1 {
			return false
		} else if matrix[next1.Row][next1.Col] == '#' {
			return false
		}

		// we have a box, check recursively
		if matrix[next1.Row][next1.Col] != '.' {
			if !checkPush(matrix, next1, dir) {
				return false
			}
		}
		// blank space, no need to check

		next2 := util.Walk(Pair{Row: cur.Row, Col: cur.Col + 1}, dir)
		// if the next object is a wall, cannot pushed
		if next2.Row <= 0 || next2.Row >= len(matrix)-1 || next2.Col <= 0 || next2.Col >= len(matrix[0])-1 {
			return false
		} else if matrix[next2.Row][next2.Col] == '#' {
			return false
		}

		// we have a box, check recursively
		if matrix[next2.Row][next2.Col] != '.' {
			if !checkPush(matrix, next2, dir) {
				return false
			}
		}

		// now do the pushing
		return true
	}

	// just pushing left or right
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
		// matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		return true
	}

	if checkPush(matrix, next, dir) {
		// we have another box, recursively check if the next space can be moved
		// if yes, move the current object to the next space
		// matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		return true
	}
	return false
}

func push(matrix [][]byte, cur Pair, dir util.Direction) bool {
	fmt.Println(cur)
	// special case for pushing box north / south
	if dir == util.North || dir == util.South {
		// left normalized
		if matrix[cur.Row][cur.Col] == ']' {
			println(runtime.Caller(0))
			cur = Pair{Row: cur.Row, Col: cur.Col - 1}
		}
		next1 := util.Walk(cur, dir)

		// we have a box, check recursively
		if matrix[next1.Row][next1.Col] != '.' {
			push(matrix, next1, dir)
		}
		// blank space, no need to check

		cur2 := Pair{Row: cur.Row, Col: cur.Col + 1}
		next2 := util.Walk(cur2, dir)

		// we have a box, check recursively
		if matrix[next2.Row][next2.Col] != '.' {
			push(matrix, next2, dir)
		}

		fmt.Println("_____1______")
		util.PrintMatrix(matrix)
		matrix[next1.Row][next1.Col] = matrix[cur.Row][cur.Col]
		matrix[next2.Row][next2.Col] = matrix[cur2.Row][cur2.Col]
		// now do the pushing
		fmt.Println("_____2______")
		util.PrintMatrix(matrix)
		matrix[cur.Row][cur.Col] = '.'
		matrix[cur2.Row][cur2.Col] = '.'
		return true
	}

	// just pushing left or right
	next := util.Walk(cur, dir)

	// we have a blank space on next, move the current object to that next space
	// since we don't have previous object info, leave the moving to the parent
	if matrix[next.Row][next.Col] == '.' {
		matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		matrix[cur.Row][cur.Col] = '.'
		return true
	}

	if push(matrix, next, dir) {
		// we have another box, recursively check if the next space can be moved
		// if yes, move the current object to the next space
		matrix[next.Row][next.Col] = matrix[cur.Row][cur.Col]
		matrix[cur.Row][cur.Col] = '.'
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
