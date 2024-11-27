package chains

import (
	"iter"
)

// Each wraps a slice as an iterable. Only iteresting when applying
// higher-level functions like Map or Filter.
func Each[T any](in []T) iter.Seq[T] {
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

// Map takes an iterator and applies a function to each element.
func Map[T, V any](mapFunc func(T) V, input iter.Seq[T]) iter.Seq[V] {
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
	zero, next := FirstAndRest(input)

	return ReduceWithZero(collectFunc, zero, next)
}

// Filter takes an iterator and only yields the items that pass the filter
// function check.
func Filter[T any](filterFunc func(T) bool, input iter.Seq[T]) iter.Seq[T] {
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

// All takes an iterator and returns true if the sequence is empty or all
// items match the predicate
func All[T any](filterFunc func(T) bool, input iter.Seq[T]) bool {
	if filterFunc == nil || input == nil {
		return true
	}

	for v := range input {
		if !filterFunc(v) {
			return false
		}
	}

	return true
}

// Any takes an iterator and returns true if the sequence is empty or any
// item matches the predicate
func Any[T any](filterFunc func(T) bool, input iter.Seq[T]) bool {
	if filterFunc == nil || input == nil {
		return true
	}

	for v := range input {
		if filterFunc(v) {
			return true
		}
	}

	return false
}

// Tap visits each item with the visitor function but passes each item along.
func Tap[T any](visitor func(T), input iter.Seq[T]) iter.Seq[T] {
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

// Zip takes two sequences and combines them into one (up to length of
// shortest) via zipfunc
func Zip[T, V, K any](zipFunc func(T, V) K, input1 iter.Seq[T], input2 iter.Seq[V]) func(func(K) bool) {
	return func(yield func(K) bool) {
		nextOne, oneDone := iter.Pull(input1)
		defer oneDone()

		nextTwo, twoDone := iter.Pull(input2)
		defer twoDone()

		var one T
		var two V
		var ok bool

		for {
			if one, ok = nextOne(); !ok {
				return
			}
			if two, ok = nextTwo(); !ok {
				return
			}

			yield(zipFunc(one, two))
		}
	}
}

// ZipLongest takes two sequences and combines them into one (up to length
// of longest) via zipfunc, using fillerOne/fillerTwo as defaults if one is
// exhausted
func ZipLongest[T, V, K any](zipFunc func(T, V) K, fillerOne T, fillerTwo V, input1 iter.Seq[T], input2 iter.Seq[V]) func(func(K) bool) {
	return func(yield func(K) bool) {
		nextOne, oneDone := iter.Pull(input1)
		defer oneDone()

		nextTwo, twoDone := iter.Pull(input2)
		defer twoDone()

		var one T
		var two V
		var ok, oneDoneProcessing, twoDoneProcessing bool

		for !(oneDoneProcessing && twoDoneProcessing) {
			if !oneDoneProcessing {
				if one, ok = nextOne(); !ok {
					oneDoneProcessing = true
					one = fillerOne
				}
			}
			if !twoDoneProcessing {
				if two, ok = nextTwo(); !ok {
					twoDoneProcessing = true
					two = fillerTwo
				}
			}

			if !yield(zipFunc(one, two)) {
				return
			}
		}
	}
}

// Flatten takes any number of iterables and combines them into one
func Flatten[T any](sequences iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for seq := range sequences {
			for item := range seq {
				yield(item)
			}
		}
	}
}

// FlattenArgs takes any number of iterable args and combines them into one
func FlattenArgs[T any](sequences ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range sequences {
			for item := range seq {
				yield(item)
			}
		}
	}
}
