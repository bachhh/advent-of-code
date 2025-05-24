package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"

	"aoc2024/util"
)

type Pair = util.Pair

// Solution: build an optmized radix tree, run a bfs on each node
// each time we hit a leaf node on of the tree, a prefix part of the design can
// be cut from a towels.
// If cutting the design result in an empty string, then we have a "possible" case.
// The issue here is the size of tree: there are hundreds of node. Exponentially increase our BFS search space
//
// Optimization 1: During building of our Radix Tree, if we can make the input string from already existing node in the tree. Then we can skip this input. e.g.: if we already have "ab", and "cd" as 2 leaf node. then we do not need to also store "abcd".
//
// Optimization 2: Insertion order. If the sequence of insert are "abcd", "ab",
// "cd", when we get to "cd", "abcd" and also any leaf node with postfix "cd"
// should be removed. Since looking up postfix is significantly harder in Radix
// Tree. we can just sort input by lenght of string. This guaranteed that
// "abcd" is always inserted after "ab" and "cd", an can be pruned more easily.

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

	inputFile := flag.Arg(0)
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Read first line
	scanner.Scan()
	firstLine := scanner.Text()
	towels := strings.Split(firstLine, ",")
	for i := range towels {
		towels[i] = strings.TrimSpace(towels[i])
	}

	// Skip the blank line
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			break
		}
	}

	// Read remaining lines
	var designs []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			designs = append(designs, line)
		}
	}

	// optimize
	sort.Slice(towels, func(i, j int) bool {
		return len(towels[i]) < len(towels[j])
	})

	trie := &util.RadixTree{}
	for _, str := range towels {
		if str == "" {
			panic("empty string towel")
		}
		// Optimize 1: if we have a 2 full string encoded the radix tree: ab , cd,
		// then we don't need to insert the full string "abcd" anymore.
		// the same string can be composed of 2 smaller string
		if trie.CanCompose(str) {
			continue
		}
		trie.Insert(str, str)
	}
	util.PrintRadixTree(trie, "root", "", false)

	counter := 0
	for _, design := range designs {
		if test(design, trie) {
			counter++
			fmt.Println("OK", design)
		} else {
			fmt.Println("NG", design)
		}
	}
	fmt.Println("counter", counter)
}

func test(design string, root *util.RadixTree) bool {
	outerQ := util.NewQueue[string]()
	outerQ.Push(design)
	for !outerQ.IsEmpty() {
		outerStr, _ := outerQ.Pop()
		// fmt.Println("outerQ size", outerQ.Size())

		innerQ := util.NewQueue[innerQData]()

		innerQ.Push(innerQData{
			Node:      root,
			Remaining: outerStr,
		})

		// bfs walk the tree until we reach a leaf node
		// any tree node that matches with a prefix of the string is pushed
		// back into the queue
		// any time we reach a leaf node. The remnant is put back into the
		// outer queue and
		for !innerQ.IsEmpty() {
			next, _ := innerQ.Pop()
			node, innerStr := next.Node, next.Remaining
			if node.IsLeafNode {
				if innerStr == "" {
					return true
				}
				outerQ.Push(innerStr)
				// even if this one is a leaf node, we might be able to
				// consume more of the string
			}

			for key, child := range node.Children {
				if remain, found := strings.CutPrefix(innerStr, key); found {
					// fmt.Printf("key %s for %s matched, prefix=%s\n", key, innerStr, remain)
					innerQ.Push(innerQData{
						Node:      child,
						Remaining: remain,
					})
				}
			}
		}
	}
	return false
}

type innerQData struct {
	Node      *util.RadixTree
	Remaining string
}
