package util

func Trenary(expr func() bool, a any, b any) any {
	if expr() {
		return a
	}
	return b
}
