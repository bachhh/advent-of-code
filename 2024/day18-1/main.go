package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"aoc2024/util"
)

type Pair = util.Pair

const (
	size          = 71
	byteFallCount = 1024
)

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

	maze := util.NewMatrix[Marker](size, size)

	for i := range maze {
		for j := range maze[i] {
			maze[i][j] = Marker{FallTime: math.MaxInt, ReachTime: math.MaxInt}
		}
	}
	maze[0][0].ReachTime = 0

	for i := 1; scanner.Scan() && i <= byteFallCount; i++ {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		// dist from left edge -> col
		col, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		// dist from top edge -> row
		row, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			continue
		}
		maze[row][col].FallTime = i
	}

	qu := util.NewQueue[Move]()
	qu.Push(Move{Pair: Pair{Row: 0, Col: 0}, Time: 0})

	for !qu.IsEmpty() {
		cur, _ := qu.Pop()
		row, col := cur.Pair.Row, cur.Pair.Col
		fmt.Printf("%+v\n", cur)

		for _, dir := range util.ManhattanDirs {
			next := util.Walk(Pair{Row: row, Col: col}, dir)

			if !util.IsPairInbound(next, maze) {
				continue
			}

			nextTime := cur.Time + 1

			// cell is reachable in less than what we have
			if maze[next.Row][next.Col].ReachTime <= nextTime {
				// fmt.Printf("next %+v is reachable in less time %d than current time %d\n", next, maze[next.Row][next.Col].ReachTime, nextTime)
				continue
			}

			// byte fallen, cannot move
			if maze[next.Row][next.Col].FallTime < math.MaxInt {
				continue
			}

			if next.Row == size-1 && next.Col == size-1 {
				fmt.Println("dest reached", nextTime)
				util.PrintMatrixTransform(false, maze, func(base Marker) string {
					if base.FallTime < math.MaxInt {
						return "#"
					}
					if base.ReachTime < math.MaxInt {
						return "O"
					}
					return "."
				})
				return
			}

			qu.Push(Move{Pair: next, Time: nextTime})
			maze[next.Row][next.Col].ReachTime = nextTime
			fmt.Println("next dest ", next.Row, "-", next.Col, " time ", nextTime)
		}
	}
	fmt.Println("destination unreachable")
}

type Marker struct {
	FallTime  int
	ReachTime int
}

type (
	Move struct {
		Pair Pair
		Time int
	}
)
