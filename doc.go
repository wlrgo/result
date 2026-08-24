// Package result provides a generic Result type for representing success or
// failure.
//
// It follows the Rust Result API as closely as Go allows, rather than
// reshaping it into an idiomatic Go design. It also provides Go helpers such
// as [From], [Result.Get], [Result.GetErr], and [Unpack] where Rust has no
// equivalent.
//
// # Construction
//
// The zero value of [Result] is [Err] of the zero value of E. Use [Ok], [Err],
// or [From] to construct a [Result] from a value, an error, or a Go (T, error)
// pair.
//
// # Methods and functions
//
// Operations on a single [Result] are methods, including those that introduce
// another type, such as [Result.And], [Result.AndThen], [Result.Map], and
// [Result.Or]. Package-level functions are used for construction ([Ok], [Err],
// [From]), conversion ([Unpack]), collecting ([Collect]), flattening a nested
// [Result] ([Flatten]), and operations that need an extra constraint, such as
// [Compare] and [Equal].
//
// # API
//
// Querying: [Result.IsOk], [Result.IsErr], [Result.IsOkAnd],
// [Result.IsErrAnd].
//
// Extracting: [Result.Expect], [Result.Unwrap], [Result.UnwrapOr],
// [Result.UnwrapOrDefault], [Result.UnwrapOrElse], [Result.ExpectErr],
// [Result.UnwrapErr].
//
// Combining: [Result.And], [Result.AndThen], [Result.Or], [Result.OrElse].
//
// Transforming: [Flatten], [Result.Map], [Result.MapErr], [Result.MapOr],
// [Result.MapOrDefault], [Result.MapOrElse], [Result.Inspect],
// [Result.InspectErr].
//
// Comparing: [Compare], [Equal], [Ge], [Gt], [Le], [Lt]. [Ok] is less than
// any [Err].
//
// Iterating: [Result.Seq] yields the contained value if the result is [Ok].
// [Collect] turns a slice of results into a result of a slice.
//
// Converting: [Result.Get], [Result.GetErr], [Unpack].
//
// # Evaluation
//
// Combinators that take a fallback value evaluate it before the call:
// [Result.And], [Result.Or], [Result.UnwrapOr], and [Result.MapOr]. Use the
// Else variants to compute a fallback only when it is needed.
//
// # Unpacking
//
// [Unpack] returns a Go (T, error) pair. A nil error on the [Err] path
// returns the zero value of T and a nil error, which is indistinguishable from
// [Ok] of the zero value. Use [Result.Get] to distinguish those cases.
package result
