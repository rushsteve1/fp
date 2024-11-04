package monads

import "testing"

func TestMonadInterface(t *testing.T) {
	v := Some(1)
	// If this cast succedes the test does
	m := Monad[int](v)
	t.Log(m)
}
