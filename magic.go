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

package fp

import (
	"cmp"
	"runtime"
	"slices"
	"testing"
)

var GlobalErrorHandler = func(err error) bool {
	panic(err)
}

// Must is the first function anyone wants in Go
func Must[T any](t T, err error) T {
	Check(err)
	return t
}

// Check is the second
func Check(err error) {
	if err != nil {
		if !GlobalErrorHandler(err) {
			runtime.Goexit()
		}
	}
}

// Ptr returns a pointer of its argument
func Ptr[T any](t T) *T {
	return &t
}

// DerefOr dereferences the passed pointer otherwise returns the or value
func DerefOr[T any](ref *T, or T) T {
	if ref != nil {
		return *ref
	}
	return or
}

// DerefZero is like [DerefOr] but returns the zero value if nil
func DerefZero[T any](ref *T) T {
	var t T
	return DerefOr(ref, t)
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

func Assert(t *testing.T, cond bool) {
	if !cond {
		t.Error("Assertion failure")
	}
}

func AssertEq[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Errorf("Assertion failure: %+v != %+v", a, b)
	}
}

func AssertSliceEq[S ~[]E, E comparable](t *testing.T, a, b S) {
	if !slices.Equal(a, b) {
		t.Errorf("Assertion failure: %+v != %+v", a, b)
	}
}
