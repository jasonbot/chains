package chains

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	toFilter := []int{1, 2, 3, 4}
	filtered := []int{3, 4}

	if !slices.Equal[[]int](
		ChainFromSlice(toFilter).Filter(func(i int) bool { return i >= 3 }).Slice(),
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
	returnCodes := []int{
		200, 201, 202, 200, 200,
		302, 301,
		403,
		200, 210,
		550, 500, 535,
		200,
		404, 404, 404,
	}

	counts := map[int]int{}
	expectedCounts := map[int]int{
		200: 8,
		300: 2,
		400: 4,
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

func TestAllStreetFighterMatches(t *testing.T) {
	regularFighters := []string{"Ryu", "Chun Li", "Guile", "E. Honda"}
	bosses := []string{"Sagat", "Vega", "M. Bison"}

	allExpectedFights := []string{
		"Ryu vs. Chun Li",
		"Ryu vs. Guile",
		"Ryu vs. E. Honda",
		"Chun Li vs. Guile",
		"Chun Li vs. E. Honda",
		"Guile vs. E. Honda",
		"Ryu vs. Sagat",
		"Chun Li vs. Sagat",
		"Guile vs. Sagat",
		"E. Honda vs. Sagat",
		"Ryu vs. Vega",
		"Chun Li vs. Vega",
		"Guile vs. Vega",
		"E. Honda vs. Vega",
		"Ryu vs. M. Bison",
		"Chun Li vs. M. Bison",
		"Guile vs. M. Bison",
		"E. Honda vs. M. Bison",
	}

	// Each combination of players without replacement
	singlePlayerFights := ChainJunctionFromIterator[[]string, string](
		CombinationsOfLength(regularFighters, 2),
	).Map(
		func(names []string) string {
			return strings.Join(names, " vs. ")
		},
	)

	// Trick to get pairwise fights from two lists -- lengthen the one by
	// the number of elements in the other, then cycle.
	bossFights := Chain2FromIterator(
		Zip(
			Cycle(Each(regularFighters)),
			Lengthen(
				Each(bosses),
				len(regularFighters),
			),
		),
	).Map(
		func(p1, p2 string) string {
			return strings.Join([]string{p1, p2}, " vs. ")
		},
	)

	allFights := ChainFromIterator(
		FlattenArgs(
			singlePlayerFights.Each(),
			bossFights.Each(),
		),
	).Slice()

	if !slices.Equal(allFights, allExpectedFights) {
		t.Fatalf("%v != %v", allFights, allExpectedFights)
	}
}
