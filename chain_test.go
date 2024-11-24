package chains

import (
	"testing"
)

func TestMapBasic(t *testing.T) {
	mapFunc := func(i int) int { return i * 2 }
	array := []int{1, 2, 3, 4}
	secondArray := Each(array).Map(mapFunc).A()

	if len(array) != len(secondArray) {
		t.Fatalf("Array %v to same length as %v", array, secondArray)
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

	sumTotal := Each(array).Reduce(reduceFunc)
	if sumTotal != endValue {
		t.Fatalf("Reduce failed: %v != %v", sumTotal, endValue)
	}
}

func TestFilterBasic(t *testing.T) {
	filterFunc := func(i int) bool { return i%2 == 0 }
	array := []int{1, 2, 3, 4}
	expectedOutput := []int{2, 4}
	secondArray := Each(array).Filter(filterFunc).A()

	if len(secondArray) != len(expectedOutput) {
		t.Fatal("Array not right shape")
	}

	for index, val := range expectedOutput {
		if val != secondArray[index] {
			t.Fatalf("Items at index %v not equal: %v, %v", index, val, secondArray[index])
		}
	}
}
