// Implements Transducers for [Seq]

package transducers

import (
	. "iter"
	"reflect"
	"time"

	"github.com/rushsteve1/fp/magic"
	"github.com/rushsteve1/fp/reducers"
)

// Most of these type definitions are for illustrative purposes and are unnecessary
// Also because generic type alaises isn't in yet

// Transform takes a value and returns a new value.
type Transform[T, U any] func(T) U

// Transducer is a generalized mapping of a computation between two Sequences.
// The easiest way to create a transducer is using [magic.Curry2] on a [Transform]
type Transducer[T, U any] func(Seq[T]) Seq[U]

func Transduce[T, U, V any](tx Transducer[T, U], rx reducers.Collector[U, V], src Seq[T]) V {
	return rx(tx(src))
}

func Pass[T any](seq Seq[T]) Seq[T] {
	return seq
}

// Map is the simplest but shows how it all actually works the same as a transducer
func Map[T, U any](seq Seq[T], f func(T) U) Seq[U] {
	return func(yield func(U) bool) {
		seq(func(t T) bool {
			return yield(f(t))
		})
	}
}

// Filter has the added constraint [comparable]
func Filter[T comparable](seq Seq[T], f func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(t T) bool {
			return magic.Ternary(f(t), yield(t), true)
		})
	}
}

// Visit can be trivially defined using [Map]
func Visit[T any](seq Seq[T], f func(T)) Seq[T] {
	return Map[T, T](seq, func(t T) T {
		f(t)
		return t
	})
}

// Take returns a new iterator that stops after count elements
func Take[T any](seq Seq[T], count int) Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		seq(func(t T) bool {
			if i >= count {
				return false
			}
			i++
			return yield(t)
		})
	}
}

// Drop removes the first count elements from the sequence
func Drop[T any](seq Seq[T], count int) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(t T) bool {
			for range count {
				if !yield(t) {
					return false
				}
			}
			return yield(t)
		})
	}
}

// Fuse stops the sequence at the first nil value
func Fuse[T comparable](seq Seq[T], count int) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(t T) bool {
			// Generics fail us hereabouts
			if reflect.ValueOf(t).IsNil() {
				return false
			}
			return yield(t)
		})
	}
}

func Debounce[T any](seq Seq[T], delay time.Duration) Seq[T] {
	return func(yield func(T) bool) {
		last := time.Unix(0, 0).UTC()
		seq(func(t T) bool {
			if time.Since(last) > delay {
				last = time.Now().UTC()
				return yield(t)
			}
			// Skip the debounced elements
			return true
		})
	}
}

func Delta[T magic.Numeric](seq Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		var prev *T
		seq(func(t T) bool {
			if prev == nil {
				prev = &t
				return true
			}
			return yield(t - *prev)
		})
	}
}

func TimeDelta(seq Seq[time.Time]) Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		var prev *time.Time = nil
		seq(func(t time.Time) bool {
			if prev == nil {
				prev = &t
				return true
			}
			halt := yield(t.Sub(*prev))
			prev = &t
			return halt
		})
	}
}
