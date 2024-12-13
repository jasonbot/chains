package chains

import (
	"cmp"
	"iter"
	"slices"
)

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
