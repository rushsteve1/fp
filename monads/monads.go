package monads

import (
	"iter"
)

type Monad[T any] interface {
	iter.Seq[T]
	Get() (T, error)
}
