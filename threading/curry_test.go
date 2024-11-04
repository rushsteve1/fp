package threading

import "testing"

func TestCurry(t *testing.T) {
	foo := func(a, b, c int) int {
		return a + b + c
	}
	f := Curry[int, int](foo, 2, 3)
	v := f(1)
	if v != 6 {
		t.Errorf("%v != 6", v)
	}
}

func TestCurry2(t *testing.T) {
	foo := func(a, b int) int {
		return a + b
	}
	f := Curry2(foo, 2)
	v := f(1)
	if v != 3 {
		t.Errorf("%v != 3", v)
	}
}

func TestCurry3(t *testing.T) {
	foo := func(a, b, c int) int {
		return a + b + c
	}
	f := Curry3(foo, 2, 3)
	v := f(1)
	if v != 6 {
		t.Errorf("%v != 6", v)
	}
}
