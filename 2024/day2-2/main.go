package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func checkLine(ints []int) int {
	if len(ints) < 2 { // auto safe
		return -1
	}
	isIncrease := ints[1] > ints[0]
	for i := 0; i < len(ints)-1; i++ {
		gap := ints[i+1] - ints[i]
		if gap == 0 {
			return i
		}
		if isIncrease && (gap <= 0 || gap > 3) {
			return i
		} else if !isIncrease && (gap < -3 || gap >= 0) {
			return i
		}
	}
	return -1
}

func canFix(ints []int, tries []int) ([]int, bool) {
	for _, fi := range tries {
		cp := make([]int, len(ints))
		copy(cp, ints)
		cp = slices.Delete(cp, fi, fi+1)
		if find := checkLine(cp); find == -1 {
			return cp, true
		}
	}
	return nil, false
}

func main() {
	score := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		ints := readLineInts(scanner.Text())

		idx := checkLine(ints)
		if idx == -1 {
			score++
		} else {

			// either remove the first or the second offending index
			// check if the remaining array still valid
			// tries := []int{idx, idx + 1}
			allTries := []int{}
			for i := range ints {
				allTries = append(allTries, i)
			}

			if _, ok := canFix(ints, allTries); ok {
				score++
				goto END
			} else {
				// fmt.Printl(alltries)
			}

			fmt.Println(ints, "unsafe at index", idx)
		END:
		}
	}
	fmt.Println(score)
}
