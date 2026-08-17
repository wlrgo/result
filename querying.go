package result

// IsOk returns true if the result is [Ok].
func (res Result[T, E]) IsOk() bool {
	return res.ok
}

// IsErr returns true if the result is [Err].
func (res Result[T, E]) IsErr() bool {
	return !res.ok
}

// IsOkAnd returns true if the result is [Ok] and the value inside of it
// matches a predicate.
func (res Result[T, E]) IsOkAnd(f func(T) bool) bool {
	return res.ok && f(res.value)
}

// IsErrAnd returns true if the result is [Err] and the value inside of it
// matches a predicate.
func (res Result[T, E]) IsErrAnd(f func(E) bool) bool {
	return !res.ok && f(res.err)
}
