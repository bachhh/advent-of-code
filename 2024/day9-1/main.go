package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

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

	cp := util.CloneSlice(scanner.Bytes())
	isSpace := false
	disk := []byte{}
	for i := range cp {
		isSpace = !isSpace

		num, err := util.CharToInt(cp[i])
		if err != nil {
			panic(err)
		}

		for i := range num {
			if isSpace {
				disk = append(disk, '.')
			} else {
				disk = append(disk, cp[i])
			}
		}
	}
	fmt.Println(string(disk))
}
