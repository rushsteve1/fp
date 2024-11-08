package monads

import "github.com/rushsteve1/fp"

// Monad is a context that an operation took place in, resulting in a value.
// They can be thought of as single-value sequences.
// You can apply additional operations within the monad using transducers.
//
// I know there's a lot of discussion about
type Monad[T any] interface {
	// All monads are sequences and can be manipulated using transducers
	fp.Seq[T]
	// All monads can be "unwrapped" with Get
	Gettable[T]
	// Ok returns true if this monad is "ok"
	// The actual meaning of that value depends on the monad
	Ok() bool
}
