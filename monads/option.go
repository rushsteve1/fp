package monads

import (
	"database/sql"
	"errors"

	"github.com/rushsteve1/fp"
)

var ErrUnwrapInvalid = errors.New("Unwrapped invalid Option")

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

func (o Option[T]) Ptr() *T {
	return fp.Ternary(o.Valid, &o.V, nil)
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		sql.Null[T]{
			V:     v,
			Valid: true,
		},
	}
}

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
