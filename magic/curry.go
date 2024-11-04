package magic

import (
	"fmt"
	"reflect"
)

type OneArg[A, Out any] = func(A) Out
type TwoArg[A1, A2, Out any] = func(A1, A2) Out

func Curry[In, Out any](f any, args ...any) OneArg[In, Out] {
	vf := reflect.ValueOf(f)
	tf := vf.Type()
	if tf.Kind() != reflect.Func {
		panic("Curry first argument must be a function")
	}
	if tf.NumOut() > 1 {
		panic("Curry function must return at most 1 value")
	}
	if tf.NumIn() > len(args) {
		panic(fmt.Sprintf("Curry too many arguments for function expected %d got %d", tf.NumIn(), len(args)))
	}

	fargs := make([]reflect.Value, 0, len(args))
	for a := range args {
		fargs = append(fargs, reflect.ValueOf(a))
	}

	tout := reflect.TypeFor[OneArg[In, Out]]()

	vout := reflect.MakeFunc(tout, func(outargs []reflect.Value) (results []reflect.Value) {
		arg := outargs[0]
		return vf.Call([]reflect.Value{arg})
	})

	cast, ok := vout.Interface().(OneArg[In, Out])
	if !ok {
		panic("Curry could not cast to OneArg")
	}
	return cast
}

func Curry2[A, B, Out any](f TwoArg[A, B, Out], b B) OneArg[A, Out] {
	return func(a A) Out {
		return f(a, b)
	}
}

func Curry3[A, B, C, Out any](f func(A, B, C) Out, c C) TwoArg[A, B, Out] {
	return func(a A, b B) Out {
		return f(a, b, c)
	}
}

func Curry32[A, B, C, Out any](f func(A, B, C) Out, b B, c C) func(A) Out {
	return func(a A) Out {
		return Curry2(Curry3(f, c), b)(a)
	}
}
