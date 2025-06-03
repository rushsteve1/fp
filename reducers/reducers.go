package reducers

import (
	"cmp"
	"slices"

	. "github.com/rushsteve1/fp"
	"github.com/rushsteve1/fp/fun"
	"github.com/rushsteve1/fp/monads"
)

// Reduction is a function that applies a value to an accumulator, returning the new accumulator
type Reduction[T, Acc any] = func(T, Acc) Acc

// Reducer is a function that reduces a sequence down to a single value
type Reducer[T, Acc any] = func(Seq[T], Reduction[T, Acc]) Acc

// Collector takes a sequence and returns a single value.
// [Reducer] can be converted to Collector using [threading.Curry2]
type Collector[T, Acc any] = func(Seq[T]) Acc

// Consume is the simplest possible reducer
// It pulls every element in the sequence then discards them
func Consume[T any](seq Seq[T]) {
	for range seq.Seq {
	}
}

// Collect wraps [slices.Collect]
func Collect[T any](seq Seq[T]) []T {
	return slices.Collect[T](seq.Seq)
}

// Collect2 is the [fp.Seq2] version of [Collect]
func Collect2[K comparable, V any](seq Seq2[K, V]) map[K]V {
	out := make(map[K]V)
	for k, v := range seq.Seq2 {
		out[k] = v
	}
	return out
}

// Same as [Collect2]
func CollectKV[K comparable, V any](seq Seq[KeyValue[K, V]]) map[K]V {
	out := make(map[K]V)
	for kv := range seq.Seq {
		out[kv.Key] = kv.Value
	}
	return out
}

// Reduce consumes a sequence returning a final accumulator value
func Reduce[T, Acc any](seq Seq[T], a Acc, f Reduction[T, Acc]) Acc {
	for v := range seq.Seq {
		a = f(v, a)
	}
	return a
}

// First returns the first element of a sequence
func First[T any](seq Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[0]
	}
	return out
}

// Last returns the last element of a sequence, which may not exist
func Last[T any](seq Seq[T]) (out T) {
	c := Collect(seq)
	if len(c) > 0 {
		return c[len(c)-1]
	}
	return out
}

// Index returns the element at the given index, if it exists
func Index[T any](seq Seq[T], i int) monads.Option[T] {
	ind := 0
	for v := range seq.Seq {
		if ind == i {
			return monads.Some(v)
		}
		i++
	}
	return monads.None[T]()
}

// Length returns the number of elements in a sequence
func Length[T any](seq Seq[T]) (i int) {
	for _ = range seq.Seq {
		i++
	}
	return i
}

// Max returns the maximum value of a sequence, determined by [max]
func Max[T cmp.Ordered](seq Seq[T]) (out T) {
	for v := range seq.Seq {
		out = max(out, v)
	}
	return out
}

// Min returns the minimum value of a sequence, determined by [min]
func Min[T cmp.Ordered](seq Seq[T]) (out T) {
	for v := range seq.Seq {
		out = min(out, v)
	}
	return out
}

// Median returns the median value of the sequence
func Median[T cmp.Ordered](seq Seq[T]) T {
	var hi T
	var lo T
	for v := range seq.Seq {
		hi = max(hi, v, lo)
		lo = max(lo, v, hi)
	}
	return lo
}

// Frequency returns the frequency of each element in the sequence
func Frequency[T cmp.Ordered](seq Seq[T]) map[T]int {
	out := make(map[T]int)
	for v := range seq.Seq {
		out[v] += 1
	}
	return out
}

// Average returns the average of a numeric sequence
func Average[T Numeric](seq Seq[T]) T {
	count := 0
	var sum T
	for v := range seq.Seq {
		count++
		sum += v
	}
	return sum / T(count)
}

// Any returns true if any value in the sequence passes the predicate.
// Short-circuits on the first passing value
func Any[T any](seq Seq[T], f fun.Predicate[T]) bool {
	for v := range seq.Seq {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns true if all values in the sequence pass the predicate
func All[T any](seq Seq[T], f fun.Predicate[T]) bool {
	for v := range seq.Seq {
		if !f(v) {
			return false
		}
	}
	return true
}
