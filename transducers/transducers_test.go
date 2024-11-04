package transducers

import (
	"slices"
	"strconv"
	"testing"

	"github.com/rushsteve1/fp/iter"
)

func TestMap(t *testing.T) {
	seq := iter.SeqFunc[int](slices.Values([]int{1, 2, 3}))
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
