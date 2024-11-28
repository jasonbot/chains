package chains

import (
	"iter"
)

func oneAtATime[T any](vals []T) iter.Seq2[T, []T] {
	// { 1 2 3} -> {[1, {2, 3}], [2, {1, 3}], [3, {1, 3}], ...}
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
	// { 1 2 3} -> {[1, {1, 2, 3}], [2, {1, 2, 3}], [3, {1, 2, 3}], ...}
	return func(yield func(T, []T) bool) {
		for index := range len(vals) {
			oneVal := vals[index]

			if !yield(oneVal, vals) {
				return
			}
		}
	}
}

func oneAtATimeTail[T any](vals []T) iter.Seq2[T, []T] {
	// { 1 2 3} -> {[1, {2, 3}], [2, {3}], ...}
	return func(yield func(T, []T) bool) {
		for index := range len(vals) {
			oneVal := vals[index]

			tailVal := vals[index+1:]
			if !yield(oneVal, tailVal) {
				return
			}
		}
	}
}

func combinationsandpermutations[T any](placementArray []T, vals []T, index int, length int, returnall bool, visit func([]T) iter.Seq2[T, []T], yield func([]T) bool) bool {
	if index >= length {
		return false
	}

	for item, rest := range visit(vals) {
		placementArray[index] = item

		if (index == length-1) || returnall {
			copyToReturn := make([]T, index+1)
			copy(copyToReturn, placementArray[:index+1])

			if !yield(copyToReturn) {
				return true
			}
		}

		if len(rest) > 0 && index < length {
			if combinationsandpermutations(placementArray, rest, index+1, length, returnall, visit, yield) {
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
		placement := make([]T, len(vals))
		combinationsandpermutations(placement, vals, 0, len(vals), true, oneAtATime, yield)
	}
}

// PermutationsOfLength will yield all combinations without replacement of
// a specified length
func PermutationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsandpermutations(placement, vals, 0, length, false, oneAtATime, yield)
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
		placement := make([]T, len(vals))
		combinationsandpermutations(placement, vals, 0, length, false, oneAtATimeWithReplacement, yield)
	}
}

// PermutationsWithReplacement will yield all possible orderings of the slice
func PermutationsWithReplacement[T any](vals []T) iter.Seq[[]T] {
	return PermutationsOfLengthWithReplacement(vals, len(vals))
}

// CombinationsOfLength will yield all combinations without replacement of
// a specified length
func CombinationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsandpermutations(placement, vals, 0, length, false, oneAtATimeTail, yield)
	}
}

// Combinations will yield all combinations without replacement of
// the entire slice
func Combinations[T any](vals []T) iter.Seq[[]T] {
	return CombinationsOfLength(vals, len(vals))
}
