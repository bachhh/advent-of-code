package mai

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"

	"aoc2024/util"
)

var allowedDir = []util.Direction{util.North, util.East, util.South, util.West}

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

	var array []int64
	scanner.Scan()
	line := scanner.Text()
	split := strings.Split(line, " ")
	for i := range split {
		num, err := strconv.Atoi(split[i])
		if err != nil {
			panic(err)
		}
		array = append(array, int64(num))
	}
	fmt.Println(array)
	for i := 0; i < 25; i++ {
		eval(array)
		fmt.Println(array)
	}
}

func eval(array []int64) {
	for i := range array {
		if array[i] == 0 {
			array[i] = 1
		} else if split := splitDigits(array[i]); len(split) == 2 {
			array = slices.Insert(array, i, split[0])
			array[i+1] = split[1]
		} else {
			array[i] *= 2024
		}
	}
}

func splitDigits(n int64) []int64 {
	strNum := strconv.Itoa(int(math.Abs(float64(n))))
	digitCount := len(strNum)
	if digitCount%2 != 0 {
		return nil
	}

	mid := len(strNum) / 2

	// Convert the two halves back to integers
	firstHalf, _ := strconv.ParseInt(strNum[:mid], 10, 64)
	secondHalf, _ := strconv.ParseInt(strNum[mid:], 10, 64)
	return []int64{firstHalf, secondHalf}
}
