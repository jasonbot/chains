package chains

import (
	"fmt"
	"testing"
)

func TestGroupBy(t *testing.T) {
	intSeq := Each([]int{100, 101, 200, 202, 203, 225, 201, 300, 303, 399})

	for key, items := range GroupBy(func(i int) int { return i - (i % 100) }, intSeq) {
		fmt.Println("KEY", key)
		for item := range items {
			fmt.Println("    ITEM:", item)
		}
	}
}
