package monads

// Result is a monad that can indicate failure.
// It is the same as in Rust
type Result[T any] struct {
	V   T
	Err error
}

func (o Result[T]) Ok() bool {
	return o.Err != nil
}

func (r Result[T]) Get() (T, error) {
	return r.V, r.Err
}

func (r Result[T]) Seq(yield func(T) bool) {
	if r.Err == nil {
		yield(r.V)
	}
}

// Wrap takes a value and an error and produces a [Result].
// Two-value spreading can be used to easily compose this with many functions.
func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{
		V:   v,
		Err: err,
	}
}

// FuncWrap is like [Wrap] but you provide the function directly.
// This goes very well with the Currying tools.
func FuncWrap[In, Out any](f func(In) (Out, error)) func(In) Result[Out] {
	return func(a In) Result[Out] {
		return Wrap(f(a))
	}
}
