package result

// IsOk reports whether the result is [Ok].
func (res Result[T, E]) IsOk() bool {
	return res.ok
}

// IsErr reports whether the result is [Err].
func (res Result[T, E]) IsErr() bool {
	return !res.ok
}

// IsOkAnd reports whether the result is [Ok] and the contained value matches
// a predicate.
func (res Result[T, E]) IsOkAnd(f func(T) bool) bool {
	return res.ok && f(res.value)
}

// IsErrAnd reports whether the result is [Err] and the contained error matches
// a predicate.
func (res Result[T, E]) IsErrAnd(f func(E) bool) bool {
	return !res.ok && f(res.err)
}
