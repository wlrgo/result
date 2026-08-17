package result

import "iter"

// Collect
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

// Seq
func (res Result[T, E]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if res.IsOk() {
			yield(res.value)
		}
	}
}
