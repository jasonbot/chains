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
// { 1 2 3 }, func(x) x + 2 -> { 3 4 5 }
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

// Accumulate takes an initial value, a reduce function, and an iterable
// and returns each result of applying the function iteratively.
// { 1 2 3 }, func(x, y) { return x + y } -> { 1 3 6 }
func Accumulate[T, V any](input iter.Seq[T], collectFunc func(V, T) V, zeroValue V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if !yield(zeroValue) {
			return
		}

		if collectFunc == nil || input == nil {
			return
		}

		for v := range input {
			zeroValue = collectFunc(zeroValue, v)

			if !yield(zeroValue) {
				return
			}
		}
	}
}

// ReduceWithZero takes an initial value, a reduce function, and an iterable
// and returns the final result of applying the function iteratively.
// { 1 2 3 }, func(x, y) { return x + y }, 100 -> 106
func ReduceWithZero[T, V any](input iter.Seq[T], collectFunc func(V, T) V, zeroValue V) V {
	return Last(Accumulate(input, collectFunc, zeroValue))
}

// Reduce takes a reduce function, and an iterable and returns the final result
// of applying the function iteratively.
// { 1 2 3 }, func(x, y) { return x + y } -> 6
func Reduce[T any](input iter.Seq[T], collectFunc func(T, T) T) T {
	zero, next := FirstAndRest(input)

	return ReduceWithZero(next, collectFunc, zero)
}

// Filter takes an iterator and only yields the items that pass the filter
// function check.
// { 1 2 3 4 5 }, func(i) i % 2 == 0 -> { 2 4 }
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

// Compact yields anything that is not the zero value.
// { 0 1 2 3 0 0 1 2 } -> { 1 2 3 1 2 }
func Compact[T comparable](input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var zeroValue T

		for v := range input {
			if v != zeroValue {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// All takes an iterator and returns true if the sequence is empty or all
// items match the predicate
// { 1 2 3 4 }, func(i) i == 1 -> false
// { 1 2 3 4 }, func(i) i != 5 -> true
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
// { 1 2 3 4 }, func(i) i % 2 == 0 -> true
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

// Windows returns a slice of up to length elements.
// { 1 2 3 4 5 }, 2 -> { 1 2 } { 3 4 } { 5 }
func Windows[T any](input iter.Seq[T], length int) iter.Seq[[]T] {
	items := make([]T, length)
	return func(yield func([]T) bool) {
		index := 0
		for v := range input {
			items[index] = v

			index += 1
			if index == length {
				if !yield(items) {
					return
				}
				index = 0
			}
		}
		if !yield(items[:index]) {
			return
		}
	}
}

// SlidingWindows returns a slice of up to length elements.
// { 1 2 3 4 5 }, 2 -> { 1 2 } { 2 3 } { 3 4 } { 4 5 }
func SlidingWindows[T any](input iter.Seq[T], length int) iter.Seq[[]T] {
	items := make([]T, length)
	return func(yield func([]T) bool) {
		index := 0
		for v := range input {
			copy(items, items[1:])
			items[length-1] = v

			index += 1
			if index >= length {
				if !yield(items) {
					return
				}
			}
		}
	}
}

// Zip takes two sequences and combines them into one (up to length of
// shortest)
// { 1 2 3 }, { "a" "b" "c" "d" } -> { ( 1, "a" ) ( 2, "b" ) ( 3, "c" ) }
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
// exhausted.
// { 1 2 3 }, { "a" "b" "c" "d" }, -1, "e"  -> { ( 1, "a" ) ( 2, "b" ) ( 3, "c" ) ( -1, "d" ) }
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
// { { 1 2 3 } { 4 5 6 } } -> { 1 2 3 4 5 6 }
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
// { 1 2 3 }, { 4 5 6 } -> { 1 2 3 4 5 6 }
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
// { 1 2 3 } -> iter{ 1 2 3 }, iter{ 1 2 3 }
func Tee[T any](in iter.Seq[T]) (iter.Seq[T], iter.Seq[T]) {
	var done1, done2 bool
	var exhausted bool

	iter1Queue := []T{}
	iter2Queue := []T{}

	next, done := iter.Pull(in)

	iter1 := func(yield func(T) bool) {
		defer done()

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
					iter2Queue = append(iter2Queue, nextval)
				}
			}

			nextval := iter1Queue[0]
			iter1Queue = iter1Queue[1:]

			if !yield(nextval) {
				done1 = true
				return
			}
		}
	}

	iter2 := func(yield func(T) bool) {
		defer done()

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
				iter2Queue = append(iter2Queue, nextval)
			}

			nextval := iter2Queue[0]
			iter2Queue = iter2Queue[1:]

			if !yield(nextval) {
				done2 = true
				return
			}
		}
	}

	return iter1, iter2
}

// Partition splits one iterator into two based on the predicate function
// { 1 2 3 4 5 6 7 8 9 10 }, func(x) { x % 2 == 0 } -> iter{ 2 4 6 8 10 }, iter{ 1 3 5 7 9 }
func Partition[T any](in iter.Seq[T], predicateFunction func(T) bool) (iter.Seq[T], iter.Seq[T]) {
	var done1, done2 bool
	var exhausted bool

	iter1Queue := []T{}
	iter2Queue := []T{}

	next, done := iter.Pull(in)

	iter1 := func(yield func(T) bool) {
		defer done()

		for {
			for len(iter1Queue) == 0 {
				if exhausted {
					return
				}

				nextval, ok := next()
				if !ok {
					exhausted = true
					return
				}

				fitsHere := predicateFunction(nextval)

				if fitsHere {
					iter1Queue = append(iter1Queue, nextval)
				} else if !done2 {
					iter2Queue = append(iter2Queue, nextval)
				}
			}

			nextval := iter1Queue[0]
			iter1Queue = iter1Queue[1:]

			if !yield(nextval) {
				done1 = true
				return
			}
		}
	}

	iter2 := func(yield func(T) bool) {
		defer done()

		for {
			for len(iter2Queue) == 0 {
				if exhausted {
					return
				}

				nextval, ok := next()
				if !ok {
					exhausted = true
					return
				}

				fitsHere := !predicateFunction(nextval)

				if fitsHere {
					iter2Queue = append(iter2Queue, nextval)
				} else if !done1 {
					iter1Queue = append(iter1Queue, nextval)
				}
			}

			nextval := iter2Queue[0]
			iter2Queue = iter2Queue[1:]

			if !yield(nextval) {
				done2 = true
				return
			}
		}
	}

	return iter1, iter2
}
