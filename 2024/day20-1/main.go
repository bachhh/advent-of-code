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

	inputFile := flag.Arg(0)
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	maze := [][]byte{}
	start, end := Pair{}, Pair{}
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		maze = append(maze, []byte(line))
		for j := range line {
			if line[j] == 'S' {
				start = Pair{
					Row: i,
					Col: j,
				}
			} else if line[j] == 'E' {
				end = Pair{
					Row: i,
					Col: j,
				}
			}
		}
	}
	Solve(start, end, maze)
}

func Solve(start, end Pair, maze [][]byte) {
	// walk the matrix
	// for any 2 points on the path, check if it's possible to jump to another
	// point just with breaking 2 walls

	dist := map[Pair]int{}
	dist[start] = 0
	path := []Pair{start}
	cur := start
	for cur != end {
		for _, dir := range util.ManhattanDirs {
			next := util.Walk(cur, dir)
			if _, found := dist[next]; found {
				continue
			}
			if maze[next.Row][next.Col] != '#' {
				dist[next] = dist[cur] + 1
				cur = next
				path = append(path, next)
				break
			}
		}
	}
	// fmt.Printf("%+v\n", dist)
	// fmt.Printf("%+v\n", path)

	distSaved100Count := 0
	allDistSaved := map[int]int{}
	for _, cur := range path {
		for _, dir := range util.ManhattanDirs {
			next1 := util.Walk(cur, dir)
			if maze[next1.Row][next1.Col] == '#' {
				for _, dir := range util.ManhattanDirs {
					next2 := util.Walk(next1, dir)
					if !util.IsPairInbound(next2, maze) ||
						next2 == cur ||
						dist[next2] < dist[cur] ||
						false {
						continue
					}

					if maze[next2.Row][next2.Col] != '#' {
						distSaved := util.Abs(dist[cur]-dist[next2]) - 2
						if distSaved > 0 {
							if distSaved >= 100 {
								distSaved100Count++
							}
							allDistSaved[distSaved]++
							fmt.Printf("saved %d dist by skipping from %+v to %+v\n",
								distSaved,
								cur,
								next2,
							)
						}
					}
				}
			}
		}
	}
	fmt.Println(distSaved100Count)
}
