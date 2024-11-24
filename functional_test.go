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
			t.Fatalf("Items at index %v not equal: %v, %v", index, val, secondArray[index])
		}
	}
}
