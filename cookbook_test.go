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

func TestCounterWithGroupBy(t *testing.T) {
	returnCodes := []int{200, 201, 202, 200, 200, 302, 301, 403, 200, 210, 550, 500, 535, 200}

	counts := map[int]int{}
	expectedCounts := map[int]int{
		200: 8,
		300: 2,
		400: 1,
		500: 3,
	}

	r := ChainJunctionFromSlice[int, int](returnCodes).GroupBy(
		func(responseCode int) int {
			return responseCode - (responseCode % 100)
		},
	)
	for httpCodeFamily, codes := range r.Each() {
		counts[httpCodeFamily] += codes.Count()
	}

	if !maps.Equal(counts, expectedCounts) {
		t.Fatalf("%v != %v", counts, expectedCounts)
	}
}
