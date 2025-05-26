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

	minDistSaved := 100
	const cheatLimit = 20
	allDistSaved := map[int]int{}
	count := 0
	for i, cur := range path {
		if cur == end {
			continue
		}
		for _, next := range path[i+1:] {
			if next != cur {
				step := util.ManhattanDistance(cur, next)
				if step <= cheatLimit {
					distSaved := util.Abs(dist[cur]-dist[next]) - step
					if distSaved >= minDistSaved {
						allDistSaved[distSaved]++
						count++
					}
				}
			}
		}
	}
	fmt.Println(allDistSaved)
	fmt.Println(count)
}
