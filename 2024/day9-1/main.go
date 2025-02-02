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

	scanner.Scan()
	cp := util.CloneSlice(scanner.Bytes())

	isSpace := false
	disk := []string{}
	id := 0

	for i := range cp {

		num, err := util.CharToInt(cp[i])
		if err != nil {
			panic(err)
		}

		if isSpace {
			for range num {
				disk = append(disk, ".")
			}
		} else {
			for range num {
				char := strconv.Itoa((id))
				disk = append(disk, char)
			}
			id++
		}

		isSpace = !isSpace
		// fmt.Println(num, isSpace, string(cp[i]))
		// fmt.Println(string(disk))
	}
	// fmt.Println(string(disk))
	front, end := 0, len(disk)-1
	for ; disk[front] != "."; front++ {
	}
	for ; disk[end] == "."; end-- {
	}
	for front < end {
		disk[front], disk[end] = disk[end], disk[front]
		for front++; disk[front] != "."; front++ {
		}
		for end--; disk[end] == "."; end-- {
		}
		// fmt.Println(string(disk))
	}
	score := 0
	for i := range disk {
		if disk[i] == "." {
			break
		}
		num, err := strconv.Atoi(disk[i])
		if err != nil {
			panic(err)
		}
		score += num * i
	}
	fmt.Println(score)
}
