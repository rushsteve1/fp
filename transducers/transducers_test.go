package transducers_test

import (
	"bytes"
	"io"
	"slices"
	"strconv"
	"testing"
	"time"

	. "github.com/rushsteve1/fp"
	. "github.com/rushsteve1/fp/fun"
	. "github.com/rushsteve1/fp/generators"
	. "github.com/rushsteve1/fp/reducers"
	. "github.com/rushsteve1/fp/transducers"
)

func TestTransduce(t *testing.T) {
	s := Transduce(
		Integers(),
		Chain3(
			Curry2(Take[int], 5),
			Curry2(Map, func(x int) int { return x + 1 }),
			Curry2(Map, strconv.Itoa),
		),
		Max,
	)
	AssertEq(t, s, "5")
}

func BenchmarkTransduce(b *testing.B) {
	Transduce(
		Integers(),
		Chain3(
			Curry2(Take[int], 5),
			Curry2(Map, func(x int) int { return x + 1 }),
			Curry2(Map, strconv.Itoa),
		),
		Max,
	)
}

func TestTransducerSeconds(t *testing.T) {
	Transduce(
		Ticker(time.Second),
		Curry2(Take[time.Time], 5),
		Collect,
	)
}

func TestMap(t *testing.T) {
	seq := SeqFunc[int](slices.Values([]int{1, 2, 3}))
	tx1 := Map(seq, func(x int) int {
		return x * 2
	})
	tx2 := Map(tx1, func(x int) string {
		return strconv.Itoa(x)
	})

	AssertSliceEq(t, Collect(tx2), []string{"2", "4", "6"})
}

func TestTake(t *testing.T) {
	seq := SeqFunc[int](slices.Values([]int{1, 2, 3, 4, 5}))
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
	AssertEq(t, len(ar), 5)
}

func TestDebounceTransducer(t *testing.T) {
	avg := Transduce(
		Ticker(100*time.Millisecond),
		Chain3(
			Curry2(Debounce[time.Time], 1*time.Second),
			Curry2(Take[time.Time], 5),
			TimeDelta,
		),
		Average,
	)
	AssertEq(t, int(avg.Seconds()), 1)
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	ar := Transduce(
		Integers(),
		Chain4(
			Curry2(Take[int], 5),
			Curry2(Map, func(i int) []byte {
				return []byte(strconv.Itoa(i))
			}),
			Visitor(Curry2(Write, io.Writer(&buf))),
			Curry2(Map, func(b []byte) string {
				return string(b)
			}),
		),
		Collect,
	)
	t.Log(buf)
	t.Log(ar)
	AssertEq(t, buf.Len(), len(ar))
}
