package monads

import "github.com/rushsteve1/fp"

// Monad is a context that an operation took place in.
// You can apply additional operations within the monad using transducers
type Monad[T any] interface {
	// All monads are sequences and can be manipulated using transducers
	fp.Seq[T]
	// Ok returns true if this monad is "ok"
	// The actual meaning of that value depends on the monad
	Ok() bool
	// Get tries to retrieve the inner value of the monad
	Get() (T, error)
}
