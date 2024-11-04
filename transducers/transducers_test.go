package transducers_test

import (
	"slices"
	"strconv"
	"testing"

	. "github.com/rushsteve1/fp/iter"
	. "github.com/rushsteve1/fp/magic"
	. "github.com/rushsteve1/fp/reducers"
	. "github.com/rushsteve1/fp/transducers"
)

func TestTransduce(t *testing.T) {
	s := Transduce(
		Chain3(
			Curry2(Take[int], 5),
			Curry2(Map, func(x int) int { return x + 1 }),
			Curry2(Map, strconv.Itoa),
		),
		Max,
		SeqFunc[int](Integers),
	)
	if s != "6" {
		t.Errorf("%s != \"6\"", s)
	}
}

func TestMap(t *testing.T) {
	seq := SeqFunc[int](slices.Values([]int{1, 2, 3}))
	tx1 := Map(seq, func(x int) int {
		return x * 2
	})
	tx2 := Map(tx1, func(x int) string {
		return strconv.Itoa(x)
	})
	sl := slices.Collect(tx2.Seq)

	if !slices.Equal(sl, []string{"2", "4", "6"}) {
		t.Fail()
	}
}

func TestTake(t *testing.T) {
	Take(SeqFunc[int](Integers), 10)
}