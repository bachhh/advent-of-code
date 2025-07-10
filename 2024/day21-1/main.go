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
}

// movement part
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

func ArrowToArrow(input string) []string {
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
			for j := range moves {
				next = append(next, result[i]+moves[j]+"A")
			}
		}
		result = next
	}
	return result
}

func init() {
	node := ArrowToArrowTree("^^^<A")
	node.PrintTree(1)
}

func ArrowToArrowTree(input string) *util.TreeNode[string] {
	cur := byte('A')
	root := &util.TreeNode[string]{Value: input}
	curNodes := []*util.TreeNode[string]{root}

	for _, b := range input {
		moves := move(
			arrowPos[cur],
			arrowPos[byte(b)],
			Pair{Row: 0, Col: 0},
		)
		cur = byte(b)

		next := []*util.TreeNode[string]{}

		for i := range curNodes {
			for j := range moves {
				child := &util.TreeNode[string]{Value: curNodes[i].Value + moves[j] + "A"}
				curNodes[i].Children = append(curNodes[i].Children, child)
				next = append(next, child)
			}
		}
		curNodes = next
	}
	return root
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

func unfold(input []string, transformer func(string) []string) []string {
	result := []string{}
	for i := range input {
		result = append(result, transformer(input[i])...)
	}
	return result
}

var scoreCacheDepth2 = map[string]int{}

func bestMoveArrow(move string, depth int) []string {
	second := []string{move}
	for i := 0; i < depth; i++ {
		second = unfold(second,
			func(input string) []string {
				return ArrowToArrow(input)
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

	third := unfold(second, func(input string) []string {
		return ArrowToArrow(input)
	})
	slices.SortFunc(second, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})
	return third
}
