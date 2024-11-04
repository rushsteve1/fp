package threading

import (
	"fmt"
	"reflect"
)

func Curry2[A, B, Out any](f func(A, B) Out, b B) func(A) Out {
	return func(a A) Out {
		return f(a, b)
	}
}

func Curry3[A, B, C, Out any](f func(A, B, C) Out, b B, c C) func(A) Out {
	return func(a A) Out {
		return Curry2(Curry32(f, c), b)(a)
	}
}

func Curry32[A, B, C, Out any](f func(A, B, C) Out, c C) func(A, B) Out {
	return func(a A, b B) Out {
		return f(a, b, c)
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
