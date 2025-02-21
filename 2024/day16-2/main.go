package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

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

	matrix := util.ScanMatrix(scanner)
	sRow, sCol := util.FindMatrix(matrix, func(b byte) bool {
		return b == 'S'
	})

	eRow, eCol := util.FindMatrix(matrix, func(b byte) bool {
		return b == 'E'
	})

	allScore := map[string]*Step{}
	queue := util.NewQueue[Step]()

	path := map[string][]Step{}

	setLowestScore(allScore, sRow, sCol, util.East, 0)
	queue.Push(Step{Pair: Pair{Row: sRow, Col: sCol}, Dir: util.East, Score: 0})

	for !queue.IsEmpty() {
		cur, _ := queue.Pop()
		// cl := util.CloneMatrix(matrix)
		// cl[cur.Pair.Row][cur.Pair.Col] = '+'
		// util.PrintMatrix(cl)

		// util.PrintMatrixTransform(matrix, func)
		for _, nextDir := range util.ManhattanDirs {
			var isFinal bool
			next := util.Walk(cur.Pair, nextDir)
			if matrix[next.Row][next.Col] == '#' {
				continue
			}
			var turningScore int
			if nextDir != cur.Dir {
				turningScore = getTurningScore(cur.Dir, nextDir)
			}

			// score of next cell = current cell score + 1 ( move ) + turning score ( if have to turn )
			nextScore := turningScore + cur.Score + 1

			// if we reach end-cell no need to care about direction
			if next.Row == eRow && next.Col == eCol {
				isFinal = true
				nextDir = util.East // always turn east when reaches final
			}
			nextStep := Step{
				Pair: next,
				Dir:  nextDir,
			}

			// new cell, just insert
			nextLowestScoreWithStep := getLowestScore(allScore, next.Row, next.Col, nextDir)
			if nextLowestScoreWithStep == nil {
				if !isFinal {
					queue.Push(Step{
						Pair:  next,
						Dir:   nextDir,
						Score: nextScore,
					})
				}
				setLowestScore(allScore, next.Row, next.Col, nextDir, nextScore)

				path[nextStep.PosString()] = append(path[nextStep.PosString()], cur)
				continue
			}

			// check with lowest score so far
			nextLowestScore := nextLowestScoreWithStep.Score
			if nextScore == nextLowestScore {
				// add path
				path[nextStep.PosString()] = append(path[nextStep.PosString()], cur)
				continue
			}

			if nextScore < nextLowestScore {
				if !isFinal {
					queue.Push(Step{
						Pair:  next,
						Dir:   nextDir,
						Score: nextScore,
					})
				}
				setLowestScore(allScore, next.Row, next.Col, nextDir, nextScore)
				path[nextStep.PosString()] = []Step{cur}
				continue
			}
		}
	}
	// fmt.Printf("%+v\n", path)

	q2 := util.NewQueue[Step]()
	q2.Push(Step{Pair: Pair{Row: eRow, Col: eCol}, Dir: util.East})
	const mark = 'O'

	for !q2.IsEmpty() {
		cur, _ := q2.Pop()
		// fmt.Println(cur.PosString())
		// fmt.Println(path[cur.PosString()])
		if prevs, ok := path[cur.PosString()]; ok {
			// fmt.Println("_____1______")
			for _, prev := range prevs {
				// fmt.Println("_____2______")
				q2.Push(prev)
				matrix[prev.Pair.Row][prev.Pair.Col] = mark
			}
		}
	}
	total := 0

	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == mark {
				total++
			}
		}
	}
	util.PrintMatrix(matrix)
	fmt.Println(total + 1)

	// fmt.Println(getLowestScore(allScore, eRow, eCol))
}

func getTurningScore(before, after util.Direction) int {
	clockWise := util.Abs(int(after) - int(before))
	antiClockWise := 4 - util.Abs(int(before)-int(after))

	return util.Min(clockWise, antiClockWise) * 1000
}

func (s *Step) PosString() string {
	return fmt.Sprintf("%d-%d-%s", s.Pair.Row, s.Pair.Col, s.Dir)
}

type Step struct {
	Pair  Pair
	Dir   util.Direction
	Score int
}

func getLowestScore(m map[string]*Step, row, col int, dir util.Direction) *Step {
	key := fmt.Sprintf("%d-%d-%s", row, col, dir)
	if val, ok := m[key]; ok {
		return val
	}
	return nil
}

// func addPathToLowestScore(m map[string]*Step, row, col int, prev Pair, dir util.Direction) {
// 	key := fmt.Sprintf("%d-%d", row, col)
// 	if entry, ok := m[key]; ok {
// 		entry.Prev = append(entry.Prev, prev)
// 	}
// }

func setLowestScore(m map[string]*Step, row, col int, dir util.Direction, score int) {
	key := fmt.Sprintf("%d-%d-%s", row, col, dir)
	add := &Step{Score: score}
	m[key] = add
}
