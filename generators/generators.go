package generators

import (
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"time"

	. "github.com/rushsteve1/fp"
	"github.com/rushsteve1/fp/monads"
)

// A generator is any function that returns a new sequence

// Empty returns an infinite empty sequence that never yields any values
func Empty[T any]() Seq[T] {
	return SeqFunc[T](func(func(T) bool) {})
}

// Generate returns an iterator that repeatedly calls the provided function,
// yielding the values it returns
func Generate[T any](f func() T) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		for yield(f()) {
		}
	})
}

// Once returns a sequence that yields the provided value once
func Once[T any](v T) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		yield(v)
	})
}

// Forever returns an infinite sequence of the provided value
func Forever[T any](v T) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		for yield(v) {
		}
	})
}

// Integers yields an infinite sequence of integers
func Integers() Seq[int] {
	return SeqFunc[int](func(yield func(int) bool) {
		n := 0
		for yield(n) {
			n++
		}
	})
}

// Shuffle returns an infinite sequence of random elements from pick
func Shuffle[E ~[]T, T any](pick E) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		for {
			n := rand.Int() % len(pick)
			if !yield(pick[n]) {
				return
			}
		}
	})
}

// Ticker yields the current time when the duration has passed
func Ticker(d time.Duration) Seq[time.Time] {
	return SeqFunc[time.Time](func(yield func(time.Time) bool) {
		for t := range time.Tick(d) {
			if !yield(t) {
				return
			}
		}
	})
}

// Chan returns an iterator that continually yields values from the channel.
// The channel is closed if the sequence stops.
// It is the caller's responsibility to close the channel to prevent deadlocks
func Chan[T any](c chan T) Seq[T] {
	return SeqFunc[T](func(yield func(T) bool) {
		for t := range c {
			if !yield(t) {
				close(c)
				return
			}
		}
	})
}

// Reader reads from the passed [io.Reader] turning into a sequence of byte arrays.
// See its counterpart [transducers.Writer]
func Reader(r io.Reader) Seq[monads.Result[[]byte]] {
	return SeqFunc[monads.Result[[]byte]](func(yield func(monads.Result[[]byte]) bool) {
		var buf []byte
		for {
			_, err := r.Read(buf)
			if errors.Is(err, io.EOF) {
				return
			}
			if !yield(monads.Wrap(buf, err)) {
				return
			}
		}
	})
}

// Accept takes a listener and returns a sequence of accepted connections.
// The listener is closed if the sequence stops.
func Accept(l net.Listener) Seq[net.Conn] {
	return SeqFunc[net.Conn](func(yield func(net.Conn) bool) {
		for {
			c := Must(l.Accept())
			if !yield(c) {
				l.Close()
				return
			}
		}
	})
}
