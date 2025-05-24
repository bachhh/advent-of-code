package util

import (
	"fmt"
	"time"
)

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
