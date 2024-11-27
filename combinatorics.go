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

func oneAtATimeWithReplacement[T any](vals []T) iter.Seq2[T, []T] {
	return func(yield func(T, []T) bool) {
		for index := range len(vals) {
			oneVal := vals[index]

			tmp := make([]T, len(vals))
			copy(tmp, vals)

			if !yield(oneVal, tmp) {
				return
			}
		}
	}
}

func permutations[T any](sofar, vals []T, length int, yield func([]T) bool) {
	for val, rest := range oneAtATime(vals) {
		currentordering := append(sofar, val)
		if length > 0 && len(rest) > 0 {
			permutations(currentordering, rest, length-1, yield)
		} else {
			if !yield(currentordering) {
				return
			}
		}
	}
}

func permutationsWithReplacement[T any](sofar, vals []T, length int, yield func([]T) bool) {
	currentordering := make([]T, len(sofar)+1)
	copy(currentordering, sofar)

	for val, rest := range oneAtATimeWithReplacement(vals) {
		currentordering[len(currentordering)-1] = val
		if length > 1 {
			permutationsWithReplacement(currentordering, rest, length-1, yield)
		} else {
			if !yield(currentordering) {
				return
			}
		}
	}
}

// PermutationsOfLength will yield all combinations without replacement of
// a specified length
func PermutationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutations([]T{}, vals, length, yield)
	}
}

// Permutations will yield all possible orderings of the slice
func Permutations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutations([]T{}, vals, len(vals), yield)
	}
}

// PermutationsOfLengthWithReplacement will yield all combinations without replacement of
// a specified length
func PermutationsOfLengthWithReplacement[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutationsWithReplacement([]T{}, vals, length, yield)
	}
}

// PermutationsWithReplacement will yield all possible orderings of the slice
func PermutationsWithReplacement[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutationsWithReplacement([]T{}, vals, len(vals), yield)
	}
}
