package transducers_test

import (
	"slices"
	"strconv"
	"testing"
	"time"

	. "iter"

	. "github.com/rushsteve1/fp/generators"
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
		Integers,
	)
	if s != "6" {
		t.Errorf("%s != \"6\"", s)
	}
}

func TestTransducerSeconds(t *testing.T) {
	Transduce(
		Curry2(Take[time.Time], 5),
		Collect,
		Seq[time.Time](Seconds),
	)
}

func TestMap(t *testing.T) {
	seq := Seq[int](slices.Values([]int{1, 2, 3}))
	tx1 := Map(seq, func(x int) int {
		return x * 2
	})
	tx2 := Map(tx1, func(x int) string {
		return strconv.Itoa(x)
	})
	sl := slices.Collect(tx2)

	if !slices.Equal(sl, []string{"2", "4", "6"}) {
		t.Fail()
	}
}

func TestTake(t *testing.T) {
	Take(Seq[int](Integers), 10)
}

func TestTakeTransducer(t *testing.T) {
	c := Curry2(Take[int], 5)
	c(Seq[int](Integers))
}
