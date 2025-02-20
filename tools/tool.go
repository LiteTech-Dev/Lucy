package tools

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/russross/blackfriday/v2"
)

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

// Insert inserts a value into a slice at a slice[pos]. If the pos is out of
// bounds, the slice remains unchanged.
func Insert[T any](slice []T, pos int, value ...T) []T {
	if pos < 0 || pos > len(slice) {
		return slice
	}
	return append(slice[:pos], append(value, slice[pos:]...)...)
}

// CloseReader closes a reader and warns if an error occurs. Call this with a
// defer statement.
func CloseReader(reader io.ReadCloser, failAction func(error)) {
	err := reader.Close()
	if err != nil {
		failAction(err)
	}
}

const (
	NetworkTestTimeout = 5 // seconds
	NetworkTestRetries = 3
)

// NetworkTest is a simple the network connection test. You can use this before
// any operation that strictly requires a network connection.
//
// A nil value means the connection is successful.
func NetworkTest() (err error) {
	retry := NetworkTestRetries
	client := http.Client{
		Timeout: NetworkTestTimeout * time.Second,
	}
Retry:
	_, err = client.Get("https://example.com")
	if err != nil {
		retry--
		if retry > 0 {
			goto Retry
		}
		return err
	}
	return nil
}

func MarkdownToPlainText(md string) (s string) {
	s = string(blackfriday.Run([]byte(md)))
	return
}

// Decorate applies a series of decorators to a function. This is used to
// prevent nested function calls for better readability.
func Decorate[T interface{}](f T, decorators ...func(T) T) T {
	for _, decorator := range decorators {
		f = decorator(f)
	}
	return f
}
