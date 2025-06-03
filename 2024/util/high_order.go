package util

func MapFunc[T any](s []T, f func(T) T) []T {
	for i := range s {
		s[i] = f(s[i])
	}
	return s
}

func UnfoldSlice(input []string, transformer func(string) []string) []string {
	result := []string{}
	for i := range input {
		result = append(result, transformer(input[i])...)
	}
	return result
}
