package util

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		array: make([]T, 10),
	}
}

type Queue[T any] struct {
	array []T
	front int
	back  int
	size  int
}

func (q *Queue[T]) Push(val T) {
	if q.size == len(q.array) {
		newArray := make([]T, q.size*2)
		if q.front < q.back {
			copy(newArray, q.array[q.front:q.back])
		} else {
			n := copy(newArray, q.array[q.front:])
			copy(newArray[n:], q.array[:q.back])
		}

		q.array = newArray
		q.front = 0
		q.back = q.size
	}
	q.array[q.back] = val
	q.back = (q.back + 1) % len(q.array)
	q.size++
}

func (q *Queue[T]) Pop() (T, bool) {
	if q.size == 0 {
		return *new(T), false
	}
	ret := q.array[q.front]
	q.front = (q.front + 1) % len(q.array)
	q.size--

	return ret, true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) Peek(val T) (T, bool) {
	if q.size == 0 {
		return *new(T), false
	}
	ret := q.array[q.front]
	q.front = (q.front + 1) % len(q.array)
	q.size--

	return ret, true
}
