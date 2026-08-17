package result

// Expect
func (res Result[T, E]) Expect(msg string) T {
	if res.IsErr() {
		// TODO: print res.err
		panic(msg)
	}

	return res.value
}

// Unwrap
func (res Result[T, E]) Unwrap() T {
	// TODO: print res.err
	return res.Expect("called `Result::unwrap()` on an `Err` value")
}

// UnwrapOr
func (res Result[T, E]) UnwrapOr(value T) T {
	return res.UnwrapOrElse(func(E) T { return value })
}

// UnwrapOrDefault
func (res Result[T, E]) UnwrapOrDefault() T {
	return res.UnwrapOrElse(func(E) T {
		var t T
		return t
	})
}

// UnwrapOrElse
func (res Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if res.IsErr() {
		return f(res.err)
	}

	return res.value
}

// ExpectErr
func (res Result[T, E]) ExpectErr(msg string) E {
	if res.IsOk() {
		// TODO: print res.value
		panic(msg)
	}

	return res.err
}

// UnwrapErr
func (res Result[T, E]) UnwrapErr() E {
	// TODO: print res.value
	return res.ExpectErr("called `Result::unwrap_err()` on an `Ok` value")
}
