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
	"time"
)

var delay = 500 * time.Millisecond

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

	total := int64(0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ints, err := parseLine(line)
		if err != nil {
			panic(err)
		}

		target, numbers := ints[0], ints[1:]
		if eval(numbers[0], target, numbers[1:]) {
			total += target
		}
	}
	fmt.Println(total)
}

func eval(accum int64, target int64, numbers []int64) bool {
	if len(numbers) == 0 {
		return accum == target
	}
	return eval(accum*numbers[0], target, numbers[1:]) ||
		eval(accum+numbers[0], target, numbers[1:]) ||
		eval(combine(accum, numbers[0]), target, numbers[1:])
}

func combine(a, b int64) int64 {
	multiplier := int64(math.Pow10(int(math.Log10(float64(b))) + 1))

	return a*multiplier + b
}

func parseLine(input string) ([]int64, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	result := []int64{}
	firstNum, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return nil, err
	}
	result = append(result, firstNum)

	numbers := strings.Fields(parts[1])
	for _, numStr := range numbers {
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

func isBitSet(number, bit int) bool {
	return (number & (1 << bit)) != 0
}

func cloneSlice[T ~[]E, E any](slice T) T {
	cp := make(T, len(slice))
	copy(cp, slice)
	return cp
}
