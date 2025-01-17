package coll

import (
	. "github.com/rushsteve1/fp"
)

// RetainedSeq is a type of sequence that remembers what elements have already
// occured. It exposes two methods, Back and Reset for working with this history.
type RetainedSeq[T any] struct {
	inner Seq[T]
	prev  []T
}

func (rs RetainedSeq[T]) Seq(yield func(T) bool) {
	rs.inner.Seq(yield)
}

func NewRetainedSeq[T any](seq Seq[T]) (out RetainedSeq[T]) {
	out.prev = make([]T, 0, 8)

	seq = SeqFunc[T](func(yield func(T) bool) {
		seq.Seq(func(t T) bool {
			out.prev = append(out.prev, t)
			return yield(t)
		})
	})
	out.inner = seq

	return out
}

func (rs RetainedSeq[T]) Back() (out RetainedSeq[T]) {
	rs.rollback(1)
}

func (rs RetainedSeq[T]) Reset() (out RetainedSeq[T]) {
	return rs.rollback(len(rs.prev))
}

func (rs RetainedSeq[T]) rollback(i int) (out RetainedSeq[T]) {
	if len(rs.prev) == 0 {
		return rs
	}

	out.prev = make([]T, 0, 8)
	out.inner = SeqFunc[T](func(yield func(T) bool) {
		didprev := false
		rs.Seq(func(t T) bool {
			if !didprev {
				didprev = true
				for _, v := range rs.prev {
					if !yield(v) {
						return false
					}
				}
			}
			return yield(t)
		})
	})

	return out
}
