package util

import (
	"fmt"
	"sort"
	"strings"
)

type RadixTree struct {
	Children   map[string]*RadixTree
	IsLeafNode bool
	// if IsLeafNode == true, FullValue is the full value of the
	FullString string
}

func (root *RadixTree) Insert(input string, fullStr string) *RadixTree {
	if root.Children == nil {
		root.Children = map[string]*RadixTree{}
	}
	if input == "" {
		root.Children[""] = &RadixTree{
			IsLeafNode: true,
			FullString: "",
		}
		return root
	}

	prefix, key, found := root.findLongestCommonPrefixChild(input)

	switch {
	case !found:
		root.Children[input] = &RadixTree{
			Children:   map[string]*RadixTree{},
			IsLeafNode: true,
			FullString: fullStr,
		}
	case len(prefix) < len(input):
		// input needs to be split, the remnant of this is recursively insert into child node
		inputRemnant := strings.TrimPrefix(input, prefix)

		// if the key is longer than prefix, it needs to be split first
		// root -"abc"-> child
		// root -"ab"-> newchild -"c"->child
		if len(key) > len(prefix) {
			keyRemnant := strings.TrimPrefix(key, prefix)

			ogChild := root.Children[key]
			delete(root.Children, key)

			root.Children[prefix] = &RadixTree{
				Children: map[string]*RadixTree{
					keyRemnant: ogChild,
				},
			}
		}
		root.Children[prefix].Insert(inputRemnant, fullStr)

		// now the key is equal to prefix, recursively insert input to the
		// node pointed by the key
	case len(prefix) == len(input):
		// input does not need to be split, since the input is less or equal to
		// the key, we have a new leaf node
		//
		// if the key is longer than prefix, then it needs to be split first,
		// the child of this key becomes a grandchild, and this key will get a new child marked as a leaf node

		// root -"abc"-> child
		// root -"ab"-> newchild -"c"->child
		if len(key) > len(prefix) {
			keyRemnant := strings.TrimPrefix(key, prefix)
			ogChild := root.Children[key]
			delete(root.Children, key)
			root.Children[prefix] = &RadixTree{
				Children: map[string]*RadixTree{
					keyRemnant: ogChild,
				},
			}
		}

		root.Children[prefix].IsLeafNode = true
		root.Children[prefix].FullString = fullStr

	}

	return root
}

func (root *RadixTree) CanCompose(str string) bool {
	type innerQData struct {
		Node      *RadixTree
		Remaining string
	}
	outerQ := NewQueue[string]()
	outerQ.Push(str)
	for !outerQ.IsEmpty() {
		outerStr, _ := outerQ.Pop()
		// fmt.Println("outerQ size", outerQ.Size())

		innerQ := NewQueue[innerQData]()

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

func (root *RadixTree) findLongestCommonPrefixChild(input string) (string, string, bool) {
	var found bool
	var maxPf string
	var maxPfKey string

	for key := range root.Children {
		pf, length := longestCommonPrefix(input, key)
		if length > len(maxPf) {
			maxPf, maxPfKey, found = pf, key, true
		}
	}
	return maxPf, maxPfKey, found
}

func longestCommonPrefix(a, b string) (string, int) {
	minLen := min(len(a), len(b))
	commonPrefix := strings.Builder{}
	i := 0
	for ; i < minLen; i++ {
		if a[i] == b[i] {
			commonPrefix.WriteByte(a[i])
		} else {
			break
		}
	}

	return commonPrefix.String(), i
}

func PrintRadixTree(node *RadixTree, name string, prefix string, isLast bool) {
	connector := "├── "
	nextPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	leafMarker := ""
	if node.IsLeafNode {
		leafMarker = fmt.Sprintf(" (leaf) fs=%s", node.FullString)
	}

	fmt.Printf("%s%s%s%s\n", prefix, connector, name, leafMarker)

	// Get sorted keys for consistent ordering
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		isLastChild := i == len(keys)-1
		PrintRadixTree(node.Children[k], k, nextPrefix, isLastChild)
	}
}
