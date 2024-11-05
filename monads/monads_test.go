package monads_test

import (
	"testing"

	"github.com/rushsteve1/fp"
	. "github.com/rushsteve1/fp/fun"
	. "github.com/rushsteve1/fp/monads"
	. "github.com/rushsteve1/fp/reducers"
	. "github.com/rushsteve1/fp/transducers"
)

func TestOptionMonadInterface(t *testing.T) {
	v := Some(1)
	// If this cast succedes the test does
	m := Monad[int](v)
	t.Log(m)
}

func TestResultMonadInterface(t *testing.T) {
	v := Wrap(1, nil)
	// If this cast succedes the test does
	m := Monad[int](v)
	t.Log(m)
}

func TestMonadTransducer(t *testing.T) {
	a := Transduce(
		Some(10),
		Curry2(Map, func(x int) int { return x * 2 }),
		Chain2(First[int], Some),
	)

	fp.AssertEq(t, a, Some(20))
}
