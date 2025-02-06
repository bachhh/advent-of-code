package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	addition := int64(10000000000000)
	total := int64(0)
	end := false
	for !end {
		data, err := ParseNumbers(scanner)
		if err != nil {
			fmt.Println("error", err)
			break
		}
		data.XP += addition
		data.YP += addition

		// either skip blank line or check for EOF
		if !scanner.Scan() {
			end = true
		}
		// fmt.Printf("%+v\n", data)
		B := (data.YP*data.XA - data.XP*data.YA) / (data.YB*data.XA - data.XB*data.YA)
		A := (data.XP - B*data.XB) / data.XA

		f1 := A*data.XA+B*data.XB == data.XP
		f2 := A*data.YA+B*data.YB == data.YP

		if !(f1 && f2) { // both formula has to be correct
			// invalid match
			fmt.Println("not possible")
		} else {
			fmt.Println(A, B, A*3+B)
			total += A*3 + B
		}

		// a*xa + b*xb = xp
		// a = ( xp -b*xb ) /xa

		// a*ya + b*yb = yp
		//( xp -b*xb )*ya /xa + b*yb = yp
		//  ya*xp - ya*b*xb  + b*yb*xa = yp * xa
		//   b* = yp*xa - ya*xp
		// b = ( yp*xa - ya*xp ) /(yb*xa - ya*xb  )

		// a = ( xp - b*xb ) / xa
		// ( ( xp - b*xb )*ya / xa ) + b*yb = yp
		// ( ( xp - b*xb )*ya / xa ) + b*yb = yp
		// ( ( xp - b*xb )*ya + b*yb*xa / xa  = yp
		// xp*ya - b*xb*ya + b*yb*xa = yp*xa
		//  b*( yb*xa - xb*ya ) = yp*xa - xp*ya
		// b = ( yp*xa - xp*ya )/( yb*xa - xb*ya )

	}
	fmt.Println("total", total)
}

func ParseNumbers(scanner *bufio.Scanner) (crane, error) {
	// EOF
	// Regular expression to match numbers in the format: X+num, Y+num, X=num, Y=num
	re := regexp.MustCompile(`\d+`)

	var numbers []int64
	ret := crane{}

	// Read and parse each line
	for i := 0; i < 3 && scanner.Scan(); i++ {
		matches := re.FindAllString(scanner.Text(), -1)
		for _, match := range matches {
			num, err := strconv.ParseInt(match, 10, 64)
			if err != nil {
				return ret, fmt.Errorf("invalid number: %s", match)
			}
			numbers = append(numbers, num)
		}
		// fmt.Println(scanner.Text())
	}

	// Ensure we have exactly 6 numbers
	if len(numbers) != 6 {
		return ret, fmt.Errorf("expected 6 numbers, got %d", len(numbers))
	}
	ret = crane{
		XA: numbers[0],
		YA: numbers[1],
		XB: numbers[2],
		YB: numbers[3],
		XP: numbers[4],
		YP: numbers[5],
	}

	return ret, nil
}

type crane struct {
	XA int64
	YA int64
	XB int64
	YB int64
	XP int64
	YP int64
}
