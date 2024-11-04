package transducers_test

import (
	"slices"
	"strconv"
	"testing"
	"time"

	. "iter"

	. "github.com/rushsteve1/fp/generators"
	. "github.com/rushsteve1/fp/reducers"
	. "github.com/rushsteve1/fp/threading"
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
		Integers(),
	)
	if s != "5" {
		t.Errorf("%s != \"6\"", s)
	}
}

func BenchmarkTransduce(b *testing.B) {
	Transduce(
		Chain3(
			Curry2(Take[int], 5),
			Curry2(Map, func(x int) int { return x + 1 }),
			Curry2(Map, strconv.Itoa),
		),
		Max,
		Integers(),
	)
}

func TestTransducerSeconds(t *testing.T) {
	Transduce(
		Curry2(Take[time.Time], 5),
		Collect,
		Ticker(time.Second),
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
	seq := slices.Values([]int{1, 2, 3, 4, 5})
	ar := Collect(Take(seq, 3))
	if !slices.Equal(ar, []int{1, 2, 3}) {
		t.Errorf("%v is not right", ar)
	}
}

func TestTakeInfinite(t *testing.T) {
	ar := Collect(Take(Integers(), 5))
	if !slices.Equal(ar, []int{0, 1, 2, 3, 4}) {
		t.Errorf("%v is not right", ar)
	}
}

func TestTakeTransducer(t *testing.T) {
	c := Curry2(Take[int], 5)
	c(Seq[int](Integers()))
}

func TestDebounce(t *testing.T) {
	seq := Ticker(100 * time.Millisecond)
	seq = Debounce(seq, 1*time.Second)
	seq = Take(seq, 5)
	ar := Collect(seq)
	t.Log(ar)
	if len(ar) != 5 {
		t.Fail()
	}
}

func TestDebounceTransducer(t *testing.T) {
	avg := Transduce(
		Chain3(
			Curry2(Debounce[time.Time], 1*time.Second),
			Curry2(Take[time.Time], 5),
			TimeDelta,
		),
		Average,
		Ticker(100*time.Millisecond),
	)
	t.Log(avg)
	if int(avg.Seconds()) != 1 {
		t.Fail()
	}
}
