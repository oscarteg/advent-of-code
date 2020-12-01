package main

import (
	"github.com/oscarteg/advent-of-code-2020/src/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDayOne(t *testing.T) {
	t.Run("it should use brute force", func(t *testing.T) {
		input := []int{1721, 979, 366, 299, 675, 1456}
		assert.Equal(t, FindMultipliesToSimple(input, 2020), 514579)
	})

	t.Run("it should use brute force for input ", func(t *testing.T) {
		input, _ := utils.ReadLines("./day_1_input.txt")
		assert.Equal(t, FindMultipliesToSimple(input, 2020), 319531)
	})
}
