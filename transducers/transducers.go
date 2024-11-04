// Implements Transducers for [Seq]

package transducers

import (
	"io"
	"time"

	. "github.com/rushsteve1/fp"
	"github.com/rushsteve1/fp/reducers"
)

// Most of these type definitions are for illustrative purposes and are unnecessary
// Also because generic type alaises isn't in yet we can't actually use them much

// Transform takes a value and returns a new value
type Transform[T, U any] func(T) U

// Predicate takes a value and returns a bool
type Predicate[T any] func(T) bool

// Transducer is a generalized mapping of a computation between two Sequences.
// It is higher-kinded than a normal HO transform.
// Rust's iter and Elixir's Stream are transducers, but Elixir's Enum is not.
// The easiest way to create a transducer is using [threading.Curry2].
type Transducer[T, U any] func(Seq[T]) Seq[U]

// Transduce is the main event, the rest of this library exists to support it.
// It allows you to chain complex calculations into a single sequence
// and then reduce that to a single value.
func Transduce[T, U, V any](tx Transducer[T, U], rx reducers.Collector[U, V], src Seq[T]) V {
	return rx(tx(src))
}

// Map is the simplest, but shows how it all actually works
func Map[T, U any](seq Seq[T], f Transform[T, U]) Seq[U] {
	return SeqFunc[U](func(yield func(U) bool) {
		seq.Seq(func(t T) bool {
			return yield(f(t))
		})
	})
}

// Filter has the added constraint [comparable] but only needs one generic
func Filter[T comparable](seq Seq[T], f Predicate[T]) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			return Ternary(f(t), yield(t), true)
		})
	})
}

// Each can be trivially defined using [Map]
func Each[T any](seq Seq[T], f func(T)) Seq[T] {
	return Map[T, T](seq, func(t T) T {
		f(t)
		return t
	})
}

// Take yields the first count elements in the sequence
func Take[T any](seq Seq[T], count int) Seq[T] {
	i := 0
	return TakeWhile(seq, func(t T) bool {
		o := i < count
		i++
		return o
	})
}

// TakeWhile yields elements while the predicate is true
func TakeWhile[T any](seq Seq[T], f Predicate[T]) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			if f(t) {
				return yield(t)
			}
			return false
		})
	})
}

// Drop removes the first count elements from the sequence
func Drop[T any](seq Seq[T], count int) Seq[T] {
	i := 0
	return DropWhile(seq, func(t T) bool {
		o := i < count
		i++
		return o
	})
}

// DropWhile removes elements while the predicate is true
func DropWhile[T any](seq Seq[T], f Predicate[T]) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		dropstop := false
		seq.Seq(func(t T) bool {
			if dropstop {
				return yield(t)
			}
			if f(t) {
				return true
			} else {
				dropstop = true
				return yield(t)
			}
		})
	})
}

// Append yields all the values of the first sequence, then the second
func Append[T any](seq Seq[T], next Seq[T]) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		firstdone := false
		if !firstdone {
			seq.Seq(func(t T) bool {
				firstdone = yield(t)
				return firstdone
			})
		} else {
			next.Seq(yield)
		}
	})
}

// Fuse stops the sequence at the first nil value
func Fuse[T Nilable](seq Seq[T], count int) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			if t == nil {
				return false
			}
			return yield(t)
		})
	})
}

// Debounce only yields values if the current element was yielded at least delay
// time since the last value was yielded.
// Elements that happen in-between debounces are dropped.
func Debounce[T any](seq Seq[T], delay time.Duration) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		last := time.Unix(0, 0).UTC()
		seq.Seq(func(t T) bool {
			if time.Since(last) > delay {
				last = time.Now().UTC()
				return yield(t)
			}
			// Skip the debounced elements
			return true
		})
	})
}

// Delta returns a new sequence that is the difference between adjacent elements
func Delta[T Numeric](seq Seq[T]) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		var prev *T
		seq.Seq(func(t T) bool {
			if prev == nil {
				prev = &t
				return true
			}
			return yield(t - *prev)
		})
	})
}

// TimeDelta is [Delta] but specialized for [time.Time]
func TimeDelta(seq Seq[time.Time]) Seq[time.Duration] {
	return SeqFunc[time.Duration](func(yield func(time.Duration) bool) {
		var prev *time.Time = nil
		seq.Seq(func(t time.Time) bool {
			if prev == nil {
				prev = &t
				return true
			}
			halt := yield(t.Sub(*prev))
			prev = &t
			return halt
		})
	})
}

// Enumerate returns a new [Seq2] with indices as keys
func Enumerate[T any](seq Seq[T]) Seq2[int, T] {
	return Seq2Func[int, T](func(yield func(int, T) bool) {
		i := 0
		seq.Seq(func(t T) bool {
			stop := yield(i, t)
			i++
			return stop
		})
	})
}

// Step only yields ever step elements
func Step[T any](seq Seq[T], step int) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		i := 1
		seq.Seq(func(t T) bool {
			if i%step == 0 {
				return yield(t)
			}
			i++
			return true
		})
	})
}

// Write will write to the given writer for every element.
// See it counterpart [generators.Reader]
func Write(seq Seq[[]byte], w io.Writer) Seq[error] {
	return Map(seq, func(b []byte) error {
		_, err := w.Write(b)
		return err
	})
}
