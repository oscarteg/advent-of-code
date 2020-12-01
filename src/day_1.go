package main

func FindMultipliesToSimple(array []int, to int) int {
	for i := 0; i <= len(array); i++ {
		for j := i; j <= len(array) - 1; j++ {
			sum := array[i] + array[j]
			if sum == to {
				return array[i] * array[j]
			}
		}
	}

	return 0
}
