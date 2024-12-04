package chains

import (
	"slices"
	"testing"
)

func TestGroupBy(t *testing.T) {
	intSeq := Each([]int{100, 101, 200, 202, 203, 225, 201, 300, 303, 399, 199})

	keys := []int{100, 200, 300, 100}
	vals := [][]int{{100, 101}, {200, 202, 203, 225, 201}, {300, 303, 399}, {199}}

	index := 0
	for key, items := range GroupBy(func(i int) int { return i - (i % 100) }, intSeq) {
		if key != keys[index] {
			t.Fatalf("Items not equal: %v != %v", key, keys[index])
		}

		computedVals := []int{}
		for item := range items {
			computedVals = append(computedVals, item)
		}

		if !All(
			Map2(
				ZipLongest(
					Each(computedVals),
					Each(vals[index]),
					-1,
					-1,
				),
				func(a, b int) bool {
					return a == b
				},
			),
			func(f bool) bool {
				return f
			},
		) {
			t.Fatalf("Values not equal: %v != %v", computedVals, vals[index])
		}

		index += 1
	}
}

func TestWindows(t *testing.T) {
	items := []int{1, 2, 3, 4}
	expectedSlidingWindows := [][]int{{1, 2, 3}, {2, 3, 4}}
	expectedWindows := [][]int{{1, 2, 3}, {4}}

	index := 0
	for window := range SlidingWindows(Each(items), 3) {
		if !slices.Equal(window, expectedSlidingWindows[index]) {
			t.Fatalf("%v != %v", window, expectedSlidingWindows[index])
		}
		index += 1
	}

	index = 0
	for window := range Windows(Each(items), 3) {
		if !slices.Equal(window, expectedWindows[index]) {
			t.Fatalf("%v != %v", window, expectedWindows[index])
		}
		index += 1
	}
}
