package fun

import (
	"strconv"
	"testing"
)

func TestChain(t *testing.T) {
	f := Chain[int, int](
		func(f int) int { return f * 2 },
		strconv.Itoa,
		func(s string) int { return len(s) },
		func(x int) int { return x * 2 },
	)

	if f(2) != 2 {
		t.Fail()
	}
}

func TestChain2(t *testing.T) {
	f := Chain2(
		func(f int) int { return f * 2 },
		strconv.Itoa,
	)
	if f(2) != "4" {
		t.Fail()
	}
}

func TestChain3(t *testing.T) {
	f := Chain3(
		func(f int) int { return f * 2 },
		strconv.Itoa,
		func(s string) int { return len(s) },
	)
	if f(2) != 1 {
		t.Fail()
	}
}

func TestChain4(t *testing.T) {
	f := Chain4(
		func(f int) int { return f * 2 },
		strconv.Itoa,
		func(s string) int { return len(s) },
		func(x int) int { return x * 2 },
	)

	if f(2) != 2 {
		t.Fail()
	}
}