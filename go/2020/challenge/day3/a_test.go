package day3

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

var input = `..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(input)

	result := a(input)
	require.Equal(t, 7, result)
}
