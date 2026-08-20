package result

// Result represents a success value or an error. The zero value of Result is
// equivalent to [Err] of the zero value of E.
type Result[T, E any] struct {
	value T
	err   E
	ok    bool
}

// Err returns a [Result] containing the error.
func Err[T, E any](err E) Result[T, E] {
	return Result[T, E]{err: err, ok: false}
}

// Ok returns a [Result] containing the value.
func Ok[T, E any](value T) Result[T, E] {
	return Result[T, E]{value: value, ok: true}
}

// From returns [Ok] of v if err is nil. Otherwise, it returns [Err] of err.
func From[T any](v T, err error) Result[T, error] {
	if err != nil {
		return Err[T](err)
	}

	return Ok[T, error](v)
}
