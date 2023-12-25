package test

import (
	"golang.org/x/exp/constraints"
)

func CmpMaps[K constraints.Ordered, V constraints.Ordered](expectedMap map[K]V, actualMap map[K]V) int {
	if len(expectedMap) != len(actualMap) {
		return len(expectedMap) - len(actualMap)
	}
	for k, v := range expectedMap {
		if v != actualMap[k] {
			return -1
		}
	}
	return 0
}
