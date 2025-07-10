package util

type DisjointSet[T comparable] struct {
	parent map[T]T
	size   map[T]int
}

func NewDisjointSet[T comparable]() *DisjointSet[T] {
	return &DisjointSet[T]{
		parent: map[T]T{},
		size:   map[T]int{},
	}
}

func (d *DisjointSet[T]) MakeSet(x T) {
	d.parent[x] = x
	d.size[x] = 1
}

func (d *DisjointSet[T]) FindOrMake(x T) T {
	if parent, found := d.Find(x); found {
		return parent
	}
	d.MakeSet(x)
	return x
}

func (d *DisjointSet[T]) Find(x T) (T, bool) {
	if _, found := d.parent[x]; !found {
		var ret T
		return ret, false
	}

	if d.parent[x] == x {
		return d.parent[x], true
	}
	d.parent[x], _ = d.Find(d.parent[x])
	return d.parent[x], true
}

func (d *DisjointSet[T]) Union(a, b T) {
	a = d.FindOrMake(a)
	b = d.FindOrMake(b)
	if a != b {
		if d.size[a] < d.size[b] {
			a, b = b, a
		}
		d.parent[b] = a
		d.size[a] += d.size[b]
	}
}

func (d *DisjointSet[T]) ToMap() map[T][]T {
	m := map[T][]T{}

	for child := range d.parent {
		parent, _ := d.Find(child)
		m[parent] = append(m[parent], child)
	}
	return m
}

func (d *DisjointSet[T]) ToSlice() [][]T {
	m := d.ToMap()
	ret := [][]T{}
	for _, sl := range m {
		ret = append(ret, sl)
	}
	return ret
}
