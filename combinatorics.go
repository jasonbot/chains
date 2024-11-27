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

func permutations[T any](sofar, vals []T, length int, yield func([]T) bool, includeall bool, exhausted bool) bool {
	for val, rest := range oneAtATime(vals) {
		if exhausted {
			return exhausted
		}

		currentordering := append(sofar, val)
		if length > 1 && len(rest) > 0 {
			exhausted = exhausted || permutations(currentordering, rest, length-1, yield, includeall, exhausted)

			if !exhausted && includeall && !yield(currentordering) {
				return true
			}
		} else {
			if exhausted || !yield(currentordering) {
				return true
			}
		}
	}
	return false
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

func permutationsWithReplacement[T any](sofar, vals []T, length int, yield func([]T) bool, includeall bool, exhausted bool) bool {
	currentordering := make([]T, len(sofar)+1)
	copy(currentordering, sofar)

	if length == 0 {
		return exhausted
	}

	for val, rest := range oneAtATimeWithReplacement(vals) {
		if exhausted {
			return exhausted
		}

		currentordering[len(currentordering)-1] = val
		if length > 1 {
			exhausted = exhausted || permutationsWithReplacement(currentordering, rest, length-1, yield, includeall, exhausted)

			if !exhausted && includeall && !yield(currentordering) {
				return true
			}
		} else {
			if exhausted || !yield(currentordering) {
				return true
			}
		}
	}
	return false
}

// PermutationsOfLength will yield all combinations without replacement of
// a specified length
func PermutationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutations([]T{}, vals, length, yield, false, false)
	}
}

// Permutations will yield all possible orderings of the slice
func Permutations[T any](vals []T) iter.Seq[[]T] {
	return PermutationsOfLength(vals, len(vals))
}

// PermutationsOfLengthWithReplacement will yield all combinations without replacement of
// a specified length
func PermutationsOfLengthWithReplacement[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutationsWithReplacement([]T{}, vals, length, yield, false, false)
	}
}

// PermutationsWithReplacement will yield all possible orderings of the slice
func PermutationsWithReplacement[T any](vals []T) iter.Seq[[]T] {
	return PermutationsOfLengthWithReplacement(vals, len(vals))
}
