package result

import (
	"fmt"
	"iter"
)

// Result represents a success value or an error. The zero value of Result is
// equivalent to [Err] of the zero value of E.
type Result[T, E any] struct {
	val T
	err E
	ok  bool
}

// Err returns a [Result] containing the error.
func Err[T, E any](err E) Result[T, E] {
	return Result[T, E]{err: err, ok: false}
}

// Ok returns a [Result] containing the value.
func Ok[T, E any](val T) Result[T, E] {
	return Result[T, E]{val: val, ok: true}
}

// From returns [Ok] of val if err is nil. Otherwise, it returns [Err] of err.
func From[T any](val T, err error) Result[T, error] {
	if err != nil {
		return Err[T](err)
	}

	return Ok[T, error](val)
}

// And returns other if r is [Ok]. Otherwise, it returns [Err] of r.
//
// other is evaluated before And is called. Use [Result.AndThen] when the second
// [Result] should be computed only if r is [Ok].
func (r Result[T, E]) And[U any](other Result[U, E]) Result[U, E] {
	return r.AndThen(func(T) Result[U, E] { return other })
}

// AndThen returns the result of calling f with the contained value if r is
// [Ok]. Otherwise, it returns [Err] of r.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) AndThen[U any](f func(T) Result[U, E]) Result[U, E] {
	if r.IsOk() {
		return f(r.val)
	}

	return Err[U](r.err)
}

// Expect returns the contained value. It panics with msg and the contained
// error if r is [Err].
func (r Result[T, E]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}

	return r.val
}

// ExpectErr returns the contained error. It panics with msg and the contained
// value if r is [Ok].
func (r Result[T, E]) ExpectErr(msg string) E {
	if r.IsOk() {
		panic(fmt.Sprintf("%s: %v", msg, r.val))
	}

	return r.err
}

// Get returns the contained value and true. If r is [Err], it returns the zero
// value of T and false.
func (r Result[T, E]) Get() (T, bool) {
	if r.IsErr() {
		var t T
		return t, false
	}

	return r.val, true
}

// GetErr returns the contained error and true. If r is [Ok], it returns the
// zero value of E and false.
func (r Result[T, E]) GetErr() (E, bool) {
	if r.IsOk() {
		var e E
		return e, false
	}

	return r.err, true
}

// Inspect calls f with the contained value if r is [Ok].
// It returns r unchanged.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) Inspect(f func(T)) Result[T, E] {
	if r.IsOk() {
		f(r.val)
	}

	return r
}

// InspectErr calls f with the contained error if r is [Err].
// It returns r unchanged.
//
// The function f is not called when r is [Ok].
func (r Result[T, E]) InspectErr(f func(E)) Result[T, E] {
	if r.IsErr() {
		f(r.err)
	}

	return r
}

// IsErr reports whether r is [Err].
func (r Result[T, E]) IsErr() bool {
	return !r.ok
}

// IsErrAnd reports whether r is [Err] and the contained error matches
// a predicate.
//
// The function f is not called when r is [Ok].
func (r Result[T, E]) IsErrAnd(f func(E) bool) bool {
	return !r.ok && f(r.err)
}

// IsOk reports whether r is [Ok].
func (r Result[T, E]) IsOk() bool {
	return r.ok
}

// IsOkAnd reports whether r is [Ok] and the contained value matches
// a predicate.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) IsOkAnd(f func(T) bool) bool {
	return r.ok && f(r.val)
}

// Map applies f to the contained value and returns the result as a [Result].
// If r is [Err], it returns [Err] of r.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) Map[U any](f func(T) U) Result[U, E] {
	return r.AndThen(func(val T) Result[U, E] { return Ok[U, E](f(val)) })
}

// MapErr applies f to the contained error and returns the result as a
// [Result]. If r is [Ok], it returns [Ok] of r.
//
// The function f is not called when r is [Ok].
func (r Result[T, E]) MapErr[F any](f func(E) F) Result[T, F] {
	return r.OrElse(func(err E) Result[T, F] { return Err[T](f(err)) })
}

// MapOr applies f to the contained value and returns the result. If r is
// [Err], it returns defaultVal.
//
// defaultVal is evaluated before MapOr is called. Use [Result.MapOrElse] to
// compute a fallback only when it is needed.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) MapOr[U any](defaultVal U, f func(T) U) U {
	return r.Map(f).UnwrapOr(defaultVal)
}

