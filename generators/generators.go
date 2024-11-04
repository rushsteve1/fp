package generators

import (
	"time"
)

// Integers yields an infinite sequence of integers
func Integers(yield func(int) bool) {
	n := 0
	for yield(n) {
		n++
	}
}

// Seconds yields the current time every second
func Seconds(yield func(time.Time) bool) {
	c := time.Tick(time.Second)
	for t := range c {
		if !yield(t) {
			return
		}
	}
}
