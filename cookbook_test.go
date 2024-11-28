package chains

import (
	"fmt"
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

func TestSumAndProductWithReduce(t *testing.T) {
	sumarray := []int{1, 2, 3, 4}
	expectedSum := 10
	expectedProduct := 24

	summation := Reduce(
		Each(sumarray),
		func(a, b int) int {
			return a + b
		},
	)

	product := Reduce(
		Each(sumarray),
		func(a, b int) int {
			return a * b
		},
	)

	product0 := ReduceWithZero(
		Each(sumarray),
		func(a, b int) int {
			return a * b
		},
		100,
	)

	if summation != expectedSum {
		t.Fatalf("%v did not sum to %v", summation, expectedSum)
	}

	if product != expectedProduct {
		t.Fatalf("%v did not multiply to %v", product, expectedProduct)
	}

	if product0 != expectedProduct*100 {
		t.Fatalf("%v did not multiply to %v", product, expectedProduct*100)
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

func TestTeeAndMap(t *testing.T) {
	numbersToCompute := []int{1, 2, 3, 4, 10, 20, 50, 100}
	expectedValues := []string{
		"2 + 10 = 12",
		"4 + 20 = 24",
		"6 + 30 = 36",
		"8 + 40 = 48",
		"20 + 100 = 120",
		"40 + 200 = 240",
		"100 + 500 = 600",
		"200 + 1000 = 1200",
	}
	iter1, iter2 := Tee(Each(numbersToCompute))

	doubler := ChainFromIterator(iter1).Map(func(i int) int { return i * 2 })
	tenner := ChainFromIterator(iter2).Map(func(i int) int { return i * 10 })
	calculatedValues := ChainJunction2[int, int, string](Chain2FromIterator(
		Zip(
			doubler.Each(),
			tenner.Each(),
		),
	)).Map(
		func(a int, b int) string {
			return fmt.Sprintf("%v + %v = %v", a, b, a+b)
		},
	).Slice()
	if !slices.Equal(calculatedValues, expectedValues) {
		t.Fatalf("%v != %v", calculatedValues, expectedValues)
	}
}
