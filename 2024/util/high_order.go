package util

func MapFunc[T any](s []T, f func(T) T) []T {
	for i := range s {
		s[i] = f(s[i])
	}
	return s
}
