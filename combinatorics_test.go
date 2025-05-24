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

	for ordering, expected := range ZipLongest(Permutations(intSeq), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestAllOrderedPermutations(t *testing.T) {
	intSeq := []int{1, 2, 3, 4, 5, 6}
	expectedValues := [][]int{
		{1},
		{2},
		{3},
		{4},
		{5},
		{6},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
		{3, 4},
		{3, 5},
		{3, 6},
		{4, 5},
		{4, 6},
		{5, 6},
		{1, 2, 3},
		{1, 2, 4},
		{1, 2, 5},
		{1, 2, 6},
		{1, 3, 4},
		{1, 3, 5},
		{1, 3, 6},
		{1, 4, 5},
		{1, 4, 6},
		{1, 5, 6},
		{2, 3, 4},
		{2, 3, 5},
		{2, 3, 6},
		{2, 4, 5},
		{2, 4, 6},
		{2, 5, 6},
		{3, 4, 5},
		{3, 4, 6},
		{3, 5, 6},
		{4, 5, 6},
		{1, 2, 3, 4},
		{1, 2, 3, 5},
		{1, 2, 3, 6},
		{1, 2, 4, 5},
		{1, 2, 4, 6},
		{1, 2, 5, 6},
		{1, 3, 4, 5},
		{1, 3, 4, 6},
		{1, 3, 5, 6},
		{1, 4, 5, 6},
		{2, 3, 4, 5},
		{2, 3, 4, 6},
		{2, 3, 5, 6},
		{2, 4, 5, 6},
		{3, 4, 5, 6},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 6},
		{1, 2, 3, 5, 6},
		{1, 2, 4, 5, 6},
		{1, 3, 4, 5, 6},
		{2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5, 6},
	}

	for ordering, expected := range ZipLongest(AllOrderedPermutations(intSeq), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestOrderedPermutations(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{
		{1},
		{2},
		{3},
		{1, 2},
		{1, 3},
		{2, 3},
	}

	for ordering, expected := range ZipLongest(OrderedPermutations(intSeq, 2), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestOrderedPermutationsOfLength(t *testing.T) {
	intSeq := []int{1, 2, 3, 4, 5, 6}
	expectedValues := [][]int{
		{1, 2, 3},
		{1, 2, 4},
		{1, 2, 5},
		{1, 2, 6},
		{1, 3, 4},
		{1, 3, 5},
		{1, 3, 6},
		{1, 4, 5},
		{1, 4, 6},
		{1, 5, 6},
		{2, 3, 4},
		{2, 3, 5},
		{2, 3, 6},
		{2, 4, 5},
		{2, 4, 6},
		{2, 5, 6},
		{3, 4, 5},
		{3, 4, 6},
		{3, 5, 6},
		{4, 5, 6},
	}

	for ordering, expected := range ZipLongest(OrderedPermutationsOfLength(intSeq, 3), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestPermutationsToLength(t *testing.T) {
	intSeq := []int{1, 2, 3, 4, 5, 6}
	expectedValues := [][]int{
		{1},
		{2},
		{3},
		{4},
		{5},
		{6},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
		{2, 1},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
		{3, 1},
		{3, 2},
		{3, 4},
		{3, 5},
		{3, 6},
		{4, 1},
		{4, 2},
		{4, 3},
		{4, 5},
		{4, 6},
		{5, 1},
		{5, 2},
		{5, 3},
		{5, 4},
		{5, 6},
		{6, 1},
		{6, 2},
		{6, 3},
		{6, 4},
		{6, 5},
		{1, 2, 3},
		{1, 2, 4},
		{1, 2, 5},
		{1, 2, 6},
		{1, 3, 2},
		{1, 3, 4},
		{1, 3, 5},
		{1, 3, 6},
		{1, 4, 2},
		{1, 4, 3},
		{1, 4, 5},
		{1, 4, 6},
		{1, 5, 2},
		{1, 5, 3},
		{1, 5, 4},
		{1, 5, 6},
		{1, 6, 2},
		{1, 6, 3},
		{1, 6, 4},
		{1, 6, 5},
		{2, 1, 3},
		{2, 1, 4},
		{2, 1, 5},
		{2, 1, 6},
		{2, 3, 1},
		{2, 3, 4},
		{2, 3, 5},
		{2, 3, 6},
		{2, 4, 1},
		{2, 4, 3},
		{2, 4, 5},
		{2, 4, 6},
		{2, 5, 1},
		{2, 5, 3},
		{2, 5, 4},
		{2, 5, 6},
		{2, 6, 1},
		{2, 6, 3},
		{2, 6, 4},
		{2, 6, 5},
		{3, 1, 2},
		{3, 1, 4},
		{3, 1, 5},
		{3, 1, 6},
		{3, 2, 1},
		{3, 2, 4},
		{3, 2, 5},
		{3, 2, 6},
		{3, 4, 1},
		{3, 4, 2},
		{3, 4, 5},
		{3, 4, 6},
		{3, 5, 1},
		{3, 5, 2},
		{3, 5, 4},
		{3, 5, 6},
		{3, 6, 1},
		{3, 6, 2},
		{3, 6, 4},
		{3, 6, 5},
		{4, 1, 2},
		{4, 1, 3},
		{4, 1, 5},
		{4, 1, 6},
		{4, 2, 1},
		{4, 2, 3},
		{4, 2, 5},
		{4, 2, 6},
		{4, 3, 1},
		{4, 3, 2},
		{4, 3, 5},
		{4, 3, 6},
		{4, 5, 1},
		{4, 5, 2},
		{4, 5, 3},
		{4, 5, 6},
		{4, 6, 1},
		{4, 6, 2},
		{4, 6, 3},
		{4, 6, 5},
		{5, 1, 2},
		{5, 1, 3},
		{5, 1, 4},
		{5, 1, 6},
		{5, 2, 1},
		{5, 2, 3},
		{5, 2, 4},
		{5, 2, 6},
		{5, 3, 1},
		{5, 3, 2},
		{5, 3, 4},
		{5, 3, 6},
		{5, 4, 1},
		{5, 4, 2},
		{5, 4, 3},
		{5, 4, 6},
		{5, 6, 1},
		{5, 6, 2},
		{5, 6, 3},
		{5, 6, 4},
		{6, 1, 2},
		{6, 1, 3},
		{6, 1, 4},
		{6, 1, 5},
		{6, 2, 1},
		{6, 2, 3},
		{6, 2, 4},
		{6, 2, 5},
		{6, 3, 1},
		{6, 3, 2},
		{6, 3, 4},
		{6, 3, 5},
		{6, 4, 1},
		{6, 4, 2},
		{6, 4, 3},
		{6, 4, 5},
		{6, 5, 1},
		{6, 5, 2},
		{6, 5, 3},
		{6, 5, 4},
	}

	for ordering, expected := range ZipLongest(PermutationsToLength(intSeq, 3), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}

func TestAllCombinations(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{
		{1},
		{1, 2},
		{1, 2, 3},
		{1, 3},
		{1, 3, 2},
		{2},
		{2, 1},
		{2, 1, 3},
		{2, 3},
		{2, 3, 1},
		{3},
		{3, 1},
		{3, 1, 2},
		{3, 2},
		{3, 2, 1},
	}

	for ordering, expected := range ZipLongest(
		AllPermutations(intSeq),
		Each(expectedValues),
		nil,
		nil,
	) {
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

	for ordering, expected := range ZipLongest(PermutationsWithReplacement(intSeq), Each(expectedValues), nil, nil) {
		if !slices.Equal(ordering, expected) {
			t.Fatalf("Arrays %v != %v", ordering, expected)
		}
	}
}
