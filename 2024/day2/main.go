package main

import (
	"bufio"
	"fmt"
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

func main() {
	score := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		ints := readLineInts(scanner.Text())
		if len(ints) < 2 { // auto safe
			score++
			continue
		}
		isIncrease := ints[1] > ints[0]
		isSafe := true
		for i := 0; i < len(ints)-1; i++ {
			gap := ints[i+1] - ints[i]
			if gap == 0 {
				isSafe = false
				break
			}
			if isIncrease && (gap <= 0 || gap > 3) {
				isSafe = false
				break
			} else if !isIncrease && (gap < -3 || gap >= 0) {
				isSafe = false
				break
			}
		}
		fmt.Println(ints)
		if isSafe {
			score++
		}
	}
	println(score)
}
