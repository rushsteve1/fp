package monads

import (
	"fmt"
	"sync"

	. "github.com/rushsteve1/fp"
)

// Observable is similar to promises in other languages, but can be considered
// and async version of [Result] within this library. It is thread-safe.
// You can consider it a continuous stream of updates to the value.
//
// It is a [monads.Monad], and therefore also a [fp.Seq]
type Observable[T comparable] struct {
	mu    *sync.Mutex
	inner T
	err   error
	seq   Seq[T]
	yield func(T) bool
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
		mu:    &sync.Mutex{},
		inner: v,
		err:   nil,
		yield: func(T) bool {
			Check(fmt.Errorf("Set called on Observable with no subscribers, this will have no effect"))
			return false
		},
	}

	o.seq = SeqFunc[T](func(yield func(T) bool) {
		o.yield = yield
		yield(v)
	})

	return o
}

func (o Observable[T]) Ok() bool {
	return o.err == nil
}

func (o Observable[T]) Get() (T, error) {
	return o.inner, o.err
}

func (o Observable[T]) Seq(yield func(T) bool) {
	o.seq.Seq(yield)
}

// Set changes the value of the [Observable], which causes a new element to
// be yielded on the sequence
func (o Observable[T]) Set(v T) {
	o.mu.Lock()
	defer o.mu.Unlock()
	// Only update when the
	if v != o.inner {
		o.inner = v
		o.yield(v)
	}
}

// Subscribe starts a new goroutine that will continually listen for changes
// to the sequence of this [Observable], calling f for each new value
func (o Observable[T]) Subscribe(f func(T)) {
	go func() {
		for v := range o.Seq {
			f(v)
		}
	}()
}
