package result

import "fmt"

// Expect
func (res Result[T, E]) Expect(msg string) T {
	if res.IsErr() {
		panic(fmt.Sprintf("%s: %v", msg, res.err))
	}

	return res.value
}

// Unwrap
func (res Result[T, E]) Unwrap() T {
	return res.Expect("result: Unwrap called on Err")
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
		panic(fmt.Sprintf("%s: %v", msg, res.value))
	}

	return res.err
}

// UnwrapErr
func (res Result[T, E]) UnwrapErr() E {
	return res.ExpectErr("result: UnwrapErr called on Ok")
}
