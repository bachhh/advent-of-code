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

	visited := [][]bool{}
	matrix := [][]byte{}
	for scanner.Scan() {
		byte := scanner.Bytes()
		cp := util.CloneSlice(byte)
		matrix = append(matrix, cp)
		visited = append(visited, make([]bool, len(cp)))
	}
	util.PrintMatrix(matrix)

	total := 0
	queue := util.NewQueue[util.Pair]()
	localQueue := util.NewQueue[util.Pair]()
	queue.Push(util.Pair{Row: 0, Col: 0})
	for queue.Size() > 0 {
		head, _ := queue.Pop()

		// visited
		if visited[head.Row][head.Col] {
			continue
		}

		// util.PrintMatrix(matrix)
		localQueue.Push(head)
		area := 0
		peri := 0
		for localQueue.Size() > 0 {

			cur, _ := localQueue.Pop()
			if visited[cur.Row][cur.Col] {
				continue
			}

			area++
			for _, dir := range allDir {
				next := util.Walk(cur, dir)

				if next.Row < 0 || next.Col < 0 || next.Row >= len(matrix) || next.Col >= len(matrix[0]) {
					peri++
					continue
				}

				if matrix[next.Row][next.Col] != matrix[cur.Row][cur.Col] {
					// discover new field, push to global queue and skip
					queue.Push(next)
					peri++
					continue
				}

				if visited[cur.Row][cur.Col] {
					continue
				}
				localQueue.Push(next)
			}

			// mark as visited
			visited[cur.Row][cur.Col] = true
		}
		fmt.Println(area, peri, area*peri)
		total += (area * peri)
	}
	fmt.Println(total)
}

func bfs(matrix [][]byte, start util.Pair) int {
	queue := []util.Pair{start}
	area := 0
	peri := 0
	for len(queue) > 0 {
		head := queue[0]
		queue = queue[1:]
		area++
		// visited
		if matrix[head.Row][head.Col] == '.' {
			continue
		}
		for _, dir := range allDir {
			next := util.Walk(start, dir)
			if next.Row < 0 || next.Col < 0 || next.Row >= len(matrix) || next.Col >= len(matrix[0]) {
				peri++
			} else if matrix[next.Row][next.Col] != matrix[head.Row][head.Col] {
				peri++
			}
			queue = append(queue, next)
		}
	}
	return peri * area
}

var allDir = []util.Direction{
	util.North,
	util.East,
	util.South,
	util.West,
}
