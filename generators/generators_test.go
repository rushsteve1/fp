package generators

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	i := 0
	for _ = range Ticker(time.Millisecond * 100) {
		if i > 5 {
			break
		}
		i++
	}
}
