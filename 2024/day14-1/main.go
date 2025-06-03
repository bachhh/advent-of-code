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

	quardrantCount := make([]int, 4)
	maxX := 101
	maxY := 103

	quadrants := calcQuadrants(maxX, maxY)
	matrix := make([][]int, maxY)
	for i := range matrix {
		matrix[i] = make([]int, maxX)
		for j := range matrix[i] {
			matrix[i][j] = 0
		}
	}

	for scanner.Scan() {
		x, y, dx, dy, err := ParseLine(scanner.Text())
		if err != nil {
			panic(err)
		}
		// fmt.Println(x, y, dx, dy)
		finalPosX := (x + 100*dx)
		finalPosY := (y + 100*dy)

		// wrap around if position is negative
		finalPosX += (util.Abs(finalPosX/maxX) + 1) * maxX
		finalPosY += (util.Abs(finalPosY/maxY) + 1) * maxY

		finalPosX = finalPosX % maxX
		finalPosY = finalPosY % maxY

		matrix[finalPosY][finalPosX]++
		for i, quad := range quadrants {
			if (quad[0] <= finalPosY && finalPosY < quad[1]) && (quad[2] <= finalPosX && finalPosX < quad[3]) {
				quardrantCount[i]++
			}
		}

		matrix[finalPosY][finalPosX]++
	}

	total := 1
	for _, count := range quardrantCount {
		fmt.Println("count", count)
		total *= count
	}
	fmt.Println("factor", total)
	// q2 := []int{52, 101, 51, 101}
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
