package chains

import (
	"iter"
)

// Each wraps a slice as an iterable. Only iteresting when applying
// higher-level functions like Map or Filter.
func Each[T any](in []T) func(func(T) bool) {
	return func(yield func(T) bool) {
		if in == nil {
			return
		}

		for _, v := range in {
			if !yield(v) {
				return
			}
		}
	}
}

// Map takes ai iterator and applies a function to each element.
func Map[T, V any](mapFunc func(T) V, input iter.Seq[T]) func(func(V) bool) {
	return func(yield func(V) bool) {
		if mapFunc == nil || input == nil {
			return
		}

		for v := range input {
			if !yield(mapFunc(v)) {
				return
			}
		}
	}
}

// ReduceWithZero takes an initial value, a reduce function, and an iterable
// and returns the final result of applying the function iteratively.
func ReduceWithZero[T, V any](collectFunc func(V, T) V, zeroValue V, input iter.Seq[T]) V {
	if collectFunc == nil || input == nil {
		var zeroValue V
		return zeroValue
	}

	for v := range input {
		zeroValue = collectFunc(zeroValue, v)
	}
	return zeroValue
}

// Reduce takes a reduce function, and an iterable and returns the final result
// of applying the function iteratively.
func Reduce[T any](collectFunc func(T, T) T, input iter.Seq[T]) T {
	zero, next := Car(input)

	return ReduceWithZero(collectFunc, zero, next)
}

// Filter takes an iterator and only yields the items that pass the filter
// function check.
func Filter[T any](filterFunc func(T) bool, input iter.Seq[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		if filterFunc == nil || input == nil {
			return
		}

		for v := range input {
			if filterFunc(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Tap visits each item with the visitor function but passes each item along.
func Tap[T any](visitor func(T), input iter.Seq[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		if visitor == nil {
			for v := range input {
				if !yield(v) {
					return
				}
			}

			return
		}

		for v := range input {
			visitor(v)
			if !yield(v) {
				return
			}
		}
	}
}
