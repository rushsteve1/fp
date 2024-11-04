package magic

import (
	"strconv"
	"testing"
)

func TestChain(t *testing.T) {
	f := Chain[int, string](
		func(f int) int { return f * 2 },
		strconv.Itoa,
	)

	if f(2) != "4" {
		t.Fail()
	}
}
