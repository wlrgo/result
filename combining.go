package result

// And
func And[T, U, E any](res Result[T, E], other Result[U, E]) Result[U, E] {
	return AndThen(res, func(T) Result[U, E] { return other })
}

// AndThen
func AndThen[T, U, E any](res Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	if res.IsOk() {
		return f(res.value)
	}

	return Err[U](res.err)
}

// Or
func Or[T, E, F any](res Result[T, E], other Result[T, F]) Result[T, F] {
	return OrElse(res, func(E) Result[T, F] { return other })
}

// OrElse
func OrElse[T, E, F any](res Result[T, E], f func(E) Result[T, F]) Result[T, F] {
	if res.IsErr() {
		return f(res.err)
	}

	return Ok[T, F](res.value)
}
