package util

import "fmt"

type TreeNode[T any] struct {
	Parent   *TreeNode[T]
	Children []*TreeNode[T]
	Value    T
}

func (n *TreeNode[T]) PrintTree() {
	fmt.Printf("%v\n", n.Value) // root node with no prefix
	for i, child := range n.Children {
		child.print("", i == len(n.Children)-1)
	}
}

func (n *TreeNode[T]) print(prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	fmt.Printf("%s%s%v\n", prefix, connector, n.Value)

	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	for i, child := range n.Children {
		child.print(newPrefix, i == len(n.Children)-1)
	}
}
