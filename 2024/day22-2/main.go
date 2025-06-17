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

	total := int64(0)
	for scanner.Scan() {
		input := scanner.Text()
		integer, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(integer)
	}
	fmt.Println(total)
}

func checkPrice(b Buyer, seq []int) int {
	for i := 0; i < len(b.delta); i++ {
		for j := range seq {
			if b.delta[i] != seq[j] {
				goto SKIP
			}
		}
		// have a match
		return b.price[i+len(seq)-1]
	SKIP:
	}
	return -1
}

func toSeq(secret int64) Buyer {
	ret := Buyer{}
	prevPrice := secret % 10
	for i := 0; i < 2000; i++ {
		secret = evolve(secret)
		curPrice := secret % 10
		ret.price = append(ret.price, int(curPrice))
		ret.delta = append(ret.delta, int(curPrice-prevPrice))
		prevPrice = curPrice
	}
}

type Buyer struct {
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
