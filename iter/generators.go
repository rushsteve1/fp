package iter

// Integers yields an infinite sequence of integers
func Integers(yield func(int) bool) {
	n := 0
	for yield(n) {
		n++
	}
}
