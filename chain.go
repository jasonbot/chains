package chains

import (
	"iter"
)

// IterableSequence is an opaque wrapper on an iterator to allow for chained methods.
type IterableSequence[T any] struct {
	iterable iter.Seq[T]
}

// IterableSequenceJunction is an opaque wrapper on an iterator to allow for chained methods,
// useful when going from one type to another like doing a .Map from int to string.
type IterableSequenceJunction[T any, V comparable] struct {
	iterable iter.Seq[T]
}

// ChainFromSlice creates an chainable IterableSequence from a slice.
func ChainFromSlice[T any](in []T) *IterableSequence[T] {
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

// Chain creates an chainable IterableSequence from an existing iterator.
func ChainFromIterator[T any](inFunc func(func(T) bool)) *IterableSequence[T] {
	return &IterableSequence[T]{
		iterable: inFunc,
	}
}

// ChainJunctionFromSlice creates an chainable IterableSequence2 from an existing slice.
func ChainJunctionFromSlice[T any, V comparable](in []T) *IterableSequenceJunction[T, V] {
	return &IterableSequenceJunction[T, V]{
		iterable: func(yield func(T) bool) {
			for _, v := range in {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// ChainJunctionFromIterator creates an chainable IterableSequence2 from an existing iterator.
func ChainJunctionFromIterator[T any, V comparable](inFunc func(func(T) bool)) *IterableSequenceJunction[T, V] {
	return &IterableSequenceJunction[T, V]{
		iterable: inFunc,
	}
}

// ChainJunction is used to go from a single-type chain to a dual-type chain.
// This conversion is needed is doing a Map/Reduce that converts type.
func ChainJunction[T any, V comparable](in *IterableSequence[T]) *IterableSequenceJunction[T, V] {
	return &IterableSequenceJunction[T, V]{
		iterable: in.iterable,
	}
}

// Each is the final point to get an iterator out of an IterableSequence.
// After chaining your various .Map(...).Filter(..)... do a `range .Each()`
// to iterate over it in your code.
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

// Tap is a borrowed Rubyism -- it takes each item and passes it along, but
// feeds it to a function to visit first. Useful for calling method, sanitizing
// fields, etc.
func (iter *IterableSequence[T]) Tap(visitor func(T)) *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	iter.iterable = Tap(iter.iterable, visitor)
	return iter
}

// Map is the classic function map -- takes a function, applies it to each
// item in the iterator, and yields that result
func (iter *IterableSequence[T]) Map(mapFunc func(T) T) *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	iter.iterable = Map(iter.iterable, mapFunc)
	return iter
}

// Reduce is the classic function reduce -- takes a function, applies it to each
// item in the iterator along with its prior value, and yields that result
func (iter *IterableSequence[T]) Reduce(reduceFunc func(T, T) T) T {
	if iter == nil {
		var zeroValue T
		return zeroValue
	}

	return Reduce(iter.iterable, reduceFunc)
}

func (iter *IterableSequence[T]) ReduceWithZero(reduceFunc func(T, T) T, zeroValue T) T {
	if iter == nil {
		var zeroValue T
		return zeroValue
	}

	return ReduceWithZero(iter.iterable, reduceFunc, zeroValue)
}

func (iter *IterableSequence[T]) Filter(filterFunc func(T) bool) *IterableSequence[T] {
	iter.iterable = Filter(iter.iterable, filterFunc)
	return iter
}

func (iter *IterableSequence[T]) All(predicateFunc func(T) bool) bool {
	return All(iter.iterable, predicateFunc)
}

func (iter *IterableSequence[T]) Any(predicateFunc func(T) bool) bool {
	return Any(iter.iterable, predicateFunc)
}

func (iter *IterableSequence[T]) Count() int {
	return Count(iter.iterable)
}

func (iter *IterableSequence[T]) Zip(i *IterableSequence[T]) *IterableSequence2[T, T] {
	return &IterableSequence2[T, T]{
		iterable: Zip(iter.iterable, i.iterable),
	}
}

func (iter *IterableSequence[T]) ZipLongest(zeroValue T, i *IterableSequence[T]) *IterableSequence2[T, T] {
	return &IterableSequence2[T, T]{
		iterable: ZipLongest(iter.iterable, i.iterable, zeroValue, zeroValue),
	}
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

func (iter *IterableSequenceJunction[T, V]) Chain() *IterableSequence[T] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[T]{
		iterable: iter.iterable,
	}
}

func (iter *IterableSequenceJunction[T, V]) Map(mapFunc func(T) V) *IterableSequence[V] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[V]{
		iterable: Map(iter.iterable, mapFunc),
	}
}

func (iter *IterableSequenceJunction[T, V]) Reduce(reduceFunc func(V, T) V) V {
	if iter == nil {
		var zeroValue V
		return zeroValue
	}

	var zeroValue V
	return ReduceWithZero(iter.iterable, reduceFunc, zeroValue)
}

func (iter *IterableSequenceJunction[T, V]) ReduceWithZero(reduceFunc func(V, T) V, zeroValue V) V {
	if iter == nil {
		var zeroValue V
		return zeroValue
	}

	return ReduceWithZero(iter.iterable, reduceFunc, zeroValue)
}

func (iter *IterableSequenceJunction[T, V]) GroupBy(keyFunc func(T) V) iter.Seq2[V, *IterableSequence[T]] {
	return func(yield func(V, *IterableSequence[T]) bool) {
		for k, items := range GroupBy[T, V](keyFunc, iter.iterable) {
			if !yield(k, &IterableSequence[T]{
				iterable: items,
			}) {
				return
			}
		}
	}
}

func (iter *IterableSequenceJunction[T, V]) Zip(i *IterableSequence[V]) iter.Seq2[T, V] {
	return Zip[T, V](iter.iterable, i.iterable)
}

func (iter *IterableSequenceJunction[T, V]) A() []T {
	if iter == nil {
		return nil
	}

	returnValues := make([]T, 0, 100)

	for item := range iter.iterable {
		returnValues = append(returnValues, item)
	}

	return returnValues
}
