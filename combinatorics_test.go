package chains

import (
	"slices"
	"testing"
)

func TestPermutations(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	for ordering, expected := range ZipLongest(nil, nil, Permutations(intSeq), Each(expectedValues)) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestPermutationsWithReplacement(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{
		{1, 1, 1},
		{1, 1, 2},
		{1, 1, 3},
		{1, 2, 1},
		{1, 2, 2},
		{1, 2, 3},
		{1, 3, 1},
		{1, 3, 2},
		{1, 3, 3},
		{2, 1, 1},
		{2, 1, 2},
		{2, 1, 3},
		{2, 2, 1},
		{2, 2, 2},
		{2, 2, 3},
		{2, 3, 1},
		{2, 3, 2},
		{2, 3, 3},
		{3, 1, 1},
		{3, 1, 2},
		{3, 1, 3},
		{3, 2, 1},
		{3, 2, 2},
		{3, 2, 3},
		{3, 3, 1},
		{3, 3, 2},
		{3, 3, 3},
	}

	for ordering, expected := range ZipLongest(nil, nil, PermutationsWithReplacement(intSeq), Each(expectedValues)) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}
