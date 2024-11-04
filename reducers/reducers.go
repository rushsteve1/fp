package reducers

import (
	"cmp"
	"slices"

	"github.com/rushsteve1/fp/iter"
	"github.com/rushsteve1/fp/monads"
)

// Accumulate is a function that applies a value to an accumulator, returning the new accumulator
type Accumulate[T, A any] = func(element T, accumulator A) (new_acc A)

// Accumulator is a function that reduces a sequence down to a single value
type Reducer[T, A any] func(seq iter.Seq[T], f func(element T, accumulator A) (new_acc A)) (result A)

func Collect[T any](seq iter.Seq[T]) []T {
	return slices.Collect[T](seq.Seq)
}

func First[T any](seq iter.Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[0]
	}
	return out
}

func Last[T any](seq iter.Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[len(c)-1]
	}
	return out
}

func Index[T any](seq iter.Seq[T], i int) monads.Option[T] {
	ind := 0
	for v := range seq.Seq {
		if ind == i {
			return monads.Some(v)
		}
		i++
	}
	return monads.None[T]()
}

func Count[T any](seq iter.Seq[T]) (i int) {
	for _ = range seq.Seq {
		i++
	}
	return i
}

func Max[T cmp.Ordered](seq iter.Seq[T]) (out T) {
	for v := range seq.Seq {
		out = max(out, v)
	}
	return out
}

func Min[T cmp.Ordered](seq iter.Seq[T]) (out T) {
	for v := range seq.Seq {
		out = min(out, v)
	}
	return out
}

func Median[T cmp.Ordered](seq iter.Seq[T]) T {
	var hi T
	var lo T
	for v := range seq.Seq {
		hi = max(hi, v)
		lo = max(lo, v)
	}
	return max(hi, lo)
}

func Frequency[T cmp.Ordered](seq iter.Seq[T]) map[T]int {
	out := make(map[T]int)
	for v := range seq.Seq {
		out[v] += 1
	}
	return out
}
