package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"aoc2024/util"
)

type Pair = util.Pair

const MAX_DEPTH = 2000

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

	inputFile := flag.Arg(0)
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	maxPrice, maxSeq := -999999, ""
	for scanner.Scan() {
		input := scanner.Text()
		secret, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			panic(err)
		}

		dprice := toSeq(secret)
		mapPrice := toMapPrice(dprice)
		// fmt.Println(dprice)
		// fmt.Println(mapPrice)
		for k, v := range mapPrice {
			globalSeqToPriceMap[k] += v
		}
	}

	for seq, price := range globalSeqToPriceMap {
		if price > maxPrice {
			maxPrice, maxSeq = price, seq
		}
	}
	fmt.Println(maxPrice, maxSeq)
}

var globalSeqToPriceMap = map[string]int{}

func toMapPrice(input deltaAndPrice) map[string]int {
	m := map[string]int{}
	for i := 3; i < len(input.delta); i++ {
		seqStr := toSeqString(input.delta[i-3 : i+1])
		if _, found := m[seqStr]; found {
			// not the first time this seq exists, skip
			continue
		}
		m[seqStr] = input.price[i]
	}
	return m
}

var digits = [...]string{"-9", "-8", "-7", "-6", "-5", "-4", "-3", "-2", "-1", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func toSeqString(in []int) string {
	ret := ""
	for _, num := range in {
		ret += digits[num+9]
	}
	return ret
}

func toSeq(secret int64) deltaAndPrice {
	ret := deltaAndPrice{}
	prevPrice := secret % 10
	for i := 0; i < MAX_DEPTH; i++ {
		secret = evolve(secret)
		curPrice := secret % 10
		ret.price = append(ret.price, int(curPrice))
		ret.delta = append(ret.delta, int(curPrice-prevPrice))
		prevPrice = curPrice
	}
	return ret
}

type deltaAndPrice struct {
	delta []int
	price []int
}

func evolve(secret int64) int64 {
	// Calculate the result of multiplying the secret number by 64. Then, mix this
	// result into the secret number. Finally, prune the secret number.
	secret = prune(mix(secret, secret*64))

	// Calculate the result of dividing the secret number by 32. Round the result
	// down to the nearest integer. Then, mix this result into the secret number.
	// Finally, prune the secret number.
	secret = mix(secret, secret/32)

	// Calculate the result of multiplying the secret number by 2048. Then, mix
	// this result into the secret number. Finally, prune the secret number.
	secret = prune(mix(secret, secret*2048))
	return secret
}

// To mix a value into the secret number, calculate the bitwise XOR of the given
// value and the secret number. Then, the secret number becomes the result of that
// operation. (If the secret number is 42 and you were to mix 15 into the secret
// number, the secret number would become 37.)
func mix(a, b int64) int64 {
	return a ^ b
}

// To prune the secret number, calculate the value of the secret number modulo
// 16777216. Then, the secret number becomes the result of that operation. (If the
// secret number is 100000000 and you were to prune the secret number, the secret
// number would become 16113920.)
func prune(a int64) int64 {
	return a % 16777216
}
