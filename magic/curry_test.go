package magic

import "testing"

func TestCurry(t *testing.T) {
	foo := func(a, b, c int) int {
		return a + b + c
	}
	f := Curry[int, int](foo, 1, 2)
	v := f(2)
	if v != (1 + 2 + 3) {
		t.Fail()
	}
}
