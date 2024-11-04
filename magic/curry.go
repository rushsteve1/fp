package magic

import (
	"fmt"
	"reflect"
)

func Curry[In, Out any](f any, args ...any) func(In) Out {
	vf := reflect.ValueOf(f)
	tf := vf.Type()
	if tf.Kind() != reflect.Func {
		panic("Curry first argument must be a function")
	}
	if tf.NumOut() > 1 {
		panic("Curry function must return at most 1 value")
	}
	if (tf.NumIn()-1) > len(args) {
		panic(fmt.Sprintf("Curry too many arguments for function expected %d got %d", tf.NumIn()-1, len(args)))
	}

	fargs := make([]reflect.Value, 0, len(args))
	for a := range args {
		fargs = append(fargs, reflect.ValueOf(a))
	}

	tout := reflect.TypeFor[func(In) Out]()

	vout := reflect.MakeFunc(tout, func(outargs []reflect.Value) (results []reflect.Value) {
		args := make([]reflect.Value, 0, len(fargs) + 1)
		args = append(args, outargs[0])
		args = append(args, fargs...)
		return vf.Call(args)
	})

	cast, ok := vout.Interface().(func(In) Out)
	if !ok {
		panic("Curry could not cast to OneArg")
	}
	return cast
}

func Curry2[A, B, Out any](f func(A, B) Out, b B) func(A) Out {
	return func(a A) Out {
		return f(a, b)
	}
}

func Curry3[A, B, C, Out any](f func(A, B, C) Out, c C) func(A, B) Out {
	return func(a A, b B) Out {
		return f(a, b, c)
	}
}

func Curry32[A, B, C, Out any](f func(A, B, C) Out, b B, c C) func(A) Out {
	return func(a A) Out {
		return Curry2(Curry3(f, c), b)(a)
	}
}
