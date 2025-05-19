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

var size = 0

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
	scanner := bufio.NewScanner(file)

	fallingBytes := []Pair{}

	scanner.Scan()
	size, err = strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic("string must be 2 digit separated by a comma")
		}

		// dist from left edge -> col
		col, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			panic(err)
		}

		// dist from top edge -> row
		row, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			panic(err)
		}

		fallingBytes = append(fallingBytes, Pair{Row: row, Col: col})
	}

	for i := 1; i <= len(fallingBytes); i++ {
		if test(fallingBytes[:i]) {
			fmt.Printf("maze is blocked at %d %+v\n", i, fallingBytes[i])
			break
		} else {
			fmt.Printf("maze is still reachable at %d %+v\n", i, fallingBytes[i])
		}
	}
}

func test(fallingBytes []Pair) bool {
	maze := util.NewMatrix[Marker](size, size)
	for i := range maze {
		for j := range maze[i] {
			maze[i][j] = Marker{
				IsBlocked: false,
				ReachTime: math.MaxInt,
			}
		}
	}

	maze[0][0].ReachTime = 0
	for i := range fallingBytes {
		row, col := fallingBytes[i].Row, fallingBytes[i].Col
		maze[row][col].IsBlocked = true
	}

	qu := util.NewQueue[Move]()
	qu.Push(Move{Pair: Pair{Row: 0, Col: 0}, Time: 0})
	for !qu.IsEmpty() {
		cur, _ := qu.Pop()
		row, col := cur.Pair.Row, cur.Pair.Col

		for _, dir := range util.ManhattanDirs {
			next := util.Walk(Pair{Row: row, Col: col}, dir)

			if !util.IsPairInbound(next, maze) {
				continue
			}

			nextTime := cur.Time + 1

			// cell is reachable in less than what we have
			if maze[next.Row][next.Col].ReachTime <= nextTime {
				continue
			}

			// byte fallen, cannot move
			if maze[next.Row][next.Col].IsBlocked {
				continue
			}

			if next.Row == size-1 && next.Col == size-1 {
				fmt.Println("dest reached", nextTime)
				return false
			}

			qu.Push(Move{Pair: next, Time: nextTime})
			maze[next.Row][next.Col].ReachTime = nextTime
		}
	}
	return true
}

type Marker struct {
	IsBlocked bool
	ReachTime int
}

type Move struct {
	Pair Pair
	Time int
}
