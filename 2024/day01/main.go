package main

import (
	"fmt"
	"io"
	"math"
	"sort"
)

func main1() {
	list1 := []int{}
	list2 := []int{}
	var num1, num2 int
	for {
		_, err := fmt.Scanln(&num1, &num2)
		if err == io.EOF {
			break
		}
		list1 = append(list1, num1)
		list2 = append(list2, num2)
		fmt.Println(num1, num2)
	}
	sort.Ints(list1)
	sort.Ints(list2)
	d := 0
	for i := range list1 {
		d += int(math.Abs(float64(list1[i] - list2[i])))
	}
	fmt.Println(d)
}

func main2() {
	list1 := []int{}
	list2 := []int{}
	var num1, num2 int
	for {
		_, err := fmt.Scanln(&num1, &num2)
		if err == io.EOF {
			break
		}
		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	sort.Ints(list2)
	var score int64
	for i := range list1 {
		idx1 := sort.SearchInts(list2, list1[i])
		if idx1 == len(list2) { // not found
			continue
		}
		if list2[idx1] == list1[i] {
			idx2 := idx1
			for ; list2[idx2] == list1[i] && idx2 < len(list2); idx2++ {
			}
			fmt.Println(list1[i], idx2-idx1)
			score += int64(list1[i] * (idx2 - idx1))
		}
	}
	fmt.Println(score)
}

func main3() {
	list1 := []int{}
	list2 := []int{}
	var num1, num2 int
	for {
		_, err := fmt.Scanln(&num1, &num2)
		if err == io.EOF {
			break
		}
		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}
}

func main() {
	main3()
}
