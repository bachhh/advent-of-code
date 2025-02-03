package main_test

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

// var testNum int64 = 1111111199999999
var testNum = map[int]int64{
	18: 111111111999999999,
	16: 1111111199999999,
	10: 1111199999,
	8:  11119999,
	4:  1199,
}

func BenchmarkSplitDigits(b *testing.B) {
	for size, input := range testNum {
		b.Run(fmt.Sprintf("BenchmarkSplitDigits1_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				splitDigits(input)
			}
		})
		b.Run(fmt.Sprintf("BenchmarkSplitDigits2_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				splitDigits2(input)
			}
		})
		b.Run(fmt.Sprintf("BenchmarkSplitDigits3_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				splitDigits3(input)
			}
		})
		b.Run(fmt.Sprintf("BenchmarkSplitDigits4_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				splitDigits4(input)
			}
		})
	}
}

func main() {
	fmt.Println("Benchmarking MyFunction:")
	testing.Benchmark(BenchmarkSplitDigits)
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

func splitDigits2(n int64) []int64 {
	digitCount := int(math.Log10(float64(n))) + 1
	if digitCount%2 != 0 {
		return nil
	}

	mid := digitCount / 2

	divisor := int64(math.Pow(10, float64(mid)))
	return []int64{n / divisor, n % divisor}
}

func splitDigits3(n int64) []int64 {
	count := 0
	cp := n
	for cp > 0 {
		cp /= 10
		count++
	}

	mid := count / 2

	divisor := int64(1)
	for mid > 0 {
		divisor *= 10
		mid--
	}
	return []int64{n / divisor, n % divisor}
}

var countSearch = []int64{
	0,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

func splitDigits4(n int64) []int64 {
	// fin the left most num which is greater than n
	count, _ := sort.Find(len(countSearch), func(idx int) int {
		if countSearch[idx] <= n {
			return 1
		}
		return 0
	})

	mid := count / 2

	divisor := countSearch[mid]
	return []int64{n / divisor, n % divisor}
}
