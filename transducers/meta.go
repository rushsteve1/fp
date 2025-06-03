package transducers

import (
	. "github.com/rushsteve1/fp"
)

// Meta-transducers are functions that take or return a transducer.
// They often can be implemented as standard transducers but are simpler this way

// Visitor takes a transducer and returns a new transducer that
// counteracts the first and returns the original sequence unaltered.
// Useful for transducers with side-effcts like [Write].
// This is a higher-kinded version of [Each].
//
// This took me a very long time to figure out
func Visitor[T, U any](tx Transducer[T, U]) Transducer[T, T] {
	// Create tne new transducer that takes the sequence
	return func(seq Seq[T]) Seq[T] {
		// Return a new sequence
		return SeqFunc[T](func(y1 func(T) bool) {
			// Create another new sequence and pass that to tx
			tx(SeqFunc[T](func(y2 func(T) bool) {
				// Call the passed in sequence
				seq.Seq(func(t T) bool {
					// Yield to both new sequences
					return y1(t) && y2(t)
				})
			}))
		})
	}
}
