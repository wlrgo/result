package result

import "iter"

// Collect returns a slice of every contained value in results. If any element
// is [Err], it returns that error.
//
// If results is empty, Collect returns [Ok] of an empty slice.
func Collect[T, E any](results []Result[T, E]) Result[[]T, E] {
	values := make([]T, 0, len(results))
	for _, res := range results {
		if res.IsErr() {
			return Err[[]T](res.err)
		}
		values = append(values, res.value)
	}

	return Ok[[]T, E](values)
}

// Seq returns an iterator over the contained value. If res is [Err], the
// iterator is empty.
func (res Result[T, E]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if res.IsOk() && !yield(res.value) {
			return
		}
	}
}
