package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
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

	for scanner.Scan() {
	}
}

func checkInBound(point Pair, matrix [][]byte) bool {
	return point.Row >= 0 && point.Row < len(matrix) && point.Col >= 0 && point.Col < len(matrix[0])
}

func GetAntinode(a, b Pair, matrix [][]byte) []Pair {
	vecAB := Pair{Row: b.Row - a.Row, Col: b.Col - a.Col}
	invVecAb := InvertVector(vecAB)

	results := []Pair{}
	for cp := AddVector(b, vecAB); checkInBound(cp, matrix); cp = AddVector(cp, vecAB) {
		// fmt.Println("gen antinode ", b, vecAB, cp)
		results = append(results, cp)
	}
	for cp := AddVector(a, invVecAb); checkInBound(cp, matrix); cp = AddVector(cp, invVecAb) {
		// fmt.Println("gen antinode ", a, invVecAb, cp)
		results = append(results, cp)
	}
	return results
}

type Pair struct {
	Row int
	Col int
}

func AddVector(point Pair, vector Pair) Pair {
	return Pair{
		Row: point.Row + vector.Row,
		Col: point.Col + vector.Col,
	}
}

func InvertVector(vec Pair) Pair {
	return Pair{
		Row: -vec.Row,
		Col: -vec.Col,
	}
}
