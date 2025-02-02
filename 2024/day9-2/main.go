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
	fileX, fileY := len(disk), 0
	lock := map[string]bool{}
	for {
		fileX, fileY = nextLeftFile(disk, fileX)
		if fileX == -1 {
			break
		}
		id := disk[fileX]

		if lock[id] {
			continue
		}
		lock[id] = true

		size := fileY - fileX + 1
		freeX, freeY := rightMostFreeSpace(disk, size)
		if freeX == -1 {
			fmt.Printf("no size %d found\n", size)
			continue
		} else if fileX < freeY {
			fmt.Printf("size %d found but not on left of file\n", size)
			continue
		}
		fmt.Printf("free space for file %s size %d found\n", disk[fileX], size)

		fmt.Println(freeX, freeY, fileX, fileY)
		err = util.SwapSlice(disk, freeX, freeY+1, fileX, fileY+1)
		if err != nil {
			panic(err)
		}
		fmt.Println(disk)
	}
	fmt.Println("checksum", checkSum(disk))
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
		// from--
		to = from
		for to < len(disk) && disk[to] == "." {
			to++
		}
		to--

		if size == -1 {
			return from, to
		} else if to-from+1 >= size {
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
			continue
		}
		num, err := strconv.Atoi(disk[i])
		if err != nil {
			panic(err)
		}
		score += num * i
	}
	return score
}
