package chains

import (
	"iter"
)

// IterableSequence2 is an opaque wrapper on an iterator to allow for chained methods.
type IterableSequence2[T, V any] struct {
	iterable iter.Seq2[T, V]
}

// IterableSequenceJunction is an opaque wrapper on an iterator to allow for chained methods,
// useful when going from one type to another like doing a .Map from int to string.
type IterableSequenceJunction2[T any, V any, K comparable] struct {
	iterable iter.Seq2[T, V]
}

// ChainJunction2 is used to add a third type to the chain; e.g. to map to an unrelated type.
func ChainJunction2[T any, V, K comparable](in *IterableSequence2[T, V]) *IterableSequenceJunction2[T, V, K] {
	return &IterableSequenceJunction2[T, V, K]{
		iterable: in.iterable,
	}
}

// Each is the final point to get an iterator out of an IterableSequence.
// After chaining your various .Map(...).Filter(..)... do a `range .Each()`
// to iterate over it in your code.
func (iter *IterableSequence2[T, V]) Each(yield func(T, V) bool) {
	if iter == nil {
		return
	}

	for t, v := range iter.iterable {
		if !yield(t, v) {
			return
		}
	}
}

// Tap is a borrowed Rubyism -- it takes each item and passes it along, but
// feeds it to a function to visit first. Useful for calling method, sanitizing
// fields, etc.
func (iter *IterableSequence2[T, V]) Tap(visitor func(T, V)) *IterableSequence2[T, V] {
	if iter == nil {
		return nil
	}

	iter.iterable = Tap2(iter.iterable, visitor)
	return iter
}

// Map is the classic function map -- takes a function, applies it to each
// item in the iterator, and yields that result
func (iter *IterableSequence2[T, V]) Map(mapFunc func(T, V) V) *IterableSequence[V] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[V]{
		iterable: Map2(iter.iterable, mapFunc),
	}
}

func (iter *IterableSequence2[T, V]) Filter(filterFunc func(T, V) bool) *IterableSequence2[T, V] {
	iter.iterable = Filter2(iter.iterable, filterFunc)
	return iter
}

func (iter *IterableSequence2[T, V]) All(predicateFunc func(T, V) bool) bool {
	return All2(iter.iterable, predicateFunc)
}

func (iter *IterableSequence2[T, V]) Any(predicateFunc func(T, V) bool) bool {
	return Any2(iter.iterable, predicateFunc)
}

func (iter *IterableSequence2[T, V]) Count() int {
	var zeroValue V
	return Count(Map2(iter.iterable, func(T, V) V { return zeroValue }))
}

func (iter *IterableSequenceJunction2[T, V, K]) Chain() *IterableSequence2[T, V] {
	if iter == nil {
		return nil
	}

	return &IterableSequence2[T, V]{
		iterable: iter.iterable,
	}
}

func (iter *IterableSequenceJunction2[T, V, K]) Map(mapFunc func(T, V) K) *IterableSequence[K] {
	if iter == nil {
		return nil
	}

	return &IterableSequence[K]{
		iterable: Map2(iter.iterable, mapFunc),
	}
}