package main

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

	file, err := os.Open("test_input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	registers := map[string]int{"A": 0, "B": 0, "C": 0}

	var opcodes []int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register") {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "Register ")
				if _, exists := registers[key]; exists {
					value, err := strconv.Atoi(parts[1])
					if err == nil {
						registers[key] = value
					}
				}
			}
		} else if strings.HasPrefix(line, "Program:") {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				numbers := strings.Split(parts[1], ",")
				for _, num := range numbers {
					val, err := strconv.Atoi(strings.TrimSpace(num))
					if err == nil {
						opcodes = append(opcodes, val)
					}
				}
			}
		}
	}

	fmt.Println(registers, opcodes)
	maxLen := 0
	// start := 10000000000000
	start := 117440
	runCounts := 1
	// fmt.Println(printHighlightBits(117440, 3, 3))

	for a := start; a < start+runCounts; a++ {
		registers["A"] = a
		output := run(registers, opcodes)
		if len(output) > maxLen {
			maxLen = len(output)
			fmt.Println(maxLen, a, a, output)
		}
		fmt.Println("end run", fmt.Sprintf("%b", a), opcodes, output, a)
		if slices.Compare(output, opcodes) == 0 {
			// fmt.Println(a)
			break
		}
	}

	// if slices.Compare(opcodes, intOutput) == 0 {
	// }
}

func run(registers map[string]int, opcodes []int) []int {
	var output []int
	reg := util.CloneMap(registers)
	fmt.Println("before run ", printHighlightBits(reg["A"], 999, 3))
	insPtr := 0
	for {
		if insPtr >= len(opcodes) {
			break
		}
		mod8xor3 := (reg["A"] % 8) ^ 3
		aHighLight := printHighlightBits(reg["A"], mod8xor3, 3)

		aHighLight2 := printHighlightBits(reg["A"]/(1<<mod8xor3), 99, 9)
		before := fmt.Sprintf("%s -%d - %03b\n", aHighLight, mod8xor3, reg["A"]%8)
		if Eval(reg, opcodes[insPtr], opcodes[insPtr+1], &insPtr, &output) {
			fmt.Println(aHighLight2)
			fmt.Println(before, output, fmt.Sprintf("C:%03b B:%03b", reg["C"], reg["B"]))
		}
		// fmt.Println("step", reg["A"])
	}
	return output
}

func Eval(reg map[string]int, opcode, operand int, insPtr *int, output *[]int) bool {
	hasOutput := false
	comboOp := comboOperand(reg, operand)
	switch Opcode(opcode) {
	case adv:
		result := reg["A"] / (1 << comboOp)
		result = trunc(result)
		reg["A"] = result

	case bdv:
		result := reg["A"] / (1 << comboOp)
		result = trunc(result)
		reg["B"] = result

	case cdv:
		result := reg["A"] / (1 << comboOp)
		result = trunc(result)
		reg["C"] = result

	case bxl:
		result := reg["B"] ^ operand
		reg["B"] = result

	case bst:
		result := comboOp % 8
		reg["B"] = result

	case jnz:
		if reg["A"] == 0 {
			break
		}
		*insPtr = operand
		goto DONE

	case bxc:
		results := reg["B"] ^ reg["C"]
		reg["B"] = results

	case out:
		*output = append(*output, comboOp%8)
		hasOutput = true

	}
	*insPtr = *insPtr + 2

DONE:
	return hasOutput
}

func comboOperand(registers map[string]int, i int) int {
	switch i {
	case 0, 1, 2, 3:
		return i
	case 4:
		return registers["A"]
	case 5:
		return registers["B"]
	case 6:
		return registers["C"]
	case 7:
		panic("invalid combo operand 7")
	}
	panic(fmt.Sprintf("invalid combo operand %d", i))
}

func trunc(i int) int {
	return i
	// return i % 8
}

type Opcode int

const (
	// division by 2^combowrite to a,
	adv Opcode = iota
	// bitwise XOR
	bxl
	// combine operand, mod 8, write to B
	bst

	// if A == 0 -> noop, jump to ins pointer otherwise
	jnz

	// bitwisexor(B, C) write to B
	bxc

	// calc combo, output value
	out

	// same as adv but write to a
	bdv
	// same as adv but write to c
	cdv
)

// 2,4  → B = A MOD 8
// 1,3  → B = B XOR 3
// 7,5  → C = A DIV 2^B
// 0,3  → A = A DIV 2^3
// 1,5  → B = B XOR B  (effectively setting B = 0)
// 4,1  → B = B XOR C  (B = C)
// 5,5  → OUTPUT B MOD 8
// 3,0  → IF A ≠ 0 THEN GOTO 0

// 001 XOR 011 = 010
// 010 XOR 011 = 001
// 011 XOR 011 = 000
// 100 XOR 011 = 111
// 101 XOR 011 = 110
// 110 XOR 011 = 101
// 111 XOR 011 = 100

// 0 3  5 4 3 0
//  3    4   5   3  0

// 0b

// 000 011 100 101 011 000 000
// 000 011 100 101 011 000
// 000 011 100 101 011
// 000 011 100 101
// 000 011 100
// 000 011
// 000
//

// 000 3   4    5   3  0    0
// 000 011 100 101 011 000 000

// 1. mod 8 = 000 -> XOR 3 = 011 => C = a DIV 2^1 = 000 mod 8 = 0
// 100 011 = 111

// 11100101011000
// 1. mod 8 = 000 -> XOR 3 = 011 => C = a DIV 2^3 = 011000 mod 8 = 0
// 001 xOR 011 = 010 = 2
// 111001 >> 2 =  1110
// 1 <<  2 = 100
// 111001

// ANSI escape codes for red color
const (
	redStart = "\033[31m" // Start red text
	redEnd   = "\033[0m"  // Reset color
)

// print A in binary form, but highlight <count> bits, offset by <offset> from the least significant bits of A.
// print in big endian

// printHighlightBits prints A in binary form (big-endian),
// highlighting <count> bits starting at <offset> from LSB in red.
func printHighlightBits(A int, offset int, count int) string {
	var padLen int
	if A == 0 {
		padLen = 1
	} else {
		padLen = int(math.Log2(float64(A))) + 1
	}

	if padLen%3 != 0 {
		padLen = padLen + (3 - padLen%3)
	}

	// Convert A to a binary string (ensure 32-bit representation)
	binStr := fmt.Sprintf("%0*b", padLen, A)

	// Reverse indexing to locate bits from the right
	length := len(binStr)
	result := ""

	for i := 0; i < length; i++ {
		if i >= length-1-offset-count+1 && i <= length-1-offset {
			// Highlight in red
			result += redStart + string(binStr[i]) + redEnd
		} else {
			result += string(binStr[i])
		}
	}

	return result
}
