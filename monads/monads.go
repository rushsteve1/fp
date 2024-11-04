package monads

import "github.com/rushsteve1/fp"

// Monad is a context that an operation took place in.
// You can apply additional operations to
type Monad[T any] interface {
	fp.Seq[T]
	Ok() bool
	Get() (T, error)
}
