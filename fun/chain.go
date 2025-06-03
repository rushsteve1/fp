package fun

import (
	"fmt"
	"reflect"
)

// Chain implementents function chaining or threading, where each function is
// passed the return value of the previous.
// This is particularly useful in Transducers that all operate over [fp.Seq]

// Reflective version
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

// Generics up to 4, could extend with go:generate
func Chain2[In, B, Out any](a func(In) B, b func(B) Out) func(In) Out {
	return func(x In) Out {
		return b(a(x))
	}
}

func Chain3[In, B, C, Out any](a func(In) B, b func(B) C, c func(C) Out) func(In) Out {
	return func(x In) Out {
		return c(b(a(x)))
	}
}

func Chain4[In, B, C, D, Out any](a func(In) B, b func(B) C, c func(C) D, d func(D) Out) func(In) Out {
	return func(x In) Out {
		return d(c(b(a(x))))
	}
}
