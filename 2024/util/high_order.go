package util

func MapFunc[T, S any](slice []T, f func(T) S) []S {
	ret := make([]S, len(slice))
	for i := range slice {
		ret[i] = f(slice[i])
	}
	return ret
}

func UnfoldSlice(input []string, transformer func(string) []string) []string {
	result := []string{}
	for i := range input {
		result = append(result, transformer(input[i])...)
	}
	return result
}

func FoldFunc[T any](sl []T, folder func(a, b T) T) T {
	var ret T
	if len(sl) == 0 {
		return ret
	}
	ret = sl[0]
	for i := 1; i < len(sl); i++ {
		ret = folder(ret, sl[i])
	}
	return ret
}

func Add[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~int | ~int8 | ~int16 | ~int32 | ~int64 | string](a, b T) T {
	return a + b
}
