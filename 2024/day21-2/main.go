package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"

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
/*
    +---+---+
    | ^ | A |
+---+---+---+
| < | v | > |
+---+---+---+
*/
type Pair = util.Pair

// numeric -> arrow(first) -> arrow x N-1 -> final

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
	total := int64(0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// input := scanner.Text()
		// result := foo(input)
		// minLenStr := minLen(result)

		// arrowScore := SolveCompact(input, 27)
		// fmt.Println(
		// 	inputScore(input), arrowScore,
		// 	input, ":",
		// )
		// total += inputScore(input) * arrowScore
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
				next = append(next, result[i]+moves[j])
			}
		}
		result = next
		cur = byte(b)
	}

	return result
}

func AtACache(input string) []string {
	result := []string{""}
	for {
		before, after, found := strings.Cut(input, "A")
		if !found {
			break
		}
		before, input = before+"A", after
		next := []string{}
		// fmt.Println("input", input)

		moves, found := manualCache[before]
		if !found {
			panic("manualCache[before] " + before)
		}
		for i := range result {
			for j := range moves {
				next = append(next, result[i]+moves[j])
			}
		}
		result = next

	}
	return result
}

var AtAManualCache = map[string][]string{}

func AtAManual(input string) []string {
	if entry, found := AtAManualCache[input]; found {
		return entry
	}
	cur := byte('A')
	result := []string{""}
	for _, char := range input {
		b := byte(char)
		next := []string{}

		moves := MoveArrow(cur, b)
		cur = b
		for i := range result {
			for j := range moves {
				next = append(next, result[i]+moves[j])
			}
		}
		result = next
	}

	AtAManualCache[input] = result
	return result
}

// ['1']['9'] -> ">>^^", 1 level
// ['1']['9'] -> ">>^^" -> ">AA<^AA" 2 level
var numToArrowCache = map[byte]map[byte][]string{}

var (
	manualCache = map[string][]string{
		"A": {"A"},
	}
	optimalCache = map[string]string{}
)

func init() {
	for a1 := range arrowPos {
		for a2 := range arrowPos {
			if a1 == a2 {
				continue
			}
			moves := MoveArrow(a1, a2)
			for _, move := range moves {
				manualCache[move] = AtAManual(move)
			}
		}
	}
	// there are many solution for moving from A-B,
	fmt.Println(manualCache)

	// init cache for numeric
	for pos1 := range numericPos {
		if numToArrowCache[pos1] == nil {
			numToArrowCache[pos1] = map[byte][]string{}
		}

		for pos2 := range numericPos {
			if pos1 == pos2 {
				continue
			}
			moves := MoveNum(pos1, pos2)
			for _, move := range moves {
				numToArrowCache[pos1][pos2] = slices.Concat(
					numToArrowCache[pos1][pos2],
					AtAManual(move))
			}
		}
	}

	// now init cache for arrow position
	// we calculate the move from arrow1 to arrow2 with the lowest score
	// the score being
	// for arrow1 := range arrowPos {
	// 	for arrow2 := range arrowPos {
	// 		if arrow1 == arrow2 {
	// 			continue
	// 		}
	// 		moves := MoveArrow(arrow1, arrow2)
	// 		for _, move := range moves {
	// 			// fmt.Println(string(arrow1), string(arrow2), move)
	// 			// time.Sleep(time.Second)
	// 			// return
	// 			// calculate the lowest score of this move
	// 		}
	// 	}
	// }

	fmt.Println(lowestScore2("v<A"))
}

func lowestScore2(input string) (string, int) {
	type data map[string]int
	root := &util.TreeNode[data]{Value: data{input: 1}}

	qu := util.NewQueue[*util.TreeNode[data]]()
	qu.Push(root)

	for {
		allNode := qu.PopAll()
		for _, node := range allNode {
			moves := AtACompact(node.Value)
			fmt.Println(node.Value, moves)

			for _, move := range moves {
				newChild := node.AddChild(move)
				qu.Push(newChild)

			}
		}
		time.Sleep(time.Second)
		root.PrintTree(10)
	}
}

func lowestScore(input string) (string, int) {
	root := &util.TreeNode[string]{Value: input}
	qu := util.NewQueue[*util.TreeNode[string]]()
	qu.Push(root)

	for depth := 1; ; depth++ {
		// evaluate the next depth from the current depth
		allNode := qu.PopAll()
		// time.Sleep(time.Second)

		for _, node := range allNode {
			// fmt.Println("len", len(node.Value))
			moves := AtACache(node.Value)
			// fmt.Println("moveCount", len(moves))
			for _, move := range moves {
				newChild := node.AddChild(move)
				qu.Push(newChild)
			}
		}
		fmt.Println("evaluate", depth)

		// re-evaluate if we can decide which children is the optimal choice
		childrenScore := util.MapFunc(root.Children, func(node *util.TreeNode[string]) int {
			leaf := node.FindSmallestLeaf(func(a, b string) int {
				return cmp.Compare(len(a), len(b))
			})
			return len(leaf.Value)
		})

		smallestScore, smallestChild := math.MaxInt, (*util.TreeNode[string])(nil)
		found := false
		for i, score := range childrenScore {
			if smallestScore != math.MaxInt && score != smallestScore {
				found = true
			}
			if score < smallestScore {
				smallestScore, smallestChild = score, root.Children[i]
			}
		}
		if found || len(root.Children) < 2 {
			return smallestChild.Value, smallestScore
		}
		root.PrintTree(2)
	}
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

func inputScore(input string) int64 {
	i, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		panic(err)
	}
	return int64(i)
}

func MoveNum(pos1 byte, pos2 byte) []string {
	return move(
		numericPos[pos1],
		numericPos[pos2],
		Pair{Row: 3, Col: 0},
	)
}

func MoveArrow(pos1 byte, pos2 byte) []string {
	return move(
		arrowPos[pos1],
		arrowPos[pos2],
		Pair{Row: 0, Col: 0},
	)
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
		return []string{"A"}
	}

	result = util.MapFunc(result, func(in string) string {
		return in + "A"
	})

	return result
}

func AtACompact(input map[string]int) []map[string]int {
	fmt.Printf("%+v\n", input)
	result := []map[string]int{}
	for str, count := range input {
		moves := AtAManual(str)
		next := []map[string]int{}
		for _, move := range moves {
			moveMap := multiplyMap(CompactArrow(move), count)
			for _, cur := range result {
				next = append(next, mergeMap(moveMap, cur))
			}
		}
		result = next
	}
	return result
}

func ScoreCompact(compact map[string]int) int64 {
	t := int64(0)
	for k, v := range compact {
		t += int64(len(k) * v)
	}
	return t
}

var CompactArrowCache = map[string]map[string]int{}

// turns an arrow pattern into compact form
func CompactArrow(input string) map[string]int {
	if entry, found := CompactArrowCache[input]; found {
		return entry
	}
	m := map[string]int{}

	for {
		before, after, found := strings.Cut(input, "A")
		if !found {
			break
		}
		pattern := before + "A"
		m[pattern]++
		input = after
	}
	CompactArrowCache[input] = m
	return m
}

func multiplyMap(a map[string]int, n int) map[string]int {
	m := map[string]int{}
	for k, v := range a {
		m[k] = v * n
	}
	return m
}

func mergeMap(a, b map[string]int) map[string]int {
	result := make(map[string]int)

	// Copy all elements from map a
	for key, value := range a {
		result[key] = value
	}

	// Add or update elements from map b
	for key, value := range b {
		result[key] += value
	}

	return result
}
