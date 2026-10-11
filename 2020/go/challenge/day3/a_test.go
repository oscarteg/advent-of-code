package day3

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `..##.......
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
	input := challenge.FromLiteral(example)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 7, result)
}

func TestAEmptyInput(t *testing.T) {
	input := challenge.FromLiteral("")
	_, err := A(input)
	require.Error(t, err)
}

func TestASingleLine(t *testing.T) {
	input := challenge.FromLiteral("#..")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 1, result)
}
