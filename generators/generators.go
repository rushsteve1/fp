package generators

import (
	"iter"
	"time"
)

// Integers yields an infinite sequence of integers
func Integers() iter.Seq[int] {
	return func(yield func(int) bool) {
		n := 0
		for yield(n) {
			n++
		}
	}
}

// Ticker yields the current time every time the duration has passed
func Ticker(d time.Duration) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		for t := range time.Tick(d) {
			if !yield(t) {
				return
			}
		}
	}
}
