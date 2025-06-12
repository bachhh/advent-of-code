package util

import "fmt"

type TreeNode[T any] struct {
	Parent   *TreeNode[T]
	Children []*TreeNode[T]
	Value    T
}

func (n *TreeNode[T]) PrintTree(depthLimit int) {
	fmt.Printf("%v\n", n.Value) // root node with no prefix
	for i, child := range n.Children {
		child.print("", i == len(n.Children)-1, depthLimit-1)
	}
}

func (n *TreeNode[T]) print(prefix string, isLast bool, depthLimit int) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	fmt.Printf("%s%s%v", prefix, connector, n.Value)
	if depthLimit == 0 {
		if len(n.Children) > 0 {
			fmt.Printf(" !( %d children hidden)", len(n.Children))
		}
		fmt.Println()
		return
	}
	fmt.Println()

	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	for i, child := range n.Children {
		child.print(newPrefix, i == len(n.Children)-1, depthLimit-1)
	}
}

func (node *TreeNode[T]) FindSmallestLeaf(comp func(a, b T) int) *TreeNode[T] {
	if node == nil {
		return nil
	}

	var smallest *TreeNode[T]

	var dfs func(n *TreeNode[T])
	dfs = func(n *TreeNode[T]) {
		if len(n.Children) == 0 {
			// It's a leaf
			if smallest == nil || comp(n.Value, smallest.Value) < 0 {
				smallest = n
			}
			return
		}
		for _, child := range n.Children {
			dfs(child)
		}
	}

	dfs(node)
	return smallest
}

func (n *TreeNode[T]) AddChild(val T) *TreeNode[T] {
	child := &TreeNode[T]{Value: val}
	n.Children = append(n.Children, child)
	return child
}
