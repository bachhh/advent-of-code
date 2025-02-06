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
	queue := util.NewQueue[Pair]()
	localQueue := util.NewQueue[Pair]()
	queue.Push(Pair{Row: 0, Col: 0})
	for queue.Size() > 0 {
		head, _ := queue.Pop()

		// visited
		if visited[head.Row][head.Col] {
			continue
		}

		// util.PrintMatrix(matrix)
		localQueue.Push(head)
		area := 0
		boundary := [][]Pair{}
		for localQueue.Size() > 0 {

			cur, _ := localQueue.Pop()
			if visited[cur.Row][cur.Col] {
				continue
			}

			area++
			for _, dir := range allDir {
				next := util.Walk(cur, dir)

				if next.Row < 0 || next.Col < 0 || next.Row >= len(matrix) || next.Col >= len(matrix[0]) {
					boundary = append(boundary, []Pair{cur, next})
					continue
				}

				if matrix[next.Row][next.Col] != matrix[cur.Row][cur.Col] {
					// discover new field, push to global queue and skip
					queue.Push(next)
					boundary = append(boundary, []Pair{cur, next})
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

		sideCount := 0
		for len(boundary) > 0 {
			curFenc := boundary[len(boundary)-1]
			boundary = boundary[:len(boundary)-1]
			isUpDown := (curFenc[1].Row - curFenc[0].Row) == 0

			fmt.Println(curFenc, isUpDown)
			fmt.Println("boundary", boundary)

			// if side direction is up-down, columns should always be equal,
			// if side direction is left-right, rows should always be equal,
			if isUpDown {
				neighbor, found := searchAndRemove(&boundary, func(e []util.Pair) bool {
					fmt.Println("search neighbor", curFenc, e)
					return e[0].Col == curFenc[0].Col &&
						e[1].Col == curFenc[1].Col &&
						util.Abs(e[0].Row-curFenc[0].Row) == 1 &&
						util.Abs(e[1].Row-curFenc[1].Row) == 1
				})
				if found {
					fmt.Println("found neighbor", curFenc, neighbor)
					var up, down []util.Pair
					if neighbor[0].Row < curFenc[0].Row {
						up, down = neighbor, curFenc
					} else {
						up, down = curFenc, neighbor
					}

					for {
						up, found = searchAndRemove(&boundary, func(e []util.Pair) bool {
							return e[0].Col == up[0].Col &&
								e[1].Col == up[1].Col &&
								e[0].Row-up[0].Row == -1
						})
						if !found {
							break
						}
						fmt.Println("up", up)
					}

					for {
						down, found = searchAndRemove(&boundary, func(e []util.Pair) bool {
							return e[0].Col == down[0].Col &&
								e[1].Col == down[1].Col &&
								e[0].Row-down[0].Row == 1
						})
						if !found {
							break
						}
						fmt.Println("down", down)
					}
				}
			} else {
				neighbor, found := searchAndRemove(&boundary, func(e []util.Pair) bool {
					fmt.Println("search neighbor", curFenc, e)
					return e[0].Row == curFenc[0].Row &&
						e[1].Row == curFenc[1].Row &&
						util.Abs(e[0].Col-curFenc[0].Col) == 1 &&
						util.Abs(e[1].Col-curFenc[1].Col) == 1
				})
				if found {
					fmt.Println("found neighbor", curFenc, neighbor)
					var left, right []util.Pair
					if neighbor[0].Col < curFenc[0].Col {
						left, right = neighbor, curFenc
					} else {
						left, right = curFenc, neighbor
					}

					for {
						left, found = searchAndRemove(&boundary, func(e []util.Pair) bool {
							return e[0].Row == left[0].Row &&
								e[1].Row == left[1].Row &&
								e[0].Col-left[0].Col == -1
						})
						if !found {
							break
						}
						fmt.Println("left", left)
					}

					for {
						right, found = searchAndRemove(&boundary, func(e []util.Pair) bool {
							return e[0].Row == right[0].Row &&
								e[1].Row == right[1].Row &&
								e[0].Col-right[0].Col == 1
						})
						if !found {
							break
						}
						fmt.Println("right ", right)
					}
				}
			}
			sideCount++
			fmt.Println("done 1 side", sideCount)
		}
		fmt.Println("done one area", area, sideCount, area*sideCount)
		total += (area * sideCount)
	}
	fmt.Println(total)
}

func searchAndRemove[T any](array *[]T, check func(condition T) bool) (T, bool) {
	var ret T
	slice := *array
	fmt.Println(slice)
	defer fmt.Println(slice)
	for i := range slice {
		if check(slice[i]) {
			ret = slice[i]
			slice[i], slice[len(slice)-1] = slice[len(slice)-1], slice[i]
			*array = (*array)[:len(slice)-1]
			return ret, true
		}
	}
	return ret, false
}

var allDir = []util.Direction{
	util.North,
	util.East,
	util.South,
	util.West,
}
