package result

import "cmp"

// Compare returns -1 if a is less than b, 0 if a equals b, or +1 if a is
// greater than b.
//
// [Ok] is less than any [Err] value. Two [Err] values are compared by their
// errors. If both are [Ok], they are compared with [cmp.Compare].
//
// For floating-point values, [cmp.Compare] treats NaN as equal to NaN. [Equal]
// uses Go ==, so two [Ok] NaN values compare equal here and unequal there.
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

// Equal reports whether a and b are equal.
//
// Two [Err] values are equal if their errors are equal. An [Err] is not equal
// to any [Ok]. Two [Ok] values are equal if their contained values are equal
// with Go ==.
//
// For floating-point values, two [Ok] NaN values are not equal. [Compare]
// uses [cmp.Compare], so those same values compare equal there.
func Equal[T, E comparable](a, b Result[T, E]) bool {
	if a.IsErr() && b.IsErr() {
		return a.err == b.err
	}

	if a.IsOk() && b.IsOk() {
		return a.value == b.value
	}

	return false
}

// Ge reports whether a is greater than or equal to b.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Ge[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) >= 0
}

// Gt reports whether a is greater than b.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Gt[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) > 0
}

// Le reports whether a is less than or equal to b.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Le[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) <= 0
}

// Lt reports whether a is less than b.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Lt[T, E cmp.Ordered](a, b Result[T, E]) bool {
	return Compare(a, b) < 0
}
