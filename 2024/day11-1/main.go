package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

var distinct = map[int64]int{}

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

	scanner.Scan()
	line := scanner.Text()
	split := strings.Split(line, " ")

	var head *Node
	for i := range split {
		num, err := strconv.Atoi(split[i])
		if err != nil {
			panic(err)
		}
		InsertTail(&head, int64(num))
	}

	fmt.Println(length(head))

	start := time.Now()

	limit := 60
	counter := 1
	cur := head

	for cur != nil {
		for cur.Blink < limit {
			if blinkCount := blinkUntilSplit(cur, limit); blinkCount > 0 {
				counter++
			}
		}
		for cur != nil && cur.Blink == limit {
			cur = cur.Next
		}
	}

	fmt.Println(time.Since(start))
	// fmt.Println(len(distinct))
	// fmt.Println(distinct)
	// printLinkedList(head)
	// start = time.Now()
	fmt.Println(length(head))
	// fmt.Println("length calculate", time.Since(start))
	// lenghtList := length(head)
	fmt.Println(counter)
	// fmt.Println(lenghtList, counter, lenghtList == counter)
}

func length(cur *Node) int {
	counter := 0
	for cur != nil {
		distinct[cur.Data]++
		cur = cur.Next
		counter++
	}
	return counter
}

func printLinkedList(cur *Node) {
	for cur != nil {
		fmt.Printf("%d+%d ", cur.Data, cur.Blink)
		cur = cur.Next
	}
	fmt.Println()
}

func eval(cur *Node) {
	for cur != nil {
		if cur.Data == 0 {
			cur.Data = 1
			cur = cur.Next
		} else if split := splitDigits3(cur.Data); len(split) == 2 {
			cur.Data = split[0]
			InsertAfter(cur, split[1])
			cur = cur.Next
			cur = cur.Next
		} else {
			cur.Data *= 2024
			cur = cur.Next
		}
	}
}

func blinkOnce(cur *Node) bool {
	cur.Blink++
	if cur.Data == 0 {
		cur.Data = 1
		return false
	} else if split := splitDigits3(cur.Data); len(split) == 2 {
		cur.Data = split[0]
		InsertAfter(cur, split[1])
		return true
	} else {
		cur.Data *= 2024
		return false
	}
}

func blinkUntilSplit(cur *Node, limit int) int {
	// if return negative, outer loop should skip the node
	if entry, found := blinkCache[cur.Data]; found {
		// even if we blink, it won't split before we hit limit
		// just set current node to max limit and skip
		nextBlink := cur.Blink + entry.Blink
		if nextBlink > limit {
			cur.Blink = limit
			return -limit
		}

		cur.Blink = nextBlink
		cur.Data = entry.Split[0]

		InsertAfter(cur, entry.Split[1])
		return entry.Blink
	}

	prev := cur.Blink
	ogData := cur.Data
	for cur.Blink < limit {
		cur.Blink++
		if cur.Data == 0 {
			cur.Data = 1
		} else if split := splitDigits3(cur.Data); len(split) == 2 {
			cur.Data = split[0]
			InsertAfter(cur, split[1])
			blinkCache[ogData] = blinkCacheEntry{Blink: cur.Blink - prev, Split: split}
			return cur.Blink - prev
		} else {
			cur.Data *= 2024
		}
	}
	return -(cur.Blink - prev)
}

var blinkCache = map[int64]blinkCacheEntry{}

type blinkCacheEntry struct {
	Blink int
	Split []int64
}

type Node struct {
	Data  int64
	Next  *Node
	Blink int
}

func InsertAfter(target *Node, newData int64) {
	// if split, child node should have same blink number
	newNode := &Node{Data: newData, Blink: target.Blink}
	newNode.Next = target.Next
	target.Next = newNode
}

func InsertTail(head **Node, newData int64) {
	if *head == nil {
		*head = &Node{Data: newData}
		return
	}
	cur := *head
	for cur.Next != nil {
		cur = cur.Next
	}

	cur.Next = &Node{Data: newData}
}

func splitDigits(n int64) []int64 {
	strNum := strconv.Itoa(int(math.Abs(float64(n))))
	digitCount := len(strNum)
	if digitCount%2 != 0 {
		return nil
	}

	mid := len(strNum) / 2

	// Convert the two halves back to integers
	firstHalf, _ := strconv.ParseInt(strNum[:mid], 10, 64)
	secondHalf, _ := strconv.ParseInt(strNum[mid:], 10, 64)
	return []int64{firstHalf, secondHalf}
}

func splitDigits2(n int64) []int64 {
	digitCount := int(math.Log10(float64(n))) + 1
	if digitCount%2 != 0 {
		return nil
	}

	mid := digitCount / 2

	divisor := int64(math.Pow(10, float64(mid)))
	return []int64{n / divisor, n % divisor}
}

func splitDigits3(n int64) []int64 {
	count := 0
	cp := n
	for cp > 0 {
		cp /= 10
		count++
	}

	if count%2 != 0 {
		return nil
	}
	mid := count / 2

	divisor := int64(1)
	for mid > 0 {
		divisor *= 10
		mid--
	}
	return []int64{n / divisor, n % divisor}
}
