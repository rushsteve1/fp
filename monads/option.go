package monads

import (
	"database/sql"
	"errors"

	"github.com/rushsteve1/fp"
)

var ErrUnwrapInvalid = errors.New("Unwrapped invalid Option")

type Option[T any] struct {
	sql.Null[T]
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		sql.Null[T]{
			V:     v,
			Valid: true,
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