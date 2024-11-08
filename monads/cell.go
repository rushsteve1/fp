package monads

// Gettable is a simple interface for get-ing an inner value
type Gettable[T any] interface {
	// Get does not return a monad in order to be more compatible
	// with other libraries.
	// And to prevent an odd self-referencing type
	Get() (T, error)
}

// Mutable is a simple interface for set-ing an inner value
type Mutable[T any] interface {
	Gettable[T]
	Set(T)
}

// Cell is a type that implements the idea of "interior mutability" from Rust.
// The Cell itself is an immutable data structure that acts as a smart pointer
// to a value that *can* be updated.
type Cell[T any] struct {
	v *T
}

// NewCell creates a new [Cell]
func NewCell[T any](v T) Cell[T] {
	return Cell[T]{v: &v}
}

func (c Cell[T]) Get() (out T, err error) {
	// The only way to construct or update a Cell is with owned values
	// so this is always safe
	return *c.v, nil
}

func (c Cell[T]) Set(v T) {
	*c.v = v
}
