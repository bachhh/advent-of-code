package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func readLineInts(line string) []int {
	ints := make([]int, 0)
	for _, part := range strings.Fields(line) {
		if num, err := strconv.Atoi(part); err == nil {
			ints = append(ints, num)
		}
	}
	return ints
}

func checkLine(ints []int) bool {
	isIncrease := ints[1] > ints[0]
	for i := 0; i < len(ints)-1; i++ {
		gap := ints[i+1] - ints[i]
		if gap == 0 {
			return false
		}
		if isIncrease && (gap <= 0 || gap > 3) {
			return false
		} else if !isIncrease && (gap < -3 || gap >= 0) {
			return false
		}
	}
	return true
}

func main1() {
	score := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		ints := readLineInts(scanner.Text())
		if len(ints) < 2 { // auto safe
			score++
			continue
		}

		if checkLine(ints) {
			score++
		}
	}
}

func main2() {
	score := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		ints := readLineInts(scanner.Text())
		if len(ints) < 2 { // auto safe
			score++
			continue
		}

		if checkLine(ints) {
			score++
		}
	}
}

func main() {
}
