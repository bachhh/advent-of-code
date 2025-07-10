package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strings"

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

	ds := util.NewDisjointSet[string]()
	for scanner.Scan() {
		text := scanner.Text()
		v := strings.Split(text, "-")
		v1, v2 := v[0], v[1]
		ds.Union(v1, v2)
		fmt.Println(v1, v2, ds.ToSlice())
	}
	total := 0
	slces := ds.ToSlice()
	for _, sl := range slces {
		if len(sl) < 3 {
			continue
		}
		slices.SortFunc(sl, func(a, b string) int {
			if a[0] == 't' && b[0] != 't' {
				return -1
			}
			if a[0] != 't' && b[0] == 't' {
				return 1
			}
			if a[0] == 't' && b[0] == 't' {
				return int(a[1]) - int(b[1])
			}
			return 0
		})
		fmt.Println(sl)
		for i, str := range sl {
			if str[0] == 't' {
				// you choose 2 out of the remaining characters
				c := choose(len(sl)-1-i, 2)
				fmt.Println(total, c)
				total += c
			}
		}
	}
	fmt.Println(total)
}

// choose k of of N member
func choose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k { // Take advantage of symmetry
		k = n - k
	}
	res := 1
	for i := 1; i <= k; i++ {
		res = res * (n - i + 1) / i
	}
	return res
}
