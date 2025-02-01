package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var (
	debug         = false
	delay         = 50 * time.Millisecond
	counter int32 = 0
	// maxStep       = 10_000
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
}
