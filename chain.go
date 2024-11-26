package chains

import (
	"iter"
)

type IterableSequence[T any] struct {
	iterable iter.Seq[T]
}

type IterableSequence2[T, V any] struct {
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

func ChainFromIterator[T any](inFunc func(func(T) bool)) *IterableSequence[T] {
	return &IterableSequence[T]{
		iterable: inFunc,
	}
}

func Chain2[T, V any](in []T) *IterableSequence2[T, V] {
	return &IterableSequence2[T, V]{
		iterable: func(yield func(T) bool) {
			for _, v := range in {
				if !yield(v) {
					return
				}
			}
		},
	}
}

func Chain2FromIterator[T, V any](inFunc func(func(T) bool)) *IterableSequence2[T, V] {
	return &IterableSequence2[T, V]{
		iterable: inFunc,
	}
}

func Junction2[T, V any](in *IterableSequence[T]) *IterableSequence2[T, V] {
	return &IterableSequence2[T, V]{
		iterable: in.iterable,
	}
}

func (iter *IterableSequence[T]) Each(yield func(T) bool) {
	if iter == nil {
		return
	}

	for v := range iter.iterable {
		if !yield(v) {
			return
		}
	}
}

func (iter *IterableSequence[T]) Tap(visitor func(T)) *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	iter.iterable = Tap(visitor, iter.iterable)
	return iter
}

func (iter *IterableSequence[T]) Map(mapFunc func(T) T) *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	iter.iterable = Map(mapFunc, iter.iterable)
	return iter
}

func (iter *IterableSequence[T]) Reduce(reduceFunc func(T, T) T) T {
	if iter == nil {
		var zeroValue T
		return zeroValue
	}

	return Reduce(reduceFunc, iter.iterable)
}

func (iter *IterableSequence[T]) ReduceWithZero(reduceFunc func(T, T) T, zeroValue T) T {
	if iter == nil {
		var zeroValue T
		return zeroValue
	}

	return ReduceWithZero(reduceFunc, zeroValue, iter.iterable)
}

func (iter *IterableSequence[T]) Filter(filterFunc func(T) bool) *IterableSequence[T] {
	iter.iterable = Filter(filterFunc, iter.iterable)
	return iter
}

func (iter *IterableSequence[T]) A() []T {
	if iter == nil {
		return nil
	}

	returnValues := make([]T, 0, 100)

	for item := range iter.iterable {
		returnValues = append(returnValues, item)
	}

	return returnValues
}

func (iter *IterableSequence2[T, V]) Chain() *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[T]{
		iterable: iter.iterable,
	}
}

func (iter *IterableSequence2[T, V]) Map(mapFunc func(T) V) *IterableSequence[V] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[V]{
		iterable: Map(mapFunc, iter.iterable),
	}
}

func (iter *IterableSequence2[T, V]) Reduce(reduceFunc func(V, T) V) V {
	if iter == nil {
		var zeroValue V
		return zeroValue
	}

	var zeroValue V
	return ReduceWithZero(reduceFunc, zeroValue, iter.iterable)
}

func (iter *IterableSequence2[T, V]) ReduceWithZero(reduceFunc func(V, T) V, zeroValue V) V {
	if iter == nil {
		var zeroValue V
		return zeroValue
	}

	return ReduceWithZero(reduceFunc, zeroValue, iter.iterable)
}

func (iter *IterableSequence2[T, V]) A() []T {
	if iter == nil {
		return nil
	}

	returnValues := make([]T, 0, 100)

	for item := range iter.iterable {
		returnValues = append(returnValues, item)
	}

	return returnValues
}
