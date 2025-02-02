package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Pair struct {
	Before int
	After  int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	ordering := []Pair{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // finished reading pair
			break
		}
		a, b := readPair(line)
		ordering = append(ordering, Pair{Before: a, After: b})
	}
	slices.SortFunc(ordering, func(x Pair, y Pair) int {
		if x.Before == y.Before {
			return cmp.Compare(x.After, y.After)
		}
		return cmp.Compare(x.Before, y.Before)
	})
	fmt.Println(ordering)

	score := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // finished reading pair
			break
		}
		update := readList(line)
		if badIdx := checkOrdering(update, ordering); badIdx == -1 {
			// midIdx := len(update) / 2
			// score += update[midIdx]
			fmt.Println("good1", update)
		} else {
			// fmt.Println("bad original", update)
			update = fixOrdering(update, ordering)
			if badIdx := checkOrdering(update, ordering); badIdx == -1 {
				midIdx := len(update) / 2
				score += update[midIdx]
				fmt.Println("good2", update, update[midIdx], score)
			} else {
				fmt.Println("still bad", badIdx, update)
			}
		}
	}
	fmt.Println(score)
}

func fixOrdering(update []int, ordering []Pair) []int {
	slices.SortStableFunc(update, func(a, b int) int {
		_, found := slices.BinarySearchFunc(ordering, Pair{Before: a, After: b}, func(element Pair, target Pair) int {
			if c := cmp.Compare(element.Before, target.Before); c == 0 {
				return cmp.Compare(element.After, target.After)
			} else {
				return c
			}
		})
		if found {
			return -1
		}

		// swap the order, and check if b should be before a
		_, found = slices.BinarySearchFunc(ordering, Pair{Before: b, After: a}, func(element Pair, target Pair) int {
			if c := cmp.Compare(element.Before, target.Before); c == 0 {
				return cmp.Compare(element.After, target.After)
			} else {
				return c
			}
		})
		if found {
			return 1
		}
		return 0 // unorderable
	})
	return update
}

func checkOrdering(update []int, ordering []Pair) int {
	for i := 0; i < len(update)-1; i++ {
		for j := i + 1; j < len(update); j++ {
			a, b := update[i], update[j]
			// find the second number first
			// we have a pair a and b, check if there are any rules that require b before a
			idx, found := sort.Find(len(ordering), func(m int) int {
				return cmp.Compare(b, ordering[m].Before)
			})
			if !found || idx == len(ordering) {
				continue
			}
			for ; idx < len(ordering) && ordering[idx].Before == b; idx++ {
				if ordering[idx].After == a {
					return i
				}
			}

		}
	}
	return -1
}

func readPair(str string) (int, int) {
	strSlice := strings.Split(str, "|")
	if len(strSlice) != 2 {
		panic("invalid input " + str + "should be a|b")
	}
	a, err := strconv.Atoi(strSlice[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(strSlice[1])
	if err != nil {
		panic(err)
	}
	return a, b
}

func readList(str string) []int {
	strSlice := strings.Split(str, ",")
	intSlice := []int{}
	for i := range strSlice {
		a, err := strconv.Atoi(strSlice[i])
		if err != nil {
			panic(err)
		}
		intSlice = append(intSlice, a)
	}
	return intSlice
}
