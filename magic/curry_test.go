package magic

import "testing"

func TestCurry(t *testing.T) {
	foo := func(a, b, c int) int {
		return a + b + c
	}
	f := Curry[int, int](foo, 1, 2)
	v := f(3)
	if v != 6 {
		t.Errorf("%v != 6", v)
	}
}
