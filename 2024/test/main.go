package main

import (
	"fmt"
	"math"
	"strings"
)

// ANSI escape codes for red color
const (
	redStart = "\033[31m" // Start red text
	redEnd   = "\033[0m"  // Reset color
)

// getBitLength returns the number of bits needed to represent A in binary
func getBitLength(A int) int {
	if A == 0 {
		return 1
	}
	return int(math.Log2(float64(A))) + 1
}

// nearestMultipleOf3 finds the smallest multiple of 3 that is >= num
func nearestMultipleOf3(num int) int {
	if num%3 == 0 {
		return num
	}
	return num + (3 - num%3)
}

// formatBinaryWithSpacing inserts spaces every 3 bits in a binary string
func formatBinaryWithSpacing(binStr string) string {
	var result strings.Builder
	length := len(binStr)
	spaceCount := 0

	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(" ") // Insert space every 3 bits
			spaceCount++
		}
		result.WriteByte(binStr[i])
	}

	return result.String()
}

// printHighlightBits prints A in binary form (big-endian),
// highlighting <count> bits starting at <offset> from LSB in red.
func printHighlightBits(A int, offset int, count int) string {
	// Determine required bit length
	bitLength := getBitLength(A)
	paddedLength := nearestMultipleOf3(bitLength)

	// Convert A to a binary string with dynamic padding
	binStr := fmt.Sprintf("%0*b", paddedLength, A)

	// Apply spaces first
	spacedBinStr := formatBinaryWithSpacing(binStr)

	// Convert to a slice for modification
	runes := []byte(spacedBinStr)
	length := len(runes)

	// Apply highlighting (taking spaces into account)
	bitIndex := len(binStr) - 1 - offset
	for i := length - 1; i >= 0 && count > 0; i-- {
		if runes[i] == ' ' {
			continue // Skip spaces
		}
		if bitIndex < 0 {
			break
		}
		// Insert red ANSI escape codes
		runes[i] = 'X' // Mark bit for highlighting
		bitIndex--
		count--
	}

	// Build final string with ANSI color
	var highlightedBin strings.Builder
	bitIndex = len(binStr) - 1 - offset
	for i := 0; i < length; i++ {
		if runes[i] == 'X' {
			highlightedBin.WriteString(redStart + string(binStr[bitIndex]) + redEnd)
			bitIndex--
		} else {
			highlightedBin.WriteByte(runes[i])
		}
	}

	return highlightedBin.String()
}

func main() {
	// Example: A = 29 (binary 11101), highlight 3 bits from offset 1
	fmt.Println(printHighlightBits(12401024, 3, 3))
}
