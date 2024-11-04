package reducers

import (
	"cmp"
	"slices"

	. "iter"

	"github.com/rushsteve1/fp/magic"
	"github.com/rushsteve1/fp/monads"
)

// Accumulate is a function that applies a value to an accumulator, returning the new accumulator
type Accumulate[T, A any] func(T, A) A

// Reducer is a function that reduces a sequence down to a single value
type Reducer[T, A any] func(Seq[T], Accumulate[T, A]) A

// Collector takes a sequence and returns a single value.
// [Reducer] can be converted to Collector using [Curry]
type Collector[T, A any] func(Seq[T]) A

func Collect[T any](seq Seq[T]) []T {
	return slices.Collect[T](seq)
}

func Reduce[T, A any](seq Seq[T], a A, f Accumulate[T, A]) A {
	for v := range seq {
		a = f(v, a)
	}
	return a
}

func First[T any](seq Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[0]
	}
	return out
}

func Last[T any](seq Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[len(c)-1]
	}
	return out
}

func Index[T any](seq Seq[T], i int) monads.Option[T] {
	ind := 0
	for v := range seq {
		if ind == i {
			return monads.Some(v)
		}
		i++
	}
	return monads.None[T]()
}

func Count[T any](seq Seq[T]) (i int) {
	for _ = range seq {
		i++
	}
	return i
}

func Max[T cmp.Ordered](seq Seq[T]) (out T) {
	for v := range seq {
		out = max(out, v)
	}
	return out
}

func Min[T cmp.Ordered](seq Seq[T]) (out T) {
	for v := range seq {
		out = min(out, v)
	}
	return out
}

func Median[T cmp.Ordered](seq Seq[T]) T {
	var hi T
	var lo T
	for v := range seq {
		hi = max(hi, v)
		lo = max(lo, v)
	}
	return max(hi, lo)
}

func Frequency[T cmp.Ordered](seq Seq[T]) map[T]int {
	out := make(map[T]int)
	for v := range seq {
		out[v] += 1
	}
	return out
}

func Average[T magic.Numeric](seq Seq[T]) T {
	count := 0
	var sum T
	for v := range seq {
		count++
		sum += v
	}
	return sum / T(count)
}

func Any[T any](seq Seq[T], f func(T) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}
	return false
}

func All[T any](seq Seq[T], f func(T) bool) bool {
	for v := range seq {
		if !f(v) {
			return false
		}
	}
	return true
}
