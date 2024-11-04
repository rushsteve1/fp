package magic

import "golang.org/x/exp/constraints"

type Complex = constraints.Complex
type Integer = constraints.Integer
type Float = constraints.Float
type Ordered = constraints.Ordered
type Signed = constraints.Signed
type Unsigned = constraints.Unsigned

type Nilable interface {
	~*any | ~[]any | ~map[any]any
}

type Numeric interface {
	Integer | Float
}
