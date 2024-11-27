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
			func(f bool) bool {
				return f
			},
			ZipLongest(func(v1, v2 int) bool {
				return v1 == v2
			}, 0, 0,
				Each(computedVals),
				Each(vals[index]))) {
			t.Fatalf("Values not equal: %v != %v", computedVals, vals[index])
		}

		index += 1
	}
}
