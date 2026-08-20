package result

// Unpack returns the contained value and a nil error if res is [Ok]. If res is
// [Err], it returns the zero value of T and the contained error.
//
// The contained error may be nil. A nil error on the [Err] path returns the
// zero value of T and a nil error, which is indistinguishable from [Ok] of the
// zero value. Use [Result.Get] to distinguish those cases.
func Unpack[T any](res Result[T, error]) (T, error) {
	if res.IsErr() {
		var t T
		return t, res.err
	}

	return res.value, nil
}

// Get returns the contained value and true. If res is [Err], it returns the
// zero value of T and false.
func (res Result[T, E]) Get() (T, bool) {
	if res.IsErr() {
		var t T
		return t, false
	}

	return res.value, true
}

// GetErr returns the contained error and true. If res is [Ok], it returns the
// zero value of E and false.
func (res Result[T, E]) GetErr() (E, bool) {
	if res.IsOk() {
		var e E
		return e, false
	}

	return res.err, true
}
