package chains

import "iter"

// Count returns the length of the exhusted iterator.
func Count[T any](input iter.Seq[T]) int {
	var count int

	for range input {
		count += 1
	}

	return count
}

// Uniq returns the first value in a sequence, omitting duplicates.
// If a duplicate shows up further in the sequence, it will show up again.
// For example, {1 1 2 2 3 3 4} will yield {1 2 3 4} but
// {1 1 2 2 1 1 4} will yield {1 2 1 4}
func Uniq[T comparable](offset int, input iter.Seq[T]) iter.Seq[T] {
	first := true
	var lastValue T

	return func(yield func(T) bool) {
		for item := range input {
			if first || (item != lastValue) {
				if !yield(item) {
					return
				}
			}
			lastValue = item
			first = false
		}
	}
}

// Cycle yields every item in the sequence indefinitely, starting from the
// beginning once exhausted. Iamgine an unbound Repeat.
func Cycle[T any](input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		items := make([]T, 0, 100)

		for item := range input {
			if !yield(item) {
				return
			}
			items = append(items, item)
		}

		var index int
		for {
			if !yield(items[index]) {
				return
			}
			index += 1
			index %= len(items)
		}
	}
}

// Repeat will yield every element of the provided sequence up to repeats
// times; 3, {1 2 3} -> { 1 2 3 1 2 3 1 2 3}
func Repeat[T any](repeats int, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		items := make([]T, 0, 100)

		for item := range input {
			if !yield(item) {
				return
			}
			items = append(items, item)
		}

		var index int
		for {
			if index == 0 {
				repeats -= 1
				if repeats <= 0 {
					return
				}
			}
			if !yield(items[index]) {
				return
			}
			index += 1
			index %= len(items)
		}
	}
}

// Lengthen will yield every element of the provided sequence up to repeats
// times; 3, {1 2 3} -> { 1 1 1 2 2 2 3 3 3 }
func Lengthen[T any](repeats int, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range input {
			for range repeats {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Rotate will yield every element of the provided sequence, rotating the first
// element to the end; { 1 2 3 } -> { 2 3 1 }
func Rotate[T any](repeats int, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		first := true
		var firstItem T

		for item := range input {
			if first {
				first = false
				firstItem = item
			} else {
				if !yield(item) {
					return
				}
			}
		}

		yield(firstItem)
	}
}

// PastOffset starts iterating at the zero-based index
func PastOffset[T any](offset int, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		index := 0

		for item := range input {
			if index >= offset {
				if !yield(item) {
					return
				}
			}
			index += 1
		}
	}
}

// UntilOffset stops iterating at the zero-based index
func UntilOffset[T any](offset int, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		index := 0

		for item := range input {
			if index < offset {
				if !yield(item) {
					return
				}
				index += 1
			} else {
				return
			}
		}
	}
}

// TakeWhile stops iterating once the filterFunc returns false
func TakeWhile[T any](filterFunc func(T) bool, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range input {
			if !filterFunc(item) || !yield(item) {
				return
			}
		}
	}
}

// TakeWhile does not start iterating until the filterFunc returns true the first time
func TakeAfter[T any](filterFunc func(T) bool, input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		found := false

		for item := range input {
			if found {
				if !yield(item) {
					return
				}
			} else if filterFunc(item) {
				if !yield(item) {
					return
				}
				found = true
			}
		}
	}
}

// FirstAndRest returns CAR/CDR, the first item from the iterable followed by an
// iterable for the rest of it.
func FirstAndRest[T any](input iter.Seq[T]) (T, iter.Seq[T]) {
	next, stop := iter.Pull(input)

	first, ok := next()
	if !ok {
		var zeroValue T
		return zeroValue, nil
	}

	return first, func(yield func(T) bool) {
		defer stop()

		for {
			if next, ok := next(); ok {
				if !yield(next) {
					return
				}
			} else {
				return
			}
		}
	}
}

// GroupBy returns an iterator of iterators, with each first-level iteration
// yielding a key and an iterator of values that map to that key value in
// sequence. Like Uniq, if the same key appears disjointedly it will show up
// multiple times in iteration.
func GroupBy[T any, K comparable](keyFunc func(T) K, input iter.Seq[T]) func(func(K, iter.Seq[T]) bool) {
	return func(yield func(K, iter.Seq[T]) bool) {
		next, done := iter.Pull(input)
		currentItem, ok := next()
		key := keyFunc(currentItem)
		defer done()

		if !ok {
			return
		}

		for {
			if !yield(key, func(yieldInner func(value T) bool) {
				yieldInner(currentItem)
				for {
					currentItem, ok = next()
					if !ok {
						return
					}

					newKey := keyFunc(currentItem)
					if newKey != key {
						key = newKey
						return
					} else {
						yieldInner(currentItem)
					}
				}
			}) {
				return
			}

			if !ok {
				return
			}
		}
	}
}
