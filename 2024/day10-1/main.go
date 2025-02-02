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

	"aoc2024/util"
)

var allowedDir = []util.Direction{util.North, util.East, util.South, util.West}

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

	var matrix [][]int
	trailHeads := []util.Pair{}
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Bytes()
		cp := make([]int, len(line))
		for col := range line {
			num, err := util.CharToInt(line[col])
			if err != nil {
				num = 0
			}
			cp[col] = num

			if line[col] == '0' {
				trailHeads = append(trailHeads, util.Pair{Row: row, Col: col})
			}
		}
		matrix = append(matrix, cp)
	}

	total := 0
	for _, pos := range trailHeads {
		score := eval(matrix, pos, []util.Pair{})
		fmt.Println("pos ", pos, " reachable ", score)
		total += len(score)
	}
	fmt.Println(total)
}

func eval(matrix [][]int, pos util.Pair, visited []util.Pair) []util.Pair {
	if matrix[pos.Row][pos.Col] == 9 {
		return []util.Pair{pos}
	}
	reachable := []util.Pair{}
	for _, dir := range allowedDir {
		nextPos := util.Walk(pos, dir)
		if !util.IsPairInbound(nextPos, matrix) {
			continue
		}
		if isVisited(nextPos, visited) {
			continue
		}
		if matrix[nextPos.Row][nextPos.Col] != matrix[pos.Row][pos.Col]+1 {
			continue
		}
		visitedClone := util.CloneSlice(visited)

		fromChild := eval(matrix, nextPos, append(visitedClone, nextPos))
		reachable = append(reachable, fromChild...)
	}
	// fmt.Println(reachable)

	slices.SortFunc(reachable, pairCmp)
	return slices.CompactFunc(reachable, func(a, b util.Pair) bool {
		return a.Row == b.Row && b.Col == a.Col
	})
}

func isVisited(point util.Pair, visited []util.Pair) bool {
	for _, p := range visited {
		if p.Row == point.Row && p.Col == point.Col {
			return true
		}
	}
	return false
}

func pairCmp(a, b util.Pair) int {
	if a.Row == b.Row {
		return cmp.Compare(a.Col, b.Col)
	}
	return cmp.Compare(a.Row, b.Row)
}
