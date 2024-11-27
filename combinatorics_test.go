package chains

import (
	"testing"
)

func TestCombinations(t *testing.T) {
	intSeq := []int{1, 2, 3}
	expectedValues := [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}

	index := 0
	for c := range Combinations(intSeq) {
		for j := range len(c) {
			if c[j] != expectedValues[index][j] {
				t.Fatalf("Arrays %v != %v", c, expectedValues[index][j])
			}
		}
		index += 1
	}
}
