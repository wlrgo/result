package result

// Flatten unwraps a [Result] that contains another [Result].
func Flatten[T, E any](res Result[Result[T, E], E]) Result[T, E] {
	return AndThen(res, func(inner Result[T, E]) Result[T, E] { return inner })
}

// Map applies f to the contained value and returns the result as a [Result].
// If res is [Err], it returns that error.
//
// The function f is not called when res is [Err].
func Map[T, U, E any](res Result[T, E], f func(T) U) Result[U, E] {
	return AndThen(res, func(val T) Result[U, E] { return Ok[U, E](f(val)) })
}

// MapErr applies f to the contained error and returns the result as a
// [Result]. If res is [Ok], it returns that value.
//
// The function f is not called when res is [Ok].
func MapErr[T, E, F any](res Result[T, E], f func(E) F) Result[T, F] {
	return OrElse(res, func(err E) Result[T, F] { return Err[T](f(err)) })
}

// MapOr applies f to the contained value and returns the result. If res is
// [Err], it returns defaultValue.
//
// defaultValue is evaluated before MapOr is called. Use [MapOrElse] to compute
// a fallback only when it is needed.
//
// The function f is not called when res is [Err].
func MapOr[T, U, E any](res Result[T, E], defaultValue U, f func(T) U) U {
	return Map(res, f).UnwrapOr(defaultValue)
}

// MapOrDefault applies f to the contained value and returns the result. If res
// is [Err], it returns the zero value of U.
//
// The function f is not called when res is [Err].
func MapOrDefault[T, U, E any](res Result[T, E], f func(T) U) U {
	return Map(res, f).UnwrapOrDefault()
}

// MapOrElse applies okF to the contained value and returns the result. If res
// is [Err], it calls errF and returns that result.
//
// The function errF is not called when res is [Ok]. The function okF is not
// called when res is [Err].
func MapOrElse[T, U, E any](res Result[T, E], errF func(E) U, okF func(T) U) U {
	return Map(res, okF).UnwrapOrElse(errF)
}

// Inspect calls f with the contained value if res is [Ok].
// It returns res unchanged.
//
// The function f is not called when res is [Err].
func (res Result[T, E]) Inspect(f func(T)) Result[T, E] {
	if res.IsOk() {
		f(res.value)
	}

	return res
}

// InspectErr calls f with the contained error if res is [Err].
// It returns res unchanged.
//
// The function f is not called when res is [Ok].
func (res Result[T, E]) InspectErr(f func(E)) Result[T, E] {
	if res.IsErr() {
		f(res.err)
	}

	return res
}
