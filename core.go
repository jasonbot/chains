package chains

import (
	"iter"
)

func Each[T any](in []T) func(func(T) bool) {
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

func Map[T, V any](mapFunc func(T) V, input iter.Seq[T]) func(func(V) bool) {
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

func Reduce[T any](collectFunc func(T, T) T, input iter.Seq[T]) T {
	var zeroValue T

	return ReduceWithZero(collectFunc, zeroValue, input)
}

func Filter[T any](filterFunc func(T) bool, input iter.Seq[T]) func(func(T) bool) {
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

func Tap[T any](visitor func(T), input iter.Seq[T]) func(func(T) bool) {
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
