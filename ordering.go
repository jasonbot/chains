package chains

import (
	"cmp"
	"iter"
	"slices"
)

type cycleItem[T any] struct {
	next func() (T, bool)
	done func()
}

type itemListForCycling[T any] []cycleItem[T]

// RoundRobin takes an arbitrary number of iterators and takes one from each
// until they are all exhausted.
func RoundRobin[T cmp.Ordered](iterators ...iter.Seq[T]) iter.Seq[T] {
	itemList := make(itemListForCycling[T], len(iterators))
	for index, iterable := range iterators {
		next, stop := iter.Pull(iterable)
		itemList[index] = cycleItem[T]{
			next: next,
			done: stop,
		}
	}

	return func(yield func(T) bool) {
		defer func() {
			for _, i := range itemList {
				if i.done != nil {
					i.done()
				}
			}
		}()

		index := 0

		missedItems := 0
		for missedItems < len(itemList) {
			if itemList[index].next != nil {
				if nextVal, ok := itemList[index].next(); ok {
					if !yield(nextVal) {
						return
					}
				} else {
					itemList[index].done = nil
					itemList[index].next = nil
				}
				missedItems = 0
			} else {
				missedItems += 1
			}

			index += 1
			index %= len(itemList)
		}
	}
}

type mergeSortedItem[T cmp.Ordered] struct {
	value T
	next  func() (T, bool)
	done  func()
}

type itemListForSorting[T cmp.Ordered] []mergeSortedItem[T]

// TODO: Make this a proper heap
func sortItems[T cmp.Ordered](itemList itemListForSorting[T]) {
	slices.SortFunc(itemList, func(a, b mergeSortedItem[T]) int {
		if a.value < b.value {
			return -1
		} else if a.value == b.value {
			return 0
		} else {
			return 1
		}
	})
}

// Merged takes an arbitrary number of iterators in ascending order and attempts to merge
// them into a single sorted iterable -- this has similar limitations to Uniq in that if
// the sequences are not ordered, the iterable will not magically be sorted -- it will
// be a best effort.
func Merged[T cmp.Ordered](iterators ...iter.Seq[T]) iter.Seq[T] {
	itemList := make(itemListForSorting[T], len(iterators))
	index := 0
	for _, iterable := range iterators {
		next, stop := iter.Pull(iterable)
		firstVal, ok := next()
		if ok {
			itemList[index] = mergeSortedItem[T]{
				value: firstVal,
				next:  next,
				done:  stop,
			}
			index += 1
		}
	}

	itemList = itemList[:index]
	sortItems(itemList)

	return func(yield func(T) bool) {
		defer func() {
			for _, i := range itemList {
				i.done()
			}
		}()

		for len(itemList) > 0 {
			if ok := yield(itemList[0].value); !ok {
				return
			}

			nextVal, ok := itemList[0].next()
			if ok {
				itemList[0].value = nextVal
				if len(itemList) > 1 {
					if itemList[0].value > itemList[1].value {
						sortItems(itemList)
					}
				}
			} else {
				itemList = itemList[1:]
			}
		}
	}
}
