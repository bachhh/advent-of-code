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

	allScore := map[string]int{}
	queue := util.NewQueue[Step]()

	setLowestScore(allScore, sRow, sCol, 0)
	queue.Push(Step{Pair: Pair{Row: sRow, Col: sCol}, Dir: util.East, Score: 0})

	for !queue.IsEmpty() {
		cur, _ := queue.Pop()
		// cl := util.CloneMatrix(matrix)
		// cl[cur.Pair.Row][cur.Pair.Col] = '+'
		// util.PrintMatrix(cl)

		// util.PrintMatrixTransform(matrix, func)
		for _, nextDir := range util.ManhattanDirs {
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
			nextLowestScore := getLowestScore(allScore, next.Row, next.Col)

			// found a lower score to reach the next cell
			// only add to queue if we found a lower score
			if nextLowestScore == -1 || nextScore < nextLowestScore {
				setLowestScore(allScore, next.Row, next.Col, nextScore)
				queue.Push(Step{
					Pair:  next,
					Dir:   nextDir,
					Score: nextScore,
				})
			}
		}
	}
	fmt.Println(getLowestScore(allScore, eRow, eCol))
}

func getTurningScore(before, after util.Direction) int {
	clockWise := util.Abs(int(after) - int(before))
	antiClockWise := 4 - util.Abs(int(before)-int(after))

	return min(clockWise, antiClockWise) * 1000
}

type Step struct {
	Pair  Pair
	Dir   util.Direction
	Score int
}

func getLowestScore(m map[string]int, row, col int) int {
	key := fmt.Sprintf("%d-%d", row, col)
	if val, ok := m[key]; ok {
		return val
	}
	return -1
}

func setLowestScore(m map[string]int, row, col int, score int) {
	key := fmt.Sprintf("%d-%d", row, col)
	m[key] = score
}
