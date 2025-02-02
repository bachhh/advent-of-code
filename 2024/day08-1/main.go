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

	antenMap := map[byte][]Pair{}
	matrix := [][]byte{}
	for row := 0; scanner.Scan(); row++ {
		cp := util.CloneSlice(scanner.Bytes())
		// fmt.Println(string(cp))

		matrix = append(matrix, cp)
		for col := range cp {
			// if cp
			char := cp[col]
			// fmt.Println(char, IsAlphaNumeric(char))
			if IsAlphaNumeric(char) {
				antenMap[char] = append(antenMap[char], Pair{Row: row, Col: col})
			}
		}
	}

	util.PrintMatrix(matrix)
	counter := 0
	for _, antennas := range antenMap {
		for i := 0; i < len(antennas)-1; i++ {
			for j := i + 1; j < len(antennas); j++ {

				anti1, anti2 := GetAntinode(antennas[i], antennas[j])
				fmt.Println(antennas[i], antennas[j],
					anti1,
					checkInBound(anti1, matrix),
					anti2,
					checkInBound(anti2, matrix))

				if checkInBound(anti1, matrix) {
					if matrix[anti1.Row][anti1.Col] != '#' {
						matrix[anti1.Row][anti1.Col] = '#'
						counter++
					}
				}
				if checkInBound(anti2, matrix) {
					if matrix[anti2.Row][anti2.Col] != '#' {
						matrix[anti2.Row][anti2.Col] = '#'
						counter++
					}
				}
			}
		}
	}
	util.PrintMatrix(matrix)
	fmt.Println(counter)
}

func checkInBound(point Pair, matrix [][]byte) bool {
	return point.Row >= 0 && point.Row < len(matrix) && point.Col >= 0 && point.Col < len(matrix[0])
}

func GetAntinode(a, b Pair) (Pair, Pair) {
	vecAB := Pair{
		Row: b.Row - a.Row,
		Col: b.Col - a.Col,
	}
	invVecAb := InvertVector(vecAB)
	return AddVector(b, vecAB), AddVector(a, invVecAb)
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

func IsAlphaNumeric(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
