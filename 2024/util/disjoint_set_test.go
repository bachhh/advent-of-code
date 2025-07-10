package util_test

import (
	"fmt"
	"testing"

	"aoc2024/util"

	"github.com/stretchr/testify/require"
)

func find[T comparable](ds *util.DisjointSet[T], input T) T {
	set, _ := ds.Find(input)
	return set
}

func TestDisjointSet(t *testing.T) {
	ds := util.NewDisjointSet[string]()
	ds.Union("a", "b")
	ds.Union("c", "d")
	ds.Union("b", "c")
	aset := find(ds, "a")
	testcase := []string{"b", "c", "d"}
	for _, actual := range testcase {
		require.Equal(t, find(ds, actual), aset)
	}
}

func TestDisjointSet2(t *testing.T) {
	steps := [][]string{
		{"kh", "tc"},
		{"qp", "kh"},
		{"de", "cg"},
		{"de", "cg"},
		{"ka", "co"},
		{"yn", "aq"},
		{"qp", "ub"},
		{"cg", "tb"},
		{"vc", "aq"},
		{"tb", "ka"},
		{"wh", "tc"},
		{"yn", "cg"},
	}

	ds := util.NewDisjointSet[string]()
	for _, step := range steps {
		a, b := step[0], step[1]
		ds.Union(a, b)
		fmt.Println(a, b, ds.ToSlice())
	}

	require.Equal(t, find[string](ds, "vc"), find[string](ds, "yn"))
	fmt.Println(ds.ToSlice())
}
