package result

// Flatten
func Flatten[T, E any](res Result[Result[T, E], E]) Result[T, E] {
	return AndThen(res, func(inner Result[T, E]) Result[T, E] { return inner })
}

// Map
func Map[T, U, E any](res Result[T, E], f func(T) U) Result[U, E] {
	return AndThen(res, func(val T) Result[U, E] { return Ok[U, E](f(val)) })
}

// MapErr
func MapErr[T, E, F any](res Result[T, E], f func(E) F) Result[T, F] {
	return OrElse(res, func(err E) Result[T, F] { return Err[T](f(err)) })
}

// MapOr
func MapOr[T, U, E any](res Result[T, E], defaultValue U, f func(T) U) U {
	return Map(res, f).UnwrapOr(defaultValue)
}

// MapOrDefault
func MapOrDefault[T, U, E any](res Result[T, E], f func(T) U) U {
	return Map(res, f).UnwrapOrDefault()
}

// MapOrElse
func MapOrElse[T, U, E any](res Result[T, E], errF func(E) U, okF func(T) U) U {
	return Map(res, okF).UnwrapOrElse(errF)
}

// Inspect
func (res Result[T, E]) Inspect(f func(T)) Result[T, E] {
	if res.IsOk() {
		f(res.value)
	}

	return res
}

// InspectErr
func (res Result[T, E]) InspectErr(f func(E)) Result[T, E] {
	if res.IsErr() {
		f(res.err)
	}

	return res
}
