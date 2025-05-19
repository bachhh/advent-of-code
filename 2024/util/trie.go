package util

// prefix
type Trie[T string] struct {
	Roots []*TrieNode[T]
}

type TrieNode[T string] struct {
	Next       map[T]*TrieNode[T]
	Value      byte
	IsLeafNode bool
	// if IsLeafNode == true, FullValue is the full value of the
	// string
	FullValue T
}
