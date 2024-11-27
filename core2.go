package chains

import (
	"iter"
)

// Map2 takes an iterator and applies a function to each element.
func Map2[T, V, K any](input iter.Seq2[T, V], mapFunc func(T, V) K) iter.Seq[K] {
	return func(yield func(K) bool) {
		if mapFunc == nil || input == nil {
			return
		}

		for v, k := range input {
			if !yield(mapFunc(v, k)) {
				return
			}
		}
	}
}

// Filter2 takes an iterator and only yields the items that pass the filter
// function check.
func Filter2[T, V any](input iter.Seq2[T, V], predicateFunc func(T, V) bool) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		if predicateFunc == nil || input == nil {
			return
		}

		for t, v := range input {
			if predicateFunc(t, v) {
				if !yield(t, v) {
					return
				}
			}
		}
	}
}

// All2 takes an iterator and returns true if the sequence is empty or all
// items match the predicate
func All2[T, V any](input iter.Seq2[T, V], predicateFunc func(T, V) bool) bool {
	if predicateFunc == nil || input == nil {
		return true
	}

	for t, v := range input {
		if !predicateFunc(t, v) {
			return false
		}
	}

	return true
}

// Any2 takes an iterator and returns true if the sequence is empty or any
// item matches the predicate
func Any2[T, V any](input iter.Seq2[T, V], predicateFunc func(T, V) bool) bool {
	if predicateFunc == nil || input == nil {
		return true
	}

	for t, v := range input {
		if predicateFunc(t, v) {
			return true
		}
	}

	return false
}

// Tap2 visits each item with the visitor function but passes each item along.
func Tap2[T, V any](input iter.Seq2[T, V], visitor func(T, V)) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		if visitor == nil {
			for t, v := range input {
				if !yield(t, v) {
					return
				}
			}

			return
		}

		for t, v := range input {
			visitor(t, v)
			if !yield(t, v) {
				return
			}
		}
	}
}

// Flatten takes any number of iterables and combines them into one
func Flatten2[T, V any](sequences iter.Seq[iter.Seq2[T, V]]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		for seq := range sequences {
			for t, v := range seq {
				if !yield(t, v) {
					return
				}
			}
		}
	}
}

// FlattenArgs takes any number of iterable args and combines them into one
func FlattenArgs2[T, V any](sequences ...iter.Seq2[T, V]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		for _, seq := range sequences {
			for t, v := range seq {
				if !yield(t, v) {
					return
				}
			}
		}
	}
}
