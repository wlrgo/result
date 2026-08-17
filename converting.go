package result

// Unpack
func Unpack[T any](res Result[T, error]) (T, error) {
	if res.IsErr() {
		var t T
		return t, res.err // NOTE: res.err can be nil, so should we accept this?
	}

	return res.value, nil
}

// Get
func (res Result[T, E]) Get() (T, bool) {
	if res.IsErr() {
		var t T
		return t, false
	}

	return res.value, true
}

// GetErr
func (res Result[T, E]) GetErr() (E, bool) {
	if res.IsOk() {
		var e E
		return e, false
	}

	return res.err, true
}
