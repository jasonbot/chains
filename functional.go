package chains

import (
	"iter"
	"slices"
)

type IterableSequence[T, V any] struct {
	iterable iter.Seq[T]
}

func Each[T any](in []T) *IterableSequence[T, V] {
	return &IterableSequence[T]{
		iterable: func(yield func(T) bool) {
			for _, v := range in {
				if !yield(v) {
					return
				}
			}
		},
	}
}

func (iter *IterableSequence[T, V]) Each(yield func(T) bool) {
	for v := range iter.iterable {
		if !yield(v) {
			return
		}
	}
}

func (iter *IterableSequence[T, V]) Map(mapFunc func(T) V) {
	iterFunc := func(yield func(V) bool) {
		for v := range iter.iterable {
			if !yield(v) {
				return
			}
		}
	}

	return
}

func (iter *IterableSequence[T, V]) A() []T {
	returnValues := []T{}
	length := 0

	for item := range iter.iterable {
		if len(returnValues) <= length {
			returnValues = slices.Grow(returnValues, 100)
		}

		returnValues[length] = item
	}

	return returnValues[:length]
}
