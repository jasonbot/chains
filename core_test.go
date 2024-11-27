package chains

import (
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
				func(a, b int) bool {
					return a == b
				},
				ZipLongest(
					Each(computedVals),
					Each(vals[index]),
					-1,
					-1,
				),
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
