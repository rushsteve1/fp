package monads

import (
)

type Monad[T any] interface {
	Get() (T, error)
	Seq(yield func(T) bool)
}
