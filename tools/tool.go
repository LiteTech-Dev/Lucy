package tools

func Trenary[T any](expr func() bool, a T, b T) T {
	if expr() {
		return a
	}
	return b
}
