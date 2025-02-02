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

	fmt.Println(disk)
	var fr, to int
	fr = len(disk)
	for {
		// fr, to = nextLeftFile(disk, fr)
		// if fr == -1 {
		// 	break
		// }

		fr, to = rightMostFreeSpace(disk, -1)
		if fr == -1 {
			break
		}

		for i := range disk[fr : to+1] {
			disk[i] = "+"
		}

		fmt.Println(disk[fr : to+1])
		break
	}
	fmt.Println(checkSum(disk))
}

// look for the next file to the left of pos
func nextLeftFile(disk []string, pos int) (from, to int) {
	to = pos - 1
	for to >= 0 && disk[to] == "." {
		to--
	}
	if to < 0 {
		return -1, -1
	}
	from = to
	for from > -1 && disk[from] == disk[to] {
		from--
	}
	from++
	return from, to
}

// look for the right most free space of a certain size
func rightMostFreeSpace(disk []string, size int) (from, to int) {
	to = -1
	for from < len(disk) {
		from = to + 1
		for from < len(disk) && disk[from] != "." {
			from++
		}
		if from == len(disk) {
			return -1, -1
		}
		from--
		to = from
		for to < len(disk) && disk[to] == "." {
			to++
		}
		to--

		if size == -1 {
			return from, to
		} else if to-from+1 == size {
			return from, to
		}
	}

	if from == len(disk) {
		return -1, -1
	}
	return from, to
}

func checkSum(disk []string) int {
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
	return score
}

func diskMapToDisk(diskmap []byte) []string {
	isSpace := false
	disk := []string{}
	id := 0

	for i := range diskmap {

		num, err := util.CharToInt(diskmap[i])
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
	return disk
}
