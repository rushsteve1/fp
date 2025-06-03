package fun

import "github.com/rushsteve1/fp"

// Predicate takes a value and returns a bool
type Predicate[T any] = func(T) bool

// Identity returns what it is passed
func Identity[T any](t T) T {
	return t
}

// Errorless takes a function that can error and returns a new function that
// wraps the provided one in [fp.Must]
// It's better to use [monads.FuncWrap]
func Errorless[T, U any](f func(T) (U, error)) func(T) U {
	return func(t T) U {
		return fp.Must(f(t))
	}
}

// Discard takes a function and returns a new wrapping function without the return value
func Discard[T, U any](f func(T) U) func(T) {
	return func(t T) {
		_ = f(t)
	}
}
