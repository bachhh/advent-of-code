package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"aoc2024/util"
)

/*
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//     | 0 | A |
//     +---+---+
*/
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

	inputFile := flag.Arg(0)
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		first := NumericToArrowConvert(input)
		second := ArrowToArrowConvert(first)
		third := ArrowToArrowConvert(second)

		fmt.Println(
			inputScore(input), len(third),

			input, ":", third)
		total += inputScore(input) * len(third)
	}
	fmt.Println(total)
}

func NumericToArrowConvert(numeric string) string {
	cur := byte('A')
	result := ""
	for _, b := range numeric {
		move1, _ := util.ManhattanPathToArrow(
			numericPos[cur],
			numericPos[byte(b)],
		)
		move1 += "A"
		cur = byte(b)
		result += move1
	}
	return result
}

func ArrowToArrowConvert(numeric string) string {
	cur := byte('A')
	result := ""
	for _, b := range numeric {
		move1, _ := util.ManhattanPathToArrow(
			arrowPos[cur],
			arrowPos[byte(b)],
		)
		move1 += "A"
		// fmt.Println(move1, "\t", string(cur))
		cur = byte(b)
		result += move1
	}
	return result
}

// convert from current pos -> arrow sequence
var arrowPos = map[byte]Pair{
	'A': {Row: 0, Col: 2},
	'^': {Row: 0, Col: 1},
	'v': {Row: 1, Col: 1},
	'<': {Row: 1, Col: 0},
	'>': {Row: 1, Col: 2},
}

// convert from current pos -> arrow sequence
var numericPos = map[byte]Pair{
	'0': {Row: 3, Col: 1},
	'1': {Row: 2, Col: 0},
	'2': {Row: 2, Col: 1},
	'3': {Row: 2, Col: 2},
	'4': {Row: 1, Col: 0},
	'5': {Row: 1, Col: 1},
	'6': {Row: 1, Col: 2},
	'7': {Row: 0, Col: 0},
	'8': {Row: 0, Col: 1},
	'9': {Row: 0, Col: 2},
	'A': {Row: 3, Col: 2},
}

func inputScore(input string) int {
	i, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		panic(err)
	}
	return i
}
