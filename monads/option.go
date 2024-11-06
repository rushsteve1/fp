package monads

import (
	"database/sql"
	"errors"

	"github.com/rushsteve1/fp"
)

var ErrUnwrapInvalid = errors.New("Unwrapped invalid Option")

// Option is a monad that wraps a value that may or may not exist.
// It is the same as in Rust.
type Option[T any] struct {
	// I have no idea why I bothered implementing it this way
	// but while we're re-interpeting the stdlib might as well
	sql.Null[T]
}

func (o Option[T]) Ok() bool {
	return o.Valid
}

func (o Option[T]) Get() (T, error) {
	return o.V, fp.Ternary(o.Valid, nil, ErrUnwrapInvalid)
}

func (o Option[T]) Seq(yield func(T) bool) {
	if o.Valid {
		yield(o.V)
	}
}

// Ptr is a helper function for converting an [Option]
func (o Option[T]) Ptr() *T {
	return fp.Ternary(o.Valid, &o.V, nil)
}

// Some returns a valid [Option]
func Some[T any](v T) Option[T] {
	return Option[T]{
		sql.Null[T]{
			V:     v,
			Valid: true,
		},
	}
}

// TrySome is like [Some] but takes in a pointer,
// so can be used with many existing functions.
// It is the logical inverse of [Ptr]
func TrySome[T any](v *T) Option[T] {
	if v == nil {
		return Option[T]{
			sql.Null[T]{
				V:     *v,
				Valid: true,
			},
		}
	}

	return Option[T]{
		sql.Null[T]{
			Valid: false,
		},
	}
}

func None[T any]() Option[T] {
	var v T
	return Option[T]{
		sql.Null[T]{
			V:     v,
			Valid: false,
		},
	}
}
