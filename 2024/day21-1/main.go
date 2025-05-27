package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
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
		util.Debugln(input)
		first := NumericToArrowConvert(input)
		util.Debugln(first[0])
		second := []string{}
		second = append(second, ArrowToArrowConvert(first[0])...)
		util.Debugln(second[0])

		third := []string{}
		third = append(third, ArrowToArrowConvert(second[0])...)
		util.Debugln(third[0])
		sort.Slice(third, func(i, j int) bool {
			return len(third[i]) < len(third[j])
		})

		output := third[0]

		fmt.Println(
			inputScore(input), len(output),
			input, ":", third,
		)
		total += inputScore(input) * len(third)
	}

	fmt.Println(total)
}

func NumericToArrowConvert(input string) []string {
	cur := byte('A')
	result := ""
	for _, b := range input {
		move1, _ := util.ManhattanPathToArrow(
			numericPos[cur],
			numericPos[byte(b)],
		)
		result = result + move1 + "A"
		cur = byte(b)
	}
	return []string{result}
}

func ArrowToArrowConvert(input string) []string {
	cur := byte('A')
	result := []string{""}
	for _, b := range input {
		next := []string{}
		move1, move2 := util.ManhattanPathToArrow(
			arrowPos[cur],
			arrowPos[byte(b)],
		)
		move1 += "A"
		move2 += "A"
		cur = byte(b)
		for i := range result {
			next = append(next, result[i]+move1)
			// next = append(next, result[i]+move2)
		}
		result = next
		// util.Debugln(result)
	}
	// util.Debugln(result)
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
