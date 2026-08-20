package result

import "fmt"

// Expect returns the contained value. It panics with msg and the contained
// error if the result is [Err].
func (res Result[T, E]) Expect(msg string) T {
	if res.IsErr() {
		panic(fmt.Sprintf("%s: %v", msg, res.err))
	}

	return res.value
}

// Unwrap returns the contained value. It panics if the result is [Err].
//
// Because this function may panic, its use is generally discouraged. Panics
// are meant for unrecoverable errors, and may abort the entire program.
// Instead, use [Result.UnwrapOr], [Result.UnwrapOrElse], or
// [Result.UnwrapOrDefault].
func (res Result[T, E]) Unwrap() T {
	return res.Expect("result: Unwrap called on Err")
}

// UnwrapOr returns the contained value or a provided default.
//
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the
// result of a function call, it is recommended to use [Result.UnwrapOrElse],
// which is lazily evaluated.
func (res Result[T, E]) UnwrapOr(value T) T {
	return res.UnwrapOrElse(func(E) T { return value })
}

// UnwrapOrDefault returns the contained value. If the result is [Err], it
// returns the zero value of T.
func (res Result[T, E]) UnwrapOrDefault() T {
	return res.UnwrapOrElse(func(E) T {
		var t T
		return t
	})
}

// UnwrapOrElse returns the contained value. If the result is [Err], it calls
// f with the contained error and returns its result. The function f is not
// called when the result is [Ok].
func (res Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if res.IsErr() {
		return f(res.err)
	}

	return res.value
}

// ExpectErr returns the contained error. It panics with msg and the contained
// value if the result is [Ok].
func (res Result[T, E]) ExpectErr(msg string) E {
	if res.IsOk() {
		panic(fmt.Sprintf("%s: %v", msg, res.value))
	}

	return res.err
}

// UnwrapErr returns the contained error. It panics if the result is [Ok].
//
// Because this function may panic, its use is generally discouraged. See
// [Result.Unwrap].
func (res Result[T, E]) UnwrapErr() E {
	return res.ExpectErr("result: UnwrapErr called on Ok")
}
