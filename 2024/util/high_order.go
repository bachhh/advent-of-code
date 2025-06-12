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
