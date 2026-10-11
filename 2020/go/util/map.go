package util

func SliceToMap(elements []int) map[int]int {
	elementMap := make(map[int]int)
	for _, element := range elements {
		elementMap[element] = element
	}

	return elementMap
}
