package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"text/tabwriter"

	"aoc2024/util"
)

type Pair = util.Pair

// - a (leaf)
//    - b (leaf)
//       - c (leaf) fs: [ c, bc abc ]
// a, b, c
// a, bc
// abc

// ab cd

// a bc d
// a b c d
// input: abc

// NOTE: now this one is a bit harder, we are not allowed to prune composition of string
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

	total := 0
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, design := range designs {
		result1 := BackTrackingTest(towels, design, map[string]int{})
		// result2, detail := BruteForceTest(towels, design)
		// fmt.Fprintf(w, "%s\t%t\t%d\t%d\n", design, result1 == result2, result1, result2)
		total += result1
	}
	w.Flush()
	fmt.Println(total)
}

func BackTrackingTest(towels []string, design string, cache map[string]int) int {
	total := 0
	if design == "" {
		return 1
	}
	if count, found := cache[design]; found {
		return count
	}
	for _, word := range towels {
		// candidate to remove
		if afterCut, found := strings.CutSuffix(design, word); found {
			total += BackTrackingTest(towels, afterCut, cache)
		}
	}

	cache[design] = total
	return total
}

func BruteForceTest(towels []string, design string) (int, [][]string) {
	type data struct {
		Segments []string
		Current  string
	}
	qu := util.NewQueue[data]()
	counter := 0
	qu.Push(data{Segments: []string{}, Current: design})
	detail := [][]string{}
	for !qu.IsEmpty() {
		top, _ := qu.Pop()
		for _, prefix := range towels {
			if remain, found := strings.CutPrefix(top.Current, prefix); found {
				if remain == "" {
					result := append(top.Segments, prefix)
					detail = append(detail, result)
					counter++
				} else {
					qu.Push(data{Segments: append(top.Segments, prefix), Current: remain})
				}
			}
		}
	}

	return counter, detail
}

func RadixTreeTest(towels []string, design string) int {
	trie := &util.RadixTree{}
	for _, str := range towels {
		if str == "" {
			panic("empty string towel")
		}
		trie.Insert(str, str)
	}

	return test(design, trie)
}

func test(design string, root *util.RadixTree) int {
	type innerQData struct {
		Node      *util.RadixTree
		Remaining string
	}

	type outerData struct {
		Previous string
		Current  string
	}

	outerQ := util.NewQueue[outerData]()
	outerQ.Push(outerData{
		Previous: "<start>",
		Current:  design,
	})
	innerQ := util.NewQueue[innerQData]()
	cache := map[string]int{
		design: 1,
	}

	for !outerQ.IsEmpty() {
		outer, _ := outerQ.Pop()
		if outer.Current == "ugrrrrugrrrr" {
			fmt.Printf("1st half ended: %+v\n", cache)
		}

		innerQ.Push(innerQData{
			Node:      root,
			Remaining: outer.Current,
		})

		for !innerQ.IsEmpty() {
			next, _ := innerQ.Pop()
			node, innerStr := next.Node, next.Remaining
			if node.IsLeafNode {
				if _, found := cache[innerStr]; found {
					fmt.Printf("inner=%s, outer=%s\n", innerStr, outer.Current)
					cache[innerStr] += cache[outer.Current]
				} else {
					outerQ.Push(outerData{
						Previous: outer.Current,
						Current:  innerStr,
					})
					cache[innerStr] = cache[outer.Current]
				}
			}

			for key, child := range node.Children {
				if remain, found := strings.CutPrefix(innerStr, key); found {
					fmt.Printf("cutting from=%s, cut=%s, to=%s\n", innerStr, key, remain)
					innerQ.Push(innerQData{
						Node:      child,
						Remaining: remain,
					})
				}
			}
		}
	}
	fmt.Printf("cache %+v\n", cache)
	return cache[""]
}

func genTestFromTowels(towels []string, maxLen int) string {
	output := ""

	for len(output) < maxLen {
		output += towels[rand.Intn(len(towels))]
	}
	return output
}

// more randomized test that generate from list of alphabet and not just from
// list of words
func genTestFromAlphabet(towels []string, maxLen int) string {
	output := ""

	alphabet := []string{"r", "b", "w", "g", "u"}
	for len(output) < maxLen {
		output += towels[rand.Intn(len(alphabet))]
	}
	return output
}
