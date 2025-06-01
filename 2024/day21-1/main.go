package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"

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
		// input := scanner.Text()
		// result := foo(input)
		// minLenStr := minLen(result)

		// fmt.Println(
		// 	inputScore(input), len(minLenStr),
		// 	input, ":",
		// )
		// total += inputScore(input) * len(minLenStr)
	}
	fmt.Println(total)
	initCache()
}

func foo(input string) []string {
	first := NumToArrow('A', input)
	second := make([]string, len(first))
	copy(second, first)
	fmt.Println("len first", len(first), "len str first", len(first[0]))
	fmt.Println(first[0])
	for i := 0; i < 2; i++ {
		second = transform(second,
			func(input string) []string {
				return ArrowToArrow(input, 1)
			})
		slices.SortFunc(second, func(a, b string) int {
			return cmp.Compare(len(a), len(b))
		})

		minLen := len(second[0])
		for j := range second {
			if len(second[j]) > minLen {
				second = second[:j]
				break
			}
		}
		fmt.Println("len_second", len(second), "lenstrsecond", len(second[0]))
		fmt.Println(second[0])
	}

	third := transform(second, func(input string) []string {
		return ArrowToArrow(input, 0)
	})
	fmt.Println("len third", len(third), "len string third", len(third[0]))
	slices.SortFunc(second, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})
	fmt.Println(third[0])
	return third[:1]
}

func NumToArrow(cur byte, input string) []string {
	result := []string{""}
	for _, b := range input {
		next := []string{}
		moves := move(
			numericPos[cur],
			numericPos[byte(b)],
			Pair{Row: 3, Col: 0},
		)
		for i := range result {
			for j := range moves {
				next = append(next, result[i]+moves[j]+"A")
			}
		}
		result = next
		cur = byte(b)
	}

	return result
}

var arrowToArrowCache = map[string]string{}

func ArrowToArrow(
	input string,
	mode int, // 0 = single, 1 = multi, use all cases, 2= optimized: only pick the best score
) []string {
	result := []string{""}
	for input != "" {
		idx := strings.IndexByte(input, byte('A'))
		section := input[:idx+1]
		// do stuff
		input = input[idx+1:]
		if cache, found := arrowToArrowCache[section]; found {
			result = util.MapFunc(result, func(in string) string {
				return in + cache
			})
		}
		output := arrowToArrow(section, 0)
		arrowToArrowCache[section] = actualWork
	}
	return result
}

func arrowToArrow(
	input string,
	mode int, // 0 = single, 1 = multi, use all cases, 2= optimized: only pick the best score
) []string {
	cur := byte('A')
	result := []string{""}
	for _, b := range input {
		next := []string{}
		moves := move(
			arrowPos[cur],
			arrowPos[byte(b)],
			Pair{Row: 0, Col: 0},
		)
		cur = byte(b)
		for i := range result {
			switch mode {
			case 0:
				next = append(next, result[i]+moves[0]+"A")
			case 1:
				for j := range moves {
					moves[j] += "A"
				}
				for j := range moves {
					next = append(next, result[i]+moves[j])
				}
			case 2:
				panic(2)
			}
		}
		result = next
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

func move(from, to Pair, forbidden Pair) []string {
	rMove, cMove := "", ""
	rowDiff := to.Row - from.Row
	colDiff := to.Col - from.Col
	if rowDiff < 0 {
		rMove = strings.Repeat("^", util.Abs(rowDiff))
	} else if rowDiff > 0 {
		rMove = strings.Repeat("v", util.Abs(rowDiff))
	}
	if colDiff < 0 {
		cMove = strings.Repeat("<", util.Abs(colDiff))
	} else if colDiff > 0 {
		cMove = strings.Repeat(">", util.Abs(colDiff))
	}

	result := []string{}
	moveRowFirst := Pair{
		Row: from.Row + rowDiff,
		Col: from.Col,
	}
	if !util.IsCollinear(from, forbidden, moveRowFirst) {
		result = append(result, rMove+cMove)
	}
	moveColFirst := Pair{
		Row: from.Row,
		Col: from.Col + colDiff,
	}
	if !util.IsCollinear(from, forbidden, moveColFirst) {
		result = append(result, cMove+rMove)
	}
	return result
}

func minLen(input []string) string {
	str := ""
	for i := range input {
		if len(input[i]) < len(str) || str == "" {
			str = input[i]
		}
	}
	return str
}

// directions contains the possible move directions
var directions = []rune{'^', 'v', '<', '>'}

// gen recursively generates all move sequences up to given length
func genMoves(length int) []string {
	var results []string
	var dfs func(path []rune, depth int)
	dfs = func(path []rune, depth int) {
		if depth == length {
			results = append(results, string(path))
			return
		}
		for _, dir := range directions {
			dfs(append(path, dir), depth+1)
		}
	}
	dfs([]rune{}, 0)
	return results
}

func transform(input []string, transformer func(string) []string) []string {
	result := []string{}
	for i := range input {
		result = append(result, transformer(input[i])...)
	}
	return result
}

var scoreCacheDepth2 = map[string]int{}

func initCache() {
	allTest := []string{}
	for i := 1; i <= 3; i++ {
		allTest = append(allTest, genMoves(i)...)
	}
	allTest = append(allTest, "")
	for i := range allTest {
		allTest[i] += "A"
		output := bestMoveArrow(allTest[i], 3)
		// allTest[i], len(output[0]))
		scoreCacheDepth2[allTest[i]] = len(output[0])
		fmt.Println(allTest[i], len(output[0]))
	}
}

func bestMoveArrow(move string, depth int) []string {
	second := []string{move}
	for i := 0; i < depth; i++ {
		second = transform(second,
			func(input string) []string {
				return ArrowToArrow(input, 1)
			})
		slices.SortFunc(second, func(a, b string) int {
			return cmp.Compare(len(a), len(b))
		})

		minLen := len(second[0])
		for j := range second {
			if len(second[j]) > minLen {
				second = second[:j]
				break
			}
		}
	}

	third := transform(second, func(input string) []string {
		return ArrowToArrow(input, 0)
	})
	slices.SortFunc(second, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})
	return third
}
