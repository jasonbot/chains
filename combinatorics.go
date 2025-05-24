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

func combinationsAndPermutations[T any](placementArray []T, vals []T, index int, length int, returnall bool, visit func([]T) iter.Seq2[T, []T], yield func([]T) bool) bool {
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
			if combinationsAndPermutations(placementArray, rest, index+1, length, returnall, visit, yield) {
				return true
			}
		}
	}

	return false
}

// AllPermutations will yield all permutations without replacement of
// every subset of items in the sequence of all length
// { 1 2 3 } -> {1, 2, 3} {1 3 2} {2 1 3} {2 3 1} {3 1 2} {3 2 1}
func AllPermutations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsAndPermutations(placement, vals, 0, len(vals), true, oneAtATime, yield)
	}
}

// PermutationsOfLength will yield all orderings without replacement of
// a specified length
// { 1 2 3 }, 2 -> { 1 2 } { 1 3 } { 2 1 } { 2 3 } { 3 1 } { 3 2 }
func PermutationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsAndPermutations(placement, vals, 0, length, false, oneAtATime, yield)
	}
}

// Permutations will yield all possible orderings of the slice
// { 1 2 3 } -> { 1 2 3 } { 1 3 2 } { 2 1 3 } { 2 3 1 } { 3 1 2 } { 3 2 1 }
func Permutations[T any](vals []T) iter.Seq[[]T] {
	return PermutationsOfLength(vals, len(vals))
}

// OrderedPermutations will yield all orderings without replacement of
// a specified length
// { 1 2 3 }, 2 -> { 1 2 } { 1 3 } { 2 3 }
func OrderedPermutations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		for i := range len(vals) {
			placement := make([]T, len(vals))
			combinationsAndPermutations(placement, vals, 0, i+1, false, oneAtATimeTail, yield)
		}
	}
}

// OrderedPermutationsOfLength will yield all orderings without replacement of
// a specified length
// { 1 2 3 }, 2 -> { 1 2 } { 1 3 } { 2 3 }
func OrderedPermutationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsAndPermutations(placement, vals, 0, length, false, oneAtATimeTail, yield)
	}
}

// PermutationsOfLengthWithReplacement will yield all combinations with replacement of
// a specified length
// { 1 2 3 }, 2 -> { 1 2 } { 1 3 } { 2 1 } { 2 3 } { 3 1 } { 3 2 }
func PermutationsOfLengthWithReplacement[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsAndPermutations(placement, vals, 0, length, false, oneAtATimeWithReplacement, yield)
	}
}

// PermutationsWithReplacement will yield all possible orderings of the slice with replacement
// { 1 2 3 } -> { 1 1 1 } { 1 1 2 } { 1 1 3 } { 2 1 1 } ... { 3 3 3 }
func PermutationsWithReplacement[T any](vals []T) iter.Seq[[]T] {
	return PermutationsOfLengthWithReplacement(vals, len(vals))
}

// CombinationsOfLength will yield all combinations with replacement of
// a specified length
// { 1 2 3 }, 2 ->  { 1 } { 1 2 } { 1 3 } { 2 } { 2 1 } { 2 3 } { 3 } { 3 1 } { 3 2 }
func CombinationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		placement := make([]T, len(vals))
		combinationsAndPermutations(placement, vals, 0, length, false, oneAtATimeTail, yield)
	}
}

// Combinations will yield all combinations with replacement of
// the entire slice
// { 1 2 3 } ->  { 1 } { 1 2 } { 1 2 3 } { 1 3 } { 1 3 2 } { 2 } { 2 1 } { 2 1 3 } { 2 3 } { 2 3 1 } { 3 } { 3 1 } { 3 1 2 } { 3 2 } { 3 2 1 }
func Combinations[T any](vals []T) iter.Seq[[]T] {
	return CombinationsOfLength(vals, len(vals))
}

// Pairwise will yield all possible combinations of the two iterators
// { "a" "b" }, { 1 2 } -> iter{ ( "a" 1 ) ( "a" 2 ) ( "b" 1 ) ( "b" 2 ) }
func Pairwise[T, V any](s1 iter.Seq[T], s2 iter.Seq[V]) iter.Seq2[T, V] {
	i1, i2 := Tee(s1)

	return Zip(Cycle(i1), Lengthen(s2, Count(i2)))
}
