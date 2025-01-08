package tools

import "sync"

// Ternary gives a if expr == true, b if expr == false
func Ternary[T any](expr func() bool, a T, b T) T {
	if expr() {
		return a
	}
	return b
}

// Memoize is only used for functions that do not take any arguments and return
// a value (typically a struct) that can be treated as a constant.
func Memoize[T any](f func() T) func() T {
	var result T
	var once sync.Once
	return func() T {
		once.Do(
			func() {
				result = f()
			},
		)
		return result
	}
}
