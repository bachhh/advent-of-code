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
	for _, char := range input {
		b := byte(char)
		next := []string{}

		moves := move(
			arrowPos[cur],
			arrowPos[b],
			Pair{Row: 0, Col: 0},
		)

		cur = b
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
	input := "^>A"

	node := buildTreeOfArrow(input, 3)
	// fmt.Println(input)
	node.PrintTree()
	fmt.Println(optimalArrow(input))
}

func buildTreeOfArrow(input string, depth int) *util.TreeNode[string] {
	root := &util.TreeNode[string]{Value: input}
	curDepth := []*util.TreeNode[string]{root}

	for i := 0; i < depth; i++ {
		nextDepth := []*util.TreeNode[string]{}

		for _, curNode := range curDepth {
			moves := ArrowToArrow(curNode.Value)
			for _, move := range moves {
				newChild := &util.TreeNode[string]{Value: move, Parent: curNode}
				curNode.Children = append(curNode.Children, newChild)

				nextDepth = append(nextDepth, newChild)
			}
		}
		curDepth = nextDepth
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

	slices.Sort(result)
	slices.Compact(result)
	result = slices.DeleteFunc(result, func(in string) bool {
		return in == ""
	})
	if len(result) == 0 {
		return []string{""}
	}

	return result
}

// var scoreCacheDepth2 = map[string]int{}

// for any given arrow input
// return the next move sequence that leads to a shortest move sequence at the next 2 depth
// also return the move at next 2 depths
func optimalArrow(input string) (string, int) {
	depth := 3
	root := buildTreeOfArrow(input, depth)
	root.PrintTree()
	qu := util.NewQueue[*util.TreeNode[string]]()
	qu.Push(root)
	for !qu.IsEmpty() {
		cur, _ := qu.Peek()
		if len(cur.Children) == 0 {
			break
		}

		cur, _ = qu.Pop()
		if cur == nil {
			panic("nil cur")
		}

		for _, child := range cur.Children {
			if child != nil {
				qu.Push(child)
			}
		}
	}
	leafNodeSlice := qu.ToSlice()
	slices.SortFunc(leafNodeSlice, func(a, b *util.TreeNode[string]) int {
		return cmp.Compare(len(a.Value), len(b.Value))
	})

	optimalScore := len(leafNodeSlice[0].Value)
	curNode := leafNodeSlice[0]
	for i := 1; i < depth; i++ {
		curNode = curNode.Parent
	}

	return curNode.Value, optimalScore
}
