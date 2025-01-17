package coll

import (
	. "github.com/rushsteve1/fp"
	. "github.com/rushsteve1/fp/monads"
	"github.com/rushsteve1/fp/reducers"
)

// https://grishaev.me/en/clojure-zippers/

// type Node[T any] interface {
// 	*T | ~[]T | *Seq[T]
// }

type Zipper[T any] struct {
	cur    any
	branch RetainedSeq[T]
	parent *Zipper[T]
}

func branch[T any](x any) bool {
	switch x.(type) {
	case *Seq[T]:
		return true
	case []T:
		return true
	case *T:
		return false
	default:
		return false
	}
}

func NewZipper[T any](root Seq[T]) (out Zipper[T]) {
	out.branch = NewRetainedSeq(root)
	out.cur = reducers.First(out.branch)
	return out
}

func Up[T any](z Zipper[T]) Option[Zipper[T]] {
	return TrySome(z.parent)
}

func Down[T any](z Zipper[T]) Option[Zipper[T]] {
	if !branch[T](z.cur) {
		return Some(z)
	}
	var out Zipper[T]
	out.parent = &z
}
func Left[T any](z Zipper[T]) Zipper[T]  {}
func Right[T any](z Zipper[T]) Zipper[T] {}
func Next[T any](z Zipper[T]) Zipper[T]  {}

func (z Zipper[T]) Seq(yield func(T) bool) {

}
