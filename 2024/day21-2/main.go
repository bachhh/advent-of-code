package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"

	"aoc2024/util"
)

type NodeStr util.TreeNode[string]

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
		input := scanner.Text()
		arrowScore := foo(input, 25)
		fmt.Println(input, arrowScore, inputScore(input))
		total += inputScore(input) * arrowScore
	}
	fmt.Println(total)
}

func foo(numeric string, depth int) int64 {
	var minScore int64 = math.MaxInt64
	arrows := NumToArrow('A', numeric)
	for _, arrow := range arrows {
		score := SolveRecurse(arrow, depth)
		if score < minScore {
			minScore = score
		}
	}

	return minScore
}

func SolveRecurse(input string, depth int) int64 {
	if depth == 0 {
		return int64(len(input))
	}
	if score, found := getCache(input, depth); found {
		return score
	}

	sumScoreOfSegments := int64(0)
	segments := splitSegment(input)
	for _, seg := range segments {
		moves := ArrowToArrow(seg)
		minScoreOfSegment := int64(math.MaxInt64)
		for _, move := range moves {
			score := SolveRecurse(move, depth-1)
			if score < minScoreOfSegment {
				minScoreOfSegment = score
			}
			util.Debugln(depth, segments, score, score)
		}
		sumScoreOfSegments += minScoreOfSegment
	}

	addCache(input, depth, sumScoreOfSegments)
	return sumScoreOfSegments
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

var AtAManualCache = map[string][]string{}

func ArrowToArrow(input string) []string {
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

func ScoreCompact(compact map[string]int) int64 {
	t := int64(0)
	for k, v := range compact {
		t += int64(len(k) * v)
	}
	return t
}

var CompactArrowCache = map[string]map[string]int{}

func splitSegment(input string) []string {
	result := []string{}
	for {
		before, after, found := strings.Cut(input, "A")
		if !found {
			return result
		}
		before, input = before+"A", after

		result = append(result, before)
	}
}

var solveRecurseCache = map[string]map[int]int64{}

func addCache(str string, depth int, score int64) {
	if solveRecurseCache[str] == nil {
		solveRecurseCache[str] = map[int]int64{}
	}
	solveRecurseCache[str][depth] = score
}

func getCache(str string, depth int) (int64, bool) {
	if solveRecurseCache[str] == nil {
		return 0, false
	}
	score, found := solveRecurseCache[str][depth]
	return score, found
}
