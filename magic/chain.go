package magic

import (
	"fmt"
	"reflect"
)

func Chain[In, Out any](fs ...any) func(In) Out {
	vfs := make([]reflect.Value, 0, len(fs))
	for i, f := range fs {
		vf := reflect.ValueOf(f)
		tf := vf.Type()
		if tf.NumIn() != 1 {
			panic(fmt.Sprintf("Chain functions must take exactly 1 argument, arg %d is %s", i, tf.String()))
		}
		if tf.NumOut() != 1 {
			panic(fmt.Sprintf("Chain functions must return exactly 1 value, arg %d is %s", i, tf.String()))
		}
		vfs = append(vfs, vf)
	}

	vfirst := vfs[0]
	vlast := vfs[len(vfs)-1]

	if vfirst.Type().In(0) != reflect.TypeFor[In]() {
		panic("Chain first function argument does not match generic")
	}

	if vlast.Type().Out(0) != reflect.TypeFor[Out]() {
		panic("Chain last function return dos not match generic")
	}

	// Pointless but makes Chain() work
	if len(vfs) == 0 {
		var o Out
		return func(In) Out { return o }
	}

	tout := reflect.TypeFor[func(In) Out]()

	vout := reflect.MakeFunc(tout, func(args []reflect.Value) (results []reflect.Value) {
		arg := args[0]
		for _, vf := range vfs {
			arg = vf.Call([]reflect.Value{arg})[0]
		}
		return []reflect.Value{arg}
	})

	cast, ok := vout.Interface().(func(In) Out)
	if !ok {
		panic("Chain function could not cast to OneArg")
	}
	return cast
}

func Thread[T, U any](start T, fs ...any) U {
	return Chain[T, U](fs...)(start)
}

func Chain2[A, B, Out any](a func(A) B, b func(B) Out) func(A) Out {
	return func(x A) Out {
		return b(a(x))
	}
}

func Chain3[A, B, C, Out any](a func(A) B, b func(B) C, c func(C) Out) func(A) Out {
	return func(x A) Out {
		return c(b(a(x)))
	}
}

func Chain4[A, B, C, D, Out any](a func(A) B, b func(B) C, c func(C) D, d func(D) Out) func(A) Out {
	return func(x A) Out {
		return d(c(b(a(x))))
	}
}
