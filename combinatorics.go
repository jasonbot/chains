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

func combinations[T any](sofar, vals []T, length int, yield func([]T) bool) {
	for val, rest := range oneAtATime(vals) {
		currentcombination := append(sofar, val)
		if length > 0 && len(rest) > 0 {
			combinations(currentcombination, rest, length-1, yield)
		} else {
			if !yield(currentcombination) {
				return
			}
		}
	}
}

func CombinationsOfLength[T any](vals []T, length int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		combinations([]T{}, vals, length, yield)
	}
}

func Combinations[T any](vals []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		combinations([]T{}, vals, len(vals), yield)
	}
}
