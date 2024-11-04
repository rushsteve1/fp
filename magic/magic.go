// Magic is a package that provides a lot of things that maybe Go doesn't need
// but with generics are now possible.
//
// If you love the functional programming style you'll like this.
//
// But be warned... here be cursed things.

// Guidelines and ideas
// 1. Generics are great
// 2. Reflection is OK
// 3. Variadic parameters allow for optional arguments
// 4. Type assertions are handy

package magic

import (
	"cmp"
)

type Nilable interface {
	~*any | ~[]any | ~map[any]any
}

var GlobalErrorHandler = func(err error) {
	panic(err)
}

func Check(err error) {
	if err != nil {
		GlobalErrorHandler(err)
	}
}

// Must is the first function anyone wants in Go
func Must[T any](t T, err error) T {
	Check(err)
	return t
}

func Clamp[T cmp.Ordered](x T, lo T, hi T) T {
	return max(min(x, hi), lo)
}

func CastOr[T any](x any, or T) T {
	if cast, ok := x.(T); ok {
		return cast
	}
	return or
}

// Or re-exports [cmp.Or]
func Or[T comparable](vals ...T) T {
	return cmp.Or[T](vals...)
}

func Ternary[T any](cond bool, a T, b T) T {
	if cond {
		return a
	}
	return b
}
