package result

// Result
type Result[T, E any] struct {
	value T
	err   E
	ok    bool
}

// Err
func Err[T, E any](err E) Result[T, E] {
	return Result[T, E]{err: err, ok: false}
}

// Ok
func Ok[T, E any](value T) Result[T, E] {
	return Result[T, E]{value: value, ok: true}
}

// From
func From[T any](v T, err error) Result[T, error] {
	if err != nil {
		return Err[T](err)
	}

	return Ok[T, error](v)
}
