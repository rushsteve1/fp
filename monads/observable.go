package monads

import (
	. "github.com/rushsteve1/fp"
)

// ObserableBufSize is the size of the channel buffer used by [Observable]
var ObserableBufSize = 5

// Observable is similar to promises in other languages, but can be considered
// and async version of [Result] within this library. It is thread-safe.
// You can consider it a continuous stream of updates to the value.
//
// It is a [monads.Monad], and therefore also a [fp.Seq].
// Because sequences are lazy an Observable with no subscribers does nothing.
type Observable[T comparable] struct {
	RWLock[T]
	err error
	c   chan T
	seq Seq[T]
}

// Promise calls f on a new goroutine and immediately returns an [Observable]
// that will be updated when the provided function returns
func Promise[T comparable](f func() T) Observable[T] {
	var v T
	o := Observe(v)
	go func() { o.Set(f()) }()
	return o
}

// Observe creates an [Observable] from a single value
func Observe[T comparable](v T) Observable[T] {
	o := Observable[T]{
		RWLock: NewRWLock(v),
		err:    nil,
		c:      make(chan T, ObserableBufSize),
	}

	o.seq = SeqFunc[T](func(yield func(T) bool) {
		v, _ := o.Get()
		yield(v)
		for t := range o.c {
			o.Set(v)
			if !yield(t) {
				return
			}
		}
	})

	return o
}

func (o Observable[T]) Ok() bool {
	return o.err == nil
}

func (o Observable[T]) Seq(yield func(T) bool) {
	o.seq.Seq(yield)
}

// Set changes the value of the [Observable], which causes a new element to
// be yielded on the sequence.
//
// WARNING: The internal channel that this function uses has a buffer set by
// [ObserableBufSize] which defaults to 5.
// Calls to Set without any subscribers can overflow this buffer.
func (o Observable[T]) Set(v T) {
	// Only update when the
	cur, _ := o.Get()
	if v != cur {
		o.c <- v
	}
}

// Subscription is a handle that can be used to stop after calling [Subscribe].
// This type simply implements [io.Closer]
type Subscription struct {
	stop chan bool
}

func (s Subscription) Close() {
	s.stop <- true
}

// Subscribe starts a new goroutine that will continually listen for changes
// to the sequence of this [Observable], calling f for each new value.
//
// It returns a [Subscription] handle that can be used to close it
func (o Observable[T]) Subscribe(f func(T)) {
	sub := Subscription{
		stop: make(chan bool),
	}

	go func() {
		for v := range o.Seq {
			select {
			case <-sub.stop:
				return
			default:
				f(v)
			}
		}
	}()
}

func (o Observable[T]) Close() error {
	close(o.c)
	return nil
}
