package result

import "cmp"

// Compare returns -1 if x is less than y, 0 if x equals y, or +1 if x is
// greater than y.
//
// [Ok] is less than any [Err] value. Two [Err] values are compared by their
// errors. If both are [Ok], they are compared with [cmp.Compare].
//
// For floating-point values, [cmp.Compare] treats NaN as equal to NaN. [Equal]
// uses Go ==, so two [Ok] NaN values compare equal here and unequal there.
func Compare[T, E cmp.Ordered](x, y Result[T, E]) int {
	switch {
	case x.IsErr() && y.IsErr():
		return cmp.Compare(x.err, y.err)
	case x.IsErr():
		return 1
	case y.IsErr():
		return -1
	default:
		return cmp.Compare(x.val, y.val)
	}
}

// Equal reports whether x and y are equal.
//
// Two [Err] values are equal if their errors are equal. An [Err] is not equal
// to any [Ok]. Two [Ok] values are equal if their contained values are equal
// with Go ==.
//
// For floating-point values, two [Ok] NaN values are not equal. [Compare]
// uses [cmp.Compare], so those same values compare equal there.
func Equal[T, E comparable](x, y Result[T, E]) bool {
	if x.IsErr() && y.IsErr() {
		return x.err == y.err
	}

	if x.IsOk() && y.IsOk() {
		return x.val == y.val
	}

	return false
}

// Ge reports whether x is greater than or equal to y.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Ge[T, E cmp.Ordered](x, y Result[T, E]) bool {
	return Compare(x, y) >= 0
}

// Gt reports whether x is greater than y.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Gt[T, E cmp.Ordered](x, y Result[T, E]) bool {
	return Compare(x, y) > 0
}

// Le reports whether x is less than or equal to y.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Le[T, E cmp.Ordered](x, y Result[T, E]) bool {
	return Compare(x, y) <= 0
}

// Lt reports whether x is less than y.
//
// See [Compare] for the ordering of [Ok] and [Err].
func Lt[T, E cmp.Ordered](x, y Result[T, E]) bool {
	return Compare(x, y) < 0
}
