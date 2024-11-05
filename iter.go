// This package wraps the standard library's [iter] package, providing some
// additional features.
//
// It is intended to potentially inform future development and act as the
// backbone of this library.

package fp

import (
	"iter"
)

// SeqFunc is exactly the same as [iter.Seq] and can be trivially cast between
type SeqFunc[V any] iter.Seq[V]

// Seq borrows a trick used by [http.Handler] to define an interface and a func
// that implements that interface by calling itself [SeqFunc]
type Seq[V any] interface {
	// Seq implements push-style iteration using the yield callback.
	// See the documenation of [iter] for more information.
	// It has exactly the same signature as [SeqFunc].
	Seq(yield func(V) bool)
}

func (sf SeqFunc[V]) Seq(yield func(V) bool) {
	sf(yield)
}

// KeyValue wholes a key-value pair
type KeyValue[K, V any] struct {
	Key   K
	Value V
}

// Seq2Func is exactly the same as [iter.Seq2] and can be trivially cast between.
type Seq2Func[K, V any] iter.Seq2[K, V]

// Seq2 is to [Seq] what [iter.Seq] is to [iter.Seq2]
// but with the addition requirement that Seq2 values can be used as
type Seq2[K, V any] interface {
	// I don't like the whole [iter.Seq2] thing that the stdlib does
	// so we use this to convert [Seq2] into [Seq]
	// This trivially gives compatibility with the rest of this library
	Seq(yield func(KeyValue[K, V]) bool)
	Seq2(yield func(K, V) bool)
}

func (sf Seq2Func[K, V]) Seq(yield func(KeyValue[K, V]) bool) {
	sf(func(k K, v V) bool {
		return yield(KeyValue[K, V]{k, v})
	})
}

func (sf Seq2Func[K, V]) Seq2(yield func(K, V) bool) {
	sf(yield)
}

// Pull is a wrappr around [iter.Pull]
func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	return iter.Pull(seq.Seq)
}

// Pull2 is a wrapper around [iter.Pull2]
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	return iter.Pull2(seq.Seq2)
}

// Duet is the inverse of [Seq2.Seq] taking a [Seq] of [KeyValue]
// and returning the [Seq2] equivalent
func Duet[K, V any](seq Seq[KeyValue[K, V]]) Seq2[K, V] {
	return Seq2Func[K, V](func(yield func(K, V) bool) {
		seq.Seq(func(kv KeyValue[K, V]) bool {
			return yield(kv.Key, kv.Value)
		})
	})
}
