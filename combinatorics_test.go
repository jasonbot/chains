package chains

import (
	"slices"
	"testing"
)

func TestOrderings(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	for ordering, expected := range ZipLongest(nil, nil, Orderings(intSeq), Each(expectedValues)) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}
