package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

	file, err := os.Open("input.txt")
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

	// fmt.Println(registers)
	// fmt.Println(opcodes)
	// fmt.Println(registers)

	fmt.Println(registers, opcodes)
	maxLen := 0
	// start := 10000000000000
	start := 0
	limit := 512

	for a := start; a <= limit; a++ {
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
	insPtr := 0
	for {
		if insPtr >= len(opcodes) {
			break
		}
		Eval(reg, opcodes[insPtr], opcodes[insPtr+1], &insPtr, &output)
		// fmt.Println("step", reg["A"])
	}
	return output
}

func Eval(registers map[string]int, opcode, operand int, insPtr *int, output *[]int) bool {
	switch Opcode(opcode) {

	case adv:
		result := registers["A"] / (1 << comboOperand(registers, operand))
		result = trunc(result)
		registers["A"] = result

	case bdv:
		result := registers["A"] / (1 << comboOperand(registers, operand))
		result = trunc(result)
		registers["B"] = result

	case cdv:
		result := registers["A"] / (1 << comboOperand(registers, operand))
		result = trunc(result)
		registers["C"] = result

	case bxl:
		result := registers["B"] ^ operand
		registers["B"] = result

	case bst:
		result := comboOperand(registers, operand) % 8
		registers["B"] = result

	case jnz:
		if registers["A"] == 0 {
			break
		}
		*insPtr = operand
		goto DONE

	case bxc:
		results := registers["B"] ^ registers["C"]
		registers["B"] = results

	case out:
		*output = append(*output, comboOperand(registers, operand)%8)

	}
	*insPtr = *insPtr + 2

DONE:
	return true
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
