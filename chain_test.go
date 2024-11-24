package chains

import (
	"testing"
)

func TestMapBasic(t *testing.T) {
	mapFunc := func(i int) int { return i * 2 }
	array := []int{1, 2, 3, 4}
	secondArray := Chain(array).Map(mapFunc).A()

	if len(array) != len(secondArray) {
		t.Fatalf("Array %v not same length as %v", array, secondArray)
	}

	for index, val := range array {
		if mapFunc(val) != secondArray[index] {
			t.Fatalf("Items at index %v not equal: %v, %v", index, mapFunc(val), secondArray[index])
		}
	}
}

func TestReduceBasic(t *testing.T) {
	array := []int{1, 2, 3, 4}
	reduceFunc := func(i, j int) int { return i + j }
	endValue := 10

	sumTotal := Chain(array).Reduce(reduceFunc)
	if sumTotal != endValue {
		t.Fatalf("Reduce failed: %v != %v", sumTotal, endValue)
	}
}

func TestFilterBasic(t *testing.T) {
	filterFunc := func(i int) bool { return i%2 == 0 }
	array := []int{1, 2, 3, 4}
	expectedOutput := []int{2, 4}
	secondArray := Chain(array).Filter(filterFunc).A()

	if len(secondArray) != len(expectedOutput) {
		t.Fatal("Array not right shape")
	}

	for index, val := range expectedOutput {
		if val != secondArray[index] {
			t.Fatalf("Items at index %v not equal: %v, %v", index, val, secondArray[index])
		}
	}
}

func TestChainSome(t *testing.T) {
	array := []int{8, 10, 145, 3}
	secondArray := Chain(
		array,
	).Map(
		func(i int) int {
			return i * 3
		},
	).Filter(
		func(i int) bool {
			return i%2 == 0
		},
	).A()
	expectedOutput := []int{24, 30}

	if len(secondArray) != len(expectedOutput) {
		t.Fatal("Array not right shape")
	}

	for index, val := range expectedOutput {
		if val != secondArray[index] {
			t.Fatalf("Items at index %v not equal: %v, %v", index, val, secondArray[index])
		}
	}

	expectedSum := 54

	if sum := Chain(secondArray).Reduce(func(a, b int) int { return a + b }); sum != expectedSum {
		t.Fatalf("Sums not expected: %v != %v", sum, 18)
	}
}
