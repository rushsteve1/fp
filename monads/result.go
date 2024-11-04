package monads

type Result[T any] struct {
	V   T
	Err error
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{
		V:   v,
		Err: err,
	}
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
