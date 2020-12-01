package util

func SliceToMap(elements []int) map[int]int {
	elementMap := make(map[int]int)
	for i := 0; i < len(elements); i += 2 {
		elementMap[elements[i]] = elements[i+1]
	}

	return elementMap
}
