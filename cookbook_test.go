package chains

import (
	"maps"
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	toFilter := []int{1, 2, 3, 4}
	filtered := []int{3, 4}

	if !slices.Equal[[]int](
		ChainFromSlice(toFilter).Filter(func(i int) bool { return i >= 3 }).A(),
		filtered,
	) {
		t.Fatalf("%v did not filter to %v", toFilter, filtered)
	}
}

func TestCounterWithReduce(t *testing.T) {
	toCount := []string{"a", "b", "c", "a", "c", "a", "b", "d", "f"}
	expectedVal := map[string]int{
		"a": 3,
		"b": 2,
		"c": 2,
		"d": 1,
		"f": 1,
	}

	counter := ReduceWithZero(
		Each(toCount),
		func(counter map[string]int, s string) map[string]int {
			counter[s] += 1
			return counter
		},
		map[string]int{},
	)

	if !maps.Equal(counter, expectedVal) {
		t.Fatalf("%v did not filter to %v", toCount, expectedVal)
	}
}
