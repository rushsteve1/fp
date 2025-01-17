# Go Functional Programming

This library is cursed and you probably shouldn't use it.
It relies on the Dark Arts to enable a very non-idoimatic way of writing Go.

I was inspired by @fossker's wonderful talk at Doomconf 2022H2 and his
[Common Lisp Transducers](https://codeberg.org/fosskers/cl-transducers)
package and documentation.

But this is a different approach that aims to build upon Go's existing standard
library to introduce transducers in a simple way alongside a lot of other
FP features that I've wanted across various projects.

There are many FP/utility libraries for Go. This one is mine.

The end result is the ability to write Go code that looks like this
```go
Transduce(
	Integers(),
	Chain4(
		Curry2(Take[int], 5),
		Delta,
		Curry2(Map, func(i int) []byte {
			return []byte(strconv.Itoa(i))
		}),
		Visitor(Curry2(Write, io.Writer(&buf))),
	),
	Collect,
)
```

None of the ideas in this library are new.
I'm mostly copying features from the wonderful [Clojure](https://clojure.org)
and of course [Rust](https://rust-lang.org).

There are also some other kinda misc things in this library.
As its grown it's become more of a general purpose utility library
that I use. You can easily import just the parts you want.