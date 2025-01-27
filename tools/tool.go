package tools

import "sync"

// TernaryFunc gives a if expr == true, b if expr == false. For a simple
// bool expression, use Ternary instead.
func TernaryFunc[T any](expr func() bool, a T, b T) T {
	if expr() {
		return a
	}
	return b
}

// Ternary gives a if v == true, b if v == false. For a function parameter, use
// TernaryFunc instead.
func Ternary[T any](v bool, a T, b T) T {
	if v {
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
