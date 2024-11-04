// Implements Transducers for [iter.Seq]

package transducers

import (
	"reflect"

	"github.com/rushsteve1/fp/iter"
	"github.com/rushsteve1/fp/magic"
	"github.com/rushsteve1/fp/reducers"
)

// Most of these type definitions are for illustrative purposes and are unnecessary
// Also because generic type alaises isn't in yet

// Transform takes a value and returns a new value.
type Transform[T, U any] = func(T) U

// Visit just looks at the value of an element without modifying it
type Visitor[T any] = Transform[T, T]

// Transducer is a generalized mapping of a computation between two iter.Sequences.
// The easiest way to create a transducer is using [magic.Curry2] on a [Transform]
type Transducer[T, U any] = func(iter.Seq[T]) iter.Seq[U]

func Transduce[T, U, V any](tx Transducer[T, U], rx reducers.Reducer[U, V], src iter.Seq[T]) V {
	return rx(tx(src))
}

func Pass[T any](seq iter.Seq[T]) iter.Seq[T] {
	return seq
}

// Map is the simplest but shows how it all actually works the same as a transducer
func Map[T, U any](seq iter.Seq[T], f Transform[T, U]) iter.Seq[U] {
	return iter.SeqFunc[U](func(yield func(U) bool) {
		seq.Seq(func(t T) bool {
			return yield(f(t))
		})
	})
}

// Filter has the added constraint [comparable]
func Filter[T comparable](seq iter.Seq[T], f Transform[T, bool]) iter.Seq[T] {
	return iter.SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			return magic.Ternary(f(t), yield(t), true)
		})
	})
}

// Visit can be trivially defined using [Map]
func Visit[T any](seq iter.Seq[T], f Visitor[T]) iter.Seq[T] {
	return Map[T, T](seq, func(t T) T {
		f(t)
		return t
	})
}

// Take returns a new iterator that stops after count elements
func Take[T any](seq iter.Seq[T], count int) iter.Seq[T] {
	return iter.SeqFunc[T](func(yield func(T) bool) {
		i := 0
		seq.Seq(func(t T) bool {
			if i >= count {
				return false
			}
			return yield(t)
		})
	})
}

// Drop removes the first count elements from the sequence
func Drop[T any](seq iter.Seq[T], count int) iter.Seq[T] {
	return iter.SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			for _ = range count {
				if !yield(t) {
					return false
				}
			}
			return yield(t)
		})
	})
}

// Fuse stops the sequence at the first nil value
func Fuse[T comparable](seq iter.Seq[T], count int) iter.Seq[T] {
	return iter.SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			// Generics fail us hereabouts
			if reflect.ValueOf(t).IsNil() {
				return false
			}
			return yield(t)
		})
	})
}
