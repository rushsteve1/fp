package generators

import (
	"io"
	. "iter"
	"time"

	"github.com/rushsteve1/fp"
)

// A generator is any function that returns a new sequence

// Empty returns an infinite empty sequence that never yields any values
func Empty[T any]() Seq[T] {
	return func(func(T) bool){}
}

// Generate returns an iterator that repeatedly calls the provided function,
// yielding the values it returns
func Generate[T any](f func() T) Seq[T] {
	return func(yield func(T) bool) {
		if !yield(f()) {
			return
		}
	}
}

// Forever returns an infinite sequence of the provided value
func Forever[T any](v T) Seq[T] {
	return func(yield func(T) bool){
		if !yield(v) {
			return
		}
	}
}

// Integers yields an infinite sequence of integers
func Integers() Seq[int] {
	return func(yield func(int) bool) {
		n := 0
		for yield(n) {
			n++
		}
	}
}

// Ticker yields the current time when the duration has passed
func Ticker(d time.Duration) Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		for t := range time.Tick(d) {
			if !yield(t) {
				return
			}
		}
	}
}

// Chan returns an iterator that continually yields values from the channel.
// It is the caller's responsibility to close the channel and prevent deadlocks
func Chan[T any](c chan T) Seq[T] {
	return func(yield func(T) bool) {
		for t := range c {
			if !yield(t) {
				return
			}
		}
	}
}

func Reader(r io.Reader) Seq[[]byte] {
	return func(yield func([]byte) bool) {
		var buf []byte
		fp.Must(r.Read(buf))
		if !yield(buf) {
			return
		}
	}
}