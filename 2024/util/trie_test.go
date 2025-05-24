package util

import (
	"fmt"
	"testing"
)

func TestRadixTree(t *testing.T) {
	input := []string{
		"roman",
		"romane",
		"romanus",
		"romulus",
		"rubens",
		"ruber",
		"rubicon",
		"rubric",
	}
	root := &RadixTree{}
	for _, str := range input {
		root.Insert(str, str)
		PrintRadixTree(root, "root", "", false)
	}

	fmt.Println("done all")
	PrintRadixTree(root, "root", "", false)
}

func FuzzRadixTree(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		root := &RadixTree{}
		fmt.Println(orig)
		root.Insert(orig, orig)
	})
}