// MapOrDefault applies f to the contained value and returns the result. If r
// is [Err], it returns the zero value of U.
//
// The function f is not called when r is [Err].
func (r Result[T, E]) MapOrDefault[U any](f func(T) U) U {
	return r.Map(f).UnwrapOrDefault()
}

// MapOrElse applies okF to the contained value and returns the result. If r is
// [Err], it calls errF and returns that result.
//
// The function errF is not called when r is [Ok]. The function okF is not
// called when r is [Err].
func (r Result[T, E]) MapOrElse[U any](errF func(E) U, okF func(T) U) U {
	return r.Map(okF).UnwrapOrElse(errF)
}

// Or returns r if it is [Ok]. Otherwise, it returns other.
//
// other is evaluated before Or is called. Use [Result.OrElse] when the fallback
// [Result] should be computed only if r is [Err].
func (r Result[T, E]) Or[F any](other Result[T, F]) Result[T, F] {
	return r.OrElse(func(E) Result[T, F] { return other })
}

// OrElse returns r if it is [Ok]. Otherwise, it calls f with the contained
// error and returns the result.
//
// The function f is not called when r is [Ok].
func (r Result[T, E]) OrElse[F any](f func(E) Result[T, F]) Result[T, F] {
	if r.IsErr() {
		return f(r.err)
	}

	return Ok[T, F](r.val)
}

// Seq returns an iterator over the contained value. If r is [Err], the
// iterator is empty.
func (r Result[T, E]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if r.IsOk() && !yield(r.val) {
			return
		}
	}
}

// Unwrap returns the contained value. It panics if r is [Err].
//
// Because this method may panic, its use is generally discouraged. Panics are
// meant for unrecoverable errors, and may abort the entire program. Instead,
// use [Result.UnwrapOr], [Result.UnwrapOrElse], or [Result.UnwrapOrDefault].
func (r Result[T, E]) Unwrap() T {
	return r.Expect("result: Unwrap called on Err")
}

// UnwrapErr returns the contained error. It panics if r is [Ok].
//
// Because this method may panic, its use is generally discouraged. See
// [Result.Unwrap].
func (r Result[T, E]) UnwrapErr() E {
	return r.ExpectErr("result: UnwrapErr called on Ok")
}

// UnwrapOr returns the contained value or a provided default.
//
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the
// result of a function call, it is recommended to use [Result.UnwrapOrElse],
// which is lazily evaluated.
func (r Result[T, E]) UnwrapOr(defaultVal T) T {
	return r.UnwrapOrElse(func(E) T { return defaultVal })
}

// UnwrapOrDefault returns the contained value. If r is [Err], it returns the
// zero value of T.
func (r Result[T, E]) UnwrapOrDefault() T {
	return r.UnwrapOrElse(func(E) T {
		var t T
		return t
	})
}

// UnwrapOrElse returns the contained value. If r is [Err], it calls f with
// the contained error and returns its result. The function f is not called
// when r is [Ok].
func (r Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if r.IsErr() {
		return f(r.err)
	}

	return r.val
}

// Collect returns [Ok] of a slice of every contained value in rs. If any
// element is [Err], it returns [Err] of that error.
//
// If rs is empty, Collect returns [Ok] of an empty slice.
func Collect[T, E any](rs []Result[T, E]) Result[[]T, E] {
	vals := make([]T, 0, len(rs))
	for _, r := range rs {
		if r.IsErr() {
			return Err[[]T](r.err)
		}
		vals = append(vals, r.val)
	}

	return Ok[[]T, E](vals)
}

// Flatten unwraps a [Result] that contains another [Result].
func Flatten[T, E any](r Result[Result[T, E], E]) Result[T, E] {
	return r.AndThen(func(inner Result[T, E]) Result[T, E] { return inner })
}

// Unpack returns the contained value and a nil error if r is [Ok]. If r is
// [Err], it returns the zero value of T and the contained error.
//
// The contained error may be nil. A nil error on the [Err] path returns the
// zero value of T and a nil error, which is indistinguishable from [Ok] of the
// zero value. Use [Result.Get] to distinguish those cases.
func Unpack[T any](r Result[T, error]) (T, error) {
	if r.IsErr() {
		var t T
		return t, r.err
	}

	return r.val, nil
}
