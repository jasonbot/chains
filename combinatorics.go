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
			if !exhausted && includeall && !yield(currentordering) {
				return true
			}

			exhausted = exhausted || permutations(currentordering, rest, length-1, yield, includeall, exhausted)
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
			if !exhausted && includeall && !yield(currentordering) {
				return true
			}

			exhausted = exhausted || permutationsWithReplacement(currentordering, rest, length-1, yield, includeall, exhausted)
		} else {
			if exhausted || !yield(currentordering) {
				return true
			}
		}
	}
	return false
}

// AllPermutations will yield all permutations without replacement of
// every subset of items in the sequence of all length
func AllPermutations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		permutations([]T{}, vals, len(vals), yield, true, false)
	}
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

func oneAtATimeTail[T any](vals []T) iter.Seq2[T, []T] {
	return func(yield func(T, []T) bool) {
		for index := range len(vals) {
			oneVal := vals[index]

			// Need to copy slice so we don't overwrite it
			tmp := make([]T, len(vals))
			copy(tmp, vals)

			tailVal := vals[index+1:]
			if !yield(oneVal, tailVal) {
				return
			}
		}
	}
}

func combinations[T any](placementArray []T, vals []T, index int, length int, yield func([]T) bool) bool {
	if index >= length {
		return false
	}

	for item, rest := range oneAtATimeTail(vals) {
		placementArray[index] = item

		if index == length-1 {
			if !yield(placementArray) {
				return true
			}
		}

		if len(rest) > 0 && index < length {
			if combinations(placementArray, rest, index+1, length, yield) {
				return true
			}
		}
	}

	return false
}

// CombinationsOfLength will yield all combinations without replacement of
// a specified length
func CombinationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		endArray := make([]T, length)
		combinations(endArray, vals, 0, length, yield)
	}
}

// Combinations will yield all combinations without replacement of
// the entire slice
func Combinations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		endArray := make([]T, len(vals))
		combinations(endArray, vals, 0, len(vals), yield)
	}
}
