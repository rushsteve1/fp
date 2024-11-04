package generators_test

import (
	"testing"
	"time"
	
	. "github.com/rushsteve1/fp/generators"
)

func TestTicker(t *testing.T) {
	i := 0
	for _ = range Ticker(time.Millisecond * 100).Seq {
		if i > 5 {
			break
		}
		i++
	}
}

func TestChan(t *testing.T) {
	c := make(chan int)
	s := Chan(c)

	go func() {
		for i := range 5 {
			c <- i
		}
		close(c)
	}()

	i := 0
	for v := range s.Seq {
		t.Log(v)
		i++
	}

	if i != 5 {
		t.Fail()
	}
}
