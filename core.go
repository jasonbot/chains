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
func Map[T, V any](input iter.Seq[T], mapFunc func(T) V) iter.Seq[V] {
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
func ReduceWithZero[T, V any](input iter.Seq[T], collectFunc func(V, T) V, zeroValue V) V {
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
func Reduce[T any](input iter.Seq[T], collectFunc func(T, T) T) T {
	zero, next := FirstAndRest(input)

	return ReduceWithZero(next, collectFunc, zero)
}

// Filter takes an iterator and only yields the items that pass the filter
// function check.
func Filter[T any](input iter.Seq[T], predicateFunc func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		if predicateFunc == nil || input == nil {
			return
		}

		for v := range input {
			if predicateFunc(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// All takes an iterator and returns true if the sequence is empty or all
// items match the predicate
func All[T any](input iter.Seq[T], predicateFunc func(T) bool) bool {
	if predicateFunc == nil || input == nil {
		return true
	}

	for v := range input {
		if !predicateFunc(v) {
			return false
		}
	}

	return true
}

// Any takes an iterator and returns true if the sequence is empty or any
// item matches the predicate
func Any[T any](input iter.Seq[T], predicateFunc func(T) bool) bool {
	if predicateFunc == nil || input == nil {
		return true
	}

	for v := range input {
		if predicateFunc(v) {
			return true
		}
	}

	return false
}

// Tap visits each item with the visitor function but passes each item along.
func Tap[T any](input iter.Seq[T], visitor func(T)) iter.Seq[T] {
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
// shortest)
func Zip[T, V any](input1 iter.Seq[T], input2 iter.Seq[V]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
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

			if !yield(one, two) {
				return
			}
		}
	}
}

// ZipLongest takes two sequences and combines them into one (up to length
// of longest) via zipfunc, using fillerOne/fillerTwo as defaults if one is
// exhausted
func ZipLongest[T, V any](input1 iter.Seq[T], input2 iter.Seq[V], fillerOne T, fillerTwo V) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
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

			if !yield(one, two) {
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
				if !yield(item) {
					return
				}
			}
		}
	}
}

// FlattenArgs takes any number of iterable args and combines them into one
func FlattenArgs[T any](sequences ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range sequences {
			for item := range seq {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Tee splits one iterator into two
func Tee[T any](in iter.Seq[T]) (iter.Seq[T], iter.Seq[T]) {
	var done1, done2 bool
	var exhausted bool

	iter1Queue := []T{}
	iter2Queue := []T{}

	next, done := iter.Pull(in)
	defer done()

	iter1 := func(yield func(T) bool) {
		for {
			if len(iter1Queue) == 0 {
				if exhausted {
					return
				}

				nextval, ok := next()
				if !ok {
					exhausted = true
					return
				}

				iter1Queue = append(iter1Queue, nextval)
				if !done2 {
					iter2Queue = append(iter1Queue, nextval)
				}
			} else {
				nextval := iter1Queue[0]
				iter1Queue = iter1Queue[1:]

				if !yield(nextval) {
					done1 = true
					return
				}
			}
		}
	}

	iter2 := func(yield func(T) bool) {
		for {
			if len(iter2Queue) == 0 {
				if exhausted {
					return
				}

				nextval, ok := next()
				if !ok {
					exhausted = true
					return
				}

				if !done1 {
					iter1Queue = append(iter1Queue, nextval)
				}
				iter2Queue = append(iter1Queue, nextval)
			} else {
				nextval := iter2Queue[0]
				iter1Queue = iter2Queue[1:]

				if !yield(nextval) {
					done2 = true
					return
				}
			}
		}
	}

	return iter1, iter2
}
