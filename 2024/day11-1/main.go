package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
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

	cache := map[int64]map[int]int{}
	allRoots := []*Node{}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := strings.Split(scanner.Text(), " ")
	for _, str := range line {
		num, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		allRoots = append(allRoots, &Node{
			Data: int64(num),
		})
	}

	limit := 75
	score := 0
	for _, root := range allRoots {
		score += DFS(cache, root, limit)
	}
	fmt.Println(score)
}

type Node struct {
	Data     int64
	Children []*Node
}

func DFS(cache map[int64]map[int]int, n *Node, blink int) int {
	if blink == 0 {
		return 1
	}

	if count, found := getCache(cache, n.Data, blink); found {
		return count
	}

	newData := []int64{}
	if len(n.Children) == 0 {
		if n.Data == 0 {
			newData = []int64{1}
		} else if split := splitDigits3(n.Data); len(split) == 2 {
			newData = split
		} else {
			newData = []int64{n.Data * 2024}
		}
		for i := range newData {
			n.Children = append(n.Children, &Node{
				Data:     newData[i],
				Children: []*Node{},
			})
		}
	}

	// the leaf node of this node is sum of leaf node of it's children
	count := 0
	for _, child := range n.Children {
		count += DFS(cache, child, blink-1)
	}

	// update cache
	insertCache(cache, n.Data, blink, count)
	return count
}

func insertCache(cache map[int64]map[int]int, data int64, blink int, count int) {
	if _, found := cache[data]; !found {
		cache[data] = map[int]int{}
	}
	cache[data][blink] = count
}

func getCache(cache map[int64]map[int]int, data int64, blink int) (int, bool) {
	if entry, found := cache[data]; found {
		if count, found := entry[blink]; found {
			return count, true
		}
	}
	return 0, false
}

func splitDigits3(n int64) []int64 {
	count := 0
	cp := n
	for cp > 0 {
		cp /= 10
		count++
	}

	if count%2 != 0 {
		return nil
	}
	mid := count / 2

	divisor := int64(1)
	for mid > 0 {
		divisor *= 10
		mid--
	}
	return []int64{n / divisor, n % divisor}
}
