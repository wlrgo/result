package result

// And returns other if res is [Ok]. Otherwise, it returns [Err] of res.
//
// other is evaluated before And is called. Use [AndThen] when the second
// [Result] should be computed only if res is [Ok].
func And[T, U, E any](res Result[T, E], other Result[U, E]) Result[U, E] {
	return AndThen(res, func(T) Result[U, E] { return other })
}

// AndThen returns the result of calling f with the contained value if res is
// [Ok]. Otherwise, it returns [Err] of res.
//
// The function f is not called when res is [Err].
func AndThen[T, U, E any](res Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	if res.IsOk() {
		return f(res.value)
	}

	return Err[U](res.err)
}

// Or returns res if it is [Ok]. Otherwise, it returns other.
//
// other is evaluated before Or is called. Use [OrElse] when the fallback
// [Result] should be computed only if res is [Err].
func Or[T, E, F any](res Result[T, E], other Result[T, F]) Result[T, F] {
	return OrElse(res, func(E) Result[T, F] { return other })
}

// OrElse returns res if it is [Ok]. Otherwise, it calls f with the contained
// error and returns the result.
//
// The function f is not called when res is [Ok].
func OrElse[T, E, F any](res Result[T, E], f func(E) Result[T, F]) Result[T, F] {
	if res.IsErr() {
		return f(res.err)
	}

	return Ok[T, F](res.value)
}
