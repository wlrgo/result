package result

import "cmp"

// Compare
func Compare[T, E cmp.Ordered](a, b Result[T, E]) int {
	switch {
	case a.IsErr() && b.IsErr():
		return cmp.Compare(a.err, b.err)
	case a.IsErr():
		return 1
	case b.IsErr():
		return -1
	default:
		return cmp.Compare(a.value, b.value)
	}
}

// Equal
func Equal[T, E comparable](a, b Result[T, E]) bool {
	if a.IsErr() && b.IsErr() {
		return a.err == b.err
	}

	if a.IsOk() && b.IsOk() {
		return a.value == b.value
	}

	return false
}

// Ge
func Ge[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) >= 0
}

// Gt
func Gt[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) > 0
}

// Le
func Le[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) <= 0
}

// Lt
func Lt[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) < 0
}
