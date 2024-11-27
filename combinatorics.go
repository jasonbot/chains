package chains

import (
	"iter"
)

func oneAtATime[T any](vals []T) iter.Seq2[T, []T] {
	return func(yield func(T, []T) bool) {
		for index := range len(vals) {
			oneVal := vals[index]

			// Need to copy slice so we don't overwrite it
			tmp := make([]T, len(vals))
			copy(tmp, vals)

			tailVal := tmp[:index+copy(tmp[index:], tmp[index+1:])]
			if !yield(oneVal, tailVal) {
				return
			}
		}
	}
}

func orderings[T any](sofar, vals []T, length int, yield func([]T) bool) {
	for val, rest := range oneAtATime(vals) {
		currentordering := append(sofar, val)
		if length > 0 && len(rest) > 0 {
			orderings(currentordering, rest, length-1, yield)
		} else {
			if !yield(currentordering) {
				return
			}
		}
	}
}

// OrderingsOfLength will yield all combinations without replacement of
// a specified length
func OrderingsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		orderings([]T{}, vals, length, yield)
	}
}

// Orderings will yield all possible orderings of the slice
func Orderings[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		orderings([]T{}, vals, len(vals), yield)
	}
}
