package chains

import (
	"iter"
	"slices"
)

type IterableSequence[T any] struct {
	iterable iter.Seq[T]
}

func Chain[T any](in []T) *IterableSequence[T] {
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

func (iter *IterableSequence[T]) Each(yield func(T) bool) {
	for v := range iter.iterable {
		if !yield(v) {
			return
		}
	}
}

func (iter *IterableSequence[T]) Tap(visitor func(T)) *IterableSequence[T] {
	return &IterableSequence[T]{
		iterable: Tap(visitor, iter.iterable),
	}
}

func (iter *IterableSequence[T]) Map(mapFunc func(T) T) *IterableSequence[T] {
	return &IterableSequence[T]{
		iterable: Map(mapFunc, iter.iterable),
	}
}

func (iter *IterableSequence[T]) Reduce(reduceFunc func(T, T) T) T {
	return Reduce(reduceFunc, iter.iterable)
}

func (iter *IterableSequence[T]) Filter(filterFunc func(T) bool) *IterableSequence[T] {
	return &IterableSequence[T]{
		iterable: Filter(filterFunc, iter.iterable),
	}
}

func (iter *IterableSequence[T]) A() []T {
	returnValues := []T{}
	length := 0

	for item := range iter.iterable {
		length += 1
		if cap(returnValues) <= length {
			returnValues = slices.Grow(returnValues, 100)
		}
		returnValues = append(returnValues, item)
	}

	return returnValues
}
