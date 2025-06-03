package fun

import (
	"fmt"
	"reflect"
)

// The set of currying functions takes a series of functions and returns
// a new function that takes only the first argument of the first function.
//
// Because of a limitation of Go and the design of this library, all curried
// functions take only a single argument.

func Curry2[In, B, Out any](f func(In, B) Out, b B) func(In) Out {
	return func(a In) Out {
		return f(a, b)
	}
}

func Curry3[In, B, C, Out any](f func(In, B, C) Out, b B, c C) func(In) Out {
	return func(a In) Out {
		return f(a, b, c)
	}
}

func Curry4[In, B, C, D, Out any](f func(In, B, C, D) Out, b B, c C, d D) func(In) Out {
	return func(a In) Out {
		return f(a, b, c, d)
	}
}

func Curry[In, Out any](f any, args ...any) func(In) Out {
	vf := reflect.ValueOf(f)
	c := vf.Type().NumIn() - 1
	if c != len(args) {
		panic(fmt.Sprintf("Curry incorect args for function expected %d got %d", c, len(args)))
	}
	vargs := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		vargs = append(vargs, reflect.ValueOf(arg))
	}
	return curry[In, Out](vf, vargs)
}

func curry[In, Out any](vf reflect.Value, vargs []reflect.Value) func(In) Out {
	tf := vf.Type()
	if tf.Kind() != reflect.Func {
		panic("Curry first argument must be a function")
	}
	if tf.NumOut() > 1 {
		panic("Curry function must return at most 1 value")
	}

	for i, varg := range vargs {
		targ := varg.Type()
		tfi := tf.In(i + 1)
		if tfi != targ {
			panic(fmt.Sprintf("Curry argument %d is %s expected %s", i, targ, tfi))
		}
	}

	tout := reflect.TypeFor[func(In) Out]()
	vout := reflect.MakeFunc(tout, func(vs []reflect.Value) []reflect.Value {
		args := make([]reflect.Value, 0, len(vargs)+1)
		args = append(args, vs[0])
		args = append(args, vargs...)
		return vf.Call(args)
	})

	cast, ok := vout.Interface().(func(In) Out)
	if !ok {
		panic("Curry could not cast to OneArg")
	}
	return cast
}
