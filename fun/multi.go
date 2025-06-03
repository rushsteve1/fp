package fun

// Multi takes multiple functions, passes the same argument to all of them,
// and then returns all the results.
// This is like currying multiple functions together.

func Multi2[T any, A, B any](v T, a func(T) A, b func(T) B) (A, B) {
	return a(v), b(v)
}

func Multi3[T any, A, B, C any](v T, a func(T) A, b func(T) B, c func(T) C) (A, B, C) {
	return a(v), b(v), c(v)
}

func Multi4[T any, A, B, C, D any](v T, a func(T) A, b func(T) B, c func(T) C, d func(T) D) (A, B, C, D) {
	return a(v), b(v), c(v), d(v)
}
