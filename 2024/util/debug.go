package util

import (
	"flag"
	"fmt"
	"time"
)

var debug = flag.Bool("debug", false, "enable debug output")

// take a start time, print every <period> second
func PeriodicPrintf(start *time.Time, period time.Duration, format string, args ...interface{}) {
	t := time.Now()
	if time.Since(*start) > period {
		fmt.Printf(format, args)

		*start = t
	}
}

func PeriodicPrintln(start *time.Time, period time.Duration, input ...any) {
	t := time.Now()
	if time.Since(*start) > period {
		fmt.Println(input)

		*start = t
	}
}

func Debugln(stuffs ...any) {
	if *debug {
		fmt.Println(stuffs...)
	}
}

func Debugf(str string, args ...any) {
	if *debug {
		fmt.Printf(str, args...)
	}
}
