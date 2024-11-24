package chains

import (
	"iter"
)

func Map[T, V any](mapFunc func(T) V, input iter.Seq[T]) func(func(V) bool) {
	return func(yield func(V) bool) {
		for v := range input {
			if !yield(mapFunc(v)) {
				return
			}
		}
	}
}

func ReduceWithZero[T, V any](collectFunc func(V, T) V, zeroValue V, input iter.Seq[T]) V {
	for v := range input {
		zeroValue = collectFunc(zeroValue, v)
	}
	return zeroValue
}

func Reduce[T any](collectFunc func(T, T) T, input iter.Seq[T]) T {
	var zeroValue T

	return ReduceWithZero(collectFunc, zeroValue, input)
}

func Filter[T any](filterFunc func(T) bool, input iter.Seq[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		for v := range input {
			if filterFunc(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Tap[T any](vistor func(T), input iter.Seq[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		for v := range input {
			vistor(v)
			if !yield(v) {
				return
			}
		}
	}
}
