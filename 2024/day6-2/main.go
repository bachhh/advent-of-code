package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var (
	debug         = false
	delay         = 50 * time.Millisecond
	counter int32 = 0
	// maxStep       = 10_000
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

	tracePath := [][]int{}
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var x, y int
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		in := make([]int, len(line))
		for j, char := range line {
			if char == '#' {
				in[j] |= (1 << 5)
			} else if char == '^' {
				x, y = i, j
			}
		}
		tracePath = append(tracePath, in)
	}

	var dir int // 0,1,2,3 = north, east south, west
	dir = 0
	setTrace(tracePath, x, y, dir)
	cpPath := make([][]int, len(tracePath))
	for i := range tracePath {
		cpPath[i] = make([]int, len(tracePath[i]))
	}
	for i := 0; ; i++ {
		setTrace(tracePath, x, y, dir)

		nextX, nextY := walk(x, y, dir)
		if nextX < 0 || nextY < 0 || nextX >= len(tracePath) || nextY >= len(tracePath[0]) {
			break // out of bound
		}

		// printTracePath(tracePath)
		if checkObstacle(tracePath, nextX, nextY) {
			dir = (dir + 1) % 4
			continue
		}

		// jobChan <- jobParam{tracePath, x, y, dir}
		if canPlaceObs(tracePath, x, y, dir, debug, cpPath) {
			// waitForEnter()
			// obsX, obsY := walk(nextX, nextY, dir)
			counter++
		}
		x, y = nextX, nextY

		// time.Sleep(delay)
	}
	printTracePath(tracePath)
	fmt.Println(counter)
}

// check if we can place an obstacle in front of the guard, at the current position.
func canPlaceObs(tracePath [][]int, x, y int, dir int, debug bool, cpPath [][]int) bool {
	obsX, obsY := walk(x, y, dir)
	if isAnyFirst4BitsSet(tracePath[obsX][obsY]) {
		return false
	}

	copyMatrixDest(tracePath, cpPath)
	nX, nY := walk(x, y, dir)
	setObstacle(cpPath, nX, nY)
	// assume we have an obstacle, now turn right
	dir = (dir + 1) % 4

	// simulate a full walk
	for {
		nextX, nextY := walk(x, y, dir)
		setTrace(cpPath, x, y, dir)

		if nextX < 0 || nextY < 0 || nextX >= len(cpPath) || nextY >= len(cpPath[0]) {
			// even if we place obstacle here, guard can still walk out of map
			return false
		}

		// found a cycle
		if checkTrace(cpPath, nextX, nextY, dir) {
			if debug {
				printTracePath(cpPath)
				fmt.Printf("possible obstacle at %d, %d\n", nX, nY)
			}
			return true
		}

		if checkObstacle(cpPath, nextX, nextY) {
			dir = (dir + 1) % 4
			continue
		}

		x, y = nextX, nextY
		// time.Sleep(delay)
	}
}

func copyMatrixDest[T any](src [][]T, dst [][]T) {
	for i := range src {
		copy(dst[i], src[i])
	}
}

func printTracePath(tracePath [][]int) {
	fmt.Print("\033[H\033[2J")
	for i := range tracePath {
		for j := range tracePath[i] {
			if checkObstacle(tracePath, i, j) {
				fmt.Printf("%2c", '#')
				// } else if checkPlaceObstacle(tracePath, i, j) {
				// 	fmt.Printf("%2c", 'O')
			} else if tracePath[i][j] == 0 {
				fmt.Printf("%2c", '.')
			} else {
				fmt.Printf("%2c", pathSign[tracePath[i][j]])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// n e s w
// 1 2 3 4
// 1 2 4 8
var pathSign = map[int]rune{
	0x0: '.',
	0x1: '│',
	0x2: '─',
	0x3: '+',
	0x4: '|',
	0x5: '|',
	0x6: '+',
	0x7: '+',
	0x8: '─',
	0x9: '+',
	0xa: '─',
	0xb: '+',
	0xc: '+',
	0xd: '+',
	0xe: '+',
	0xf: '+',
}

func setBit(number, bit int) int {
	return number | (1 << bit)
}

func isBitSet(number, bit int) bool {
	return (number & (1 << bit)) != 0
}

func checkObstacle(lab [][]int, x, y int) bool {
	return isBitSet(lab[x][y], 5)
}

// use the 5th bit for notating obstacle
func setObstacle(lab [][]int, x, y int) {
	lab[x][y] = setBit(lab[x][y], 5)
}

func setTrace(tracePath [][]int, x, y, dir int) {
	tracePath[x][y] = tracePath[x][y] | (1 << dir)
}

func checkTrace(tracePath [][]int, x, y, dir int) bool {
	return (tracePath[x][y] & (1 << dir)) != 0
}

func walk(x, y int, dir int) (int, int) {
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
	return nextX, nextY
}

func waitForEnter() {
	fmt.Println("Press Enter to continue...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Waits for user input (Enter)
}

func isAnyFirst4BitsSet(n int) bool {
	return (n & 0b1111) != 0
}
