package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
	"time"

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

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	allRobot := []Robot{}

	for scanner.Scan() {
		x, y, dx, dy, err := ParseLine(scanner.Text())
		if err != nil {
			panic(err)
		}
		allRobot = append(allRobot, Robot{x, y, dx, dy})

	}
	// reader := bufio.NewReader(os.Stdin)
	limit := 8000

	for i := 5000; i < limit; i++ {
		// _, err := reader.ReadString('\n') // Waits for user input
		// if err != nil {
		// 	panic(err)
		// }
		DrawMap(allRobot, i)
		time.Sleep(2 * time.Millisecond)
		fmt.Println(i)
	}
}

// ParseLine extracts position and velocity from a given formatted string.
func ParseLine(line string) (x, y, dx, dy int, err error) {
	// Regex to capture numbers from p=x,y and v=dx,dy
	re := regexp.MustCompile(`-?\d+`)

	// Extract all numbers
	matches := re.FindAllString(line, -1)
	if len(matches) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("invalid line format: %s", line)
	}

	// Convert matches to integers
	numbers := make([]int, 4)
	for i, match := range matches {
		num, err := strconv.Atoi(match)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid number: %s", match)
		}
		numbers[i] = num
	}

	return numbers[0], numbers[1], numbers[2], numbers[3], nil
}

func calcQuadrants(maxX, maxY int) [][]int {
	// q = matrix[row:row][col:col]
	return [][]int{
		{
			0,        // min y
			maxY / 2, // maxy+1
			0,        // minx
			maxX / 2, // maxx+1
		},
		{
			0,
			maxY / 2,
			maxX/2 + 1,
			maxX,
		},
		{
			maxY/2 + 1,
			maxY,
			maxX/2 + 1,
			maxX,
		},
		{
			maxY/2 + 1,
			maxY,
			0,
			maxX / 2,
		},
	}
}

type Robot struct {
	X, Y   int
	Dx, Dy int
}

func (r *Robot) CalcPosition(sec int) (int, int) {
	x, y, dx, dy := r.X, r.Y, r.Dx, r.Dy

	// fmt.Println(x, y, dx, dy)
	finalPosX := (x + sec*dx)
	finalPosY := (y + sec*dy)

	// wrap around if position is negative
	finalPosX += (util.Abs(finalPosX/maxX) + 1) * maxX
	finalPosY += (util.Abs(finalPosY/maxY) + 1) * maxY

	return finalPosX % maxX, finalPosY % maxY
}

func DrawMap(allRobot []Robot, sec int) {
	matrix := util.NewMatrix[int](maxY, maxX)

	for _, robot := range allRobot {
		x, y := robot.CalcPosition(sec)
		matrix[y][x]++
	}

	util.PrintMatrixTransform(true, matrix, func(i int) string {
		var str string
		if i == 0 {
			str = "."
		} else {
			str = "X"
		}
		return fmt.Sprintf("%s", str)
	})
}
