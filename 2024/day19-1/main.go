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

	sort.Strings(towels)
	fmt.Println("First line tokens:", towels)
	fmt.Println("Remaining inputs:", len(designs))

	for _, design := range designs {
		if test(design, towels) {
			fmt.Println("OK", design)
		} else {
			fmt.Println("NG", design)
		}
	}
}

func test(design string, towels []string) bool {
	qu := util.NewQueue[string]()
	qu.Push(design)
	trie := util.Trie[string]{}
	for !qu.IsEmpty() {
		cur, _ := qu.Pop()

		trieQu := util.NewQueue[*util.TrieNode[string]]()
		for _, root := range trie.Roots {
			if root.Value == cur[0] {
				trieQu.Push(root)
			}
		}

		for !trieQu.IsEmpty() {
			node, _ := trieQu.Pop()
			if node.IsLeafNode {
				prefix := node.FullValue
				remain, ok := strings.CutPrefix(cur, prefix)
				if !ok {
					panic(fmt.Sprintf("string %s does not contain prefix %s",
						cur, prefix))
				}
				qu.Push(remain)
			}

			for _, next := range node.Next {
				if next.Value != cur[0] {
					continue
				}
				trieQu.Push(next)
			}
		}
	}
	return false
}
